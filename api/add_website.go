package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var allowedColours = [10]string{
	"#ea9a97",
	"#ebbcba",
	"#f0a96c",
	"#f6c177",
	"#a6da95",
	"#7fb5a0",
	"#9ccfd8",
	"#8ba3d1",
	"#c4a7e7",
	"#eb6f92",
}

func updateDB(domain string, hex string, approved bool) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	var exists int
	err = conn.QueryRow(context.Background(),
		"select 1 from cool_people where domain = $1 limit 1", domain).Scan(&exists)
	if err == nil {
		return errors.New("Website already exists")
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	_, err = conn.Exec(context.Background(),
		"INSERT INTO cool_people (domain, approved, colour) VALUES ($1, $2, $3)",
		domain, approved, hex,
	)
	if err != nil {
		return err
	}

	return nil
}

func triggerWebhook(url string, client *http.Client) error {
	webhookBody := []byte(fmt.Sprintf(
		`{"content": "<@&1418303555500114092> A website needs approving! %s"}`,
		url,
	))

	req, err := http.NewRequest(http.MethodPost, os.Getenv("WEBHOOK_URL"), bytes.NewBuffer(webhookBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func isValidToken(tok string) bool {
	// regex? I barely know 'er!

	if len(tok) < 43 {
		return false
	}
	for i := 0; i < len(tok); i++ {
		c := tok[i]
		switch {
		case 'A' <= c && c <= 'Z':
		case 'a' <= c && c <= 'z':
		case '0' <= c && c <= '9':
		case c == '-' || c == '_':
		default:
			return false
		}
	}
	return true
}

func AddWebsite(w http.ResponseWriter, r *http.Request) {
	rawDomain := r.URL.Query().Get("domain")
	expectedToken := r.URL.Query().Get("t")
	colour := r.URL.Query().Get("colour")

	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if rawDomain == "" || expectedToken == "" {
		http.Error(w, "Missing domain or token", http.StatusBadRequest)
		return
	}

	if !strings.Contains(rawDomain, "://") {
		rawDomain = "https://" + rawDomain
	}

	u, err := url.Parse(rawDomain)
	if err != nil || u.Host == "" || u.User != nil {
		http.Error(w, "Invalid domain", http.StatusBadRequest)
		return
	}

	target := u.ResolveReference(&url.URL{Path: "/.well-known/rowan"}).String()

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(target)
	if err != nil || res.StatusCode != http.StatusOK {
		http.Error(w, "Verification file not found", http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Read failed", http.StatusInternalServerError)
		return
	}

	if !isValidToken(string(bytes.TrimSpace(body))) || !bytes.Equal(bytes.TrimSpace(body), []byte(expectedToken)) {
		http.Error(w, "Token mismatch", http.StatusUnauthorized)
		return
	}

	clrIndex, err := strconv.Atoi(colour)
	if err != nil || clrIndex < 0 || clrIndex >= len(allowedColours) {
		http.Error(w, fmt.Sprintf("Colour must be a number from 0 to %d", len(allowedColours)-1), http.StatusBadRequest)
		return
	}

	hexColour := allowedColours[clrIndex]

	susDomains, err := client.Get(os.Getenv("SUS_DOMAINS_URL"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusBadGateway)
		log.Printf("ERROR: You goofed up big time: %v", err)
		return
	}
	decoder := json.NewDecoder(susDomains.Body)

	var susDomainsList []string

	err = decoder.Decode(&susDomainsList)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("ERROR: You goofed up big time: %v", err)
		return
	}

	approved := true

	for _, susDomain := range susDomainsList {
		if strings.HasSuffix(u.Hostname(), susDomain) {
			approved = false

			err := triggerWebhook(u.String(), client)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				log.Printf("ERROR: Couldn't trigger webhook: %v", err)
				return
			}
		}
	}

	err = updateDB(u.Hostname(), hexColour, approved)
	if err != nil {
		http.Error(w, "Couldn't update database", http.StatusInternalServerError)
		log.Printf("ERROR: Couldn't update database: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`"Success!"`))
}
