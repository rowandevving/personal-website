package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5"
)

type Website struct {
	Data  WebsiteData `json:"data"`
	Title string      `json:"title"`
}

type WebsiteData struct {
	Domain string `json:"domain"`
	Colour string `json:"colour"`
}

func AllWebsites(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Couldn't connect to database: %v", err)
		return
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select domain, colour from cool_people where approved = true")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Couldn't query database: %v", err)
		return
	}
	websiteDatas, err := pgx.CollectRows(rows, pgx.RowToStructByName[WebsiteData])
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error: %v", err)
		return
	}

	var websites []Website

	for _, websiteData := range websiteDatas {
		u, err := url.Parse(websiteData.Domain)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}
		u.Scheme = "https"

		res, err := http.Get(u.String())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error: %v", err)
			return
		}

		defer res.Body.Close()

		title := doc.Find("head title").Text()

		websites = append(websites, Website{websiteData, title})
	}

	data, err := json.Marshal(websites)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Couldn't marshal JSON: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300, s-maxage=86400, stale-while-revalidate=300")

	fmt.Fprint(w, string(data))
}
