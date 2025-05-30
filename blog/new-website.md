---
title: The New Website
date: 2025-05-30
tags: 
 - Web Development
---
After much deliberation on whether or not I should rewrite my whole website again, here we are. And I couldn't be more happy with how it's turned out.

At the time of writing, I still haven't bought a domain yet so you've most likely found your way here via a link. If you didn't, well, hi! I started building the first major iteration of this site back in early 2023, and it ended up looking... pretty terrible. The whole site was a garish pink colour which was painful to look at, and I'm pretty sure I shoehorned every CSS 'trick' I knew into it, from weird parallax effects to hovering cards which followed your mouse.

## The First Rewrite

After I finally decided it was time for a change, I created a new repository and started again. This time, I used [Astro](https://astro.build), which turned out to be really nice to work with. Instead of making my own theme, I decided to use the [Catppuccin](https://catppuccin.com/) colour scheme, which is what I already use in all my day-to-day apps so it felt like a natural choice. It ultimately ended up looking like this:

![A dark-themed personal portfolio page, including social media icons and a plaid sidebar.](/assets/blog/old-website.png)

Why is one-quarter of the screen taken up by a huge sidebar with absolutely no function? I've asked myself the same question.

This version got off the ground pretty quickly compared to my previous attempt, however, progress slowed as I kept procrastinating by starting and then abandoning other overly-ambitious projects.

## One more time!

I began building the website you're currently on earlier in 2025. After hearing about it from a friend, I decided to try out Óscar Otero's [Lume](https://lume.land) Static Site Generator. I think it was one of the nicest experiences I've had building a website.

I found everything I needed for my site in the various plugins that Lume includes out-of-the-box, without having to install any third-party packages.
A few notable plugins include integrations for atomic CSS engines TailwindCSS and UnoCSS, the `metas` plugin which makes it very simple to generate all the important `<head>` tags for SEO, and the `codeHighlight` plugin which, as its name suggests, introduces simple syntax highlighting for all code blocks.

This was also my first time working with the Deno runtime, and I must say that not having to configure Typescript is really nice, along with being able to use all the browser APIs on the server.

You can support Óscar's work [here](https://github.com/sponsors/oscarotero) and you can support Lume on [Open Collective](https://opencollective.com/lume) (not sponsored of course)