package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/playwright-community/playwright-go"
)

func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "?", "")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, "&", "and")
	name = strings.ReplaceAll(name, "#", "")
	return strings.TrimSpace(name)
}

func main() {
	pw, err := playwright.Run()
	if err != nil {
		panic(err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		panic(err)
	}

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		StorageStatePath: playwright.String("/auth/auth.json"), // session state
	})
	if err != nil {
		panic(err)
	}

	page, err := context.NewPage()
	if err != nil {
		panic(err)
	}

	startURL := os.Getenv("SPACE_URL")
	if startURL == "" {
		panic("Missing SPACE_URL environment variable")
	}

	fmt.Println("🔗 Starting crawl at:", startURL)

	visited := make(map[string]bool)
	links := map[string]bool{startURL: true}

	converter := markdown.NewConverter("", true, nil)

	for len(links) > 0 {
		var current string
		for url := range links {
			current = url
			break
		}
		delete(links, current)

		if visited[current] {
			continue
		}
		visited[current] = true

		fmt.Println("🌐 Visiting:", current)

		_, err := page.Goto(current, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
			Timeout:   playwright.Float(30000),
		})
		if err != nil {
			fmt.Println("⚠️  Failed to load:", current)
			continue
		}

		// Wait for content
		err = page.WaitForSelector(".wiki-content", playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(10000),
		})
		if err != nil {
			fmt.Println("⚠️  No wiki-content found:", current)
			continue
		}

		html, err := page.InnerHTML(".wiki-content")
		if err != nil {
			fmt.Println("⚠️  Could not extract content:", current)
			continue
		}

		markdown, err := converter.ConvertString(html)
		if err != nil {
			fmt.Println("⚠️  Markdown conversion failed:", current)
			continue
		}

		title, err := page.Title()
		if err != nil || title == "" {
			title = fmt.Sprintf("untitled-%d", time.Now().Unix())
		}

		safeFileName := sanitizeFileName(title) + ".md"
		filePath := filepath.Join("/data", safeFileName)
		err = os.WriteFile(filePath, []byte(markdown), 0644)
		if err != nil {
			fmt.Println("❌ Failed to write file:", filePath)
			continue
		}

		fmt.Println("✅ Saved:", filePath)

		// Discover more links
		anchors, _ := page.QuerySelectorAll("a")
		for _, anchor := range anchors {
			href, _ := anchor.GetAttribute("href")
			if href == "" {
				continue
			}
			// Normalize and validate
			if strings.HasPrefix(href, "/wiki/spaces/") && !strings.Contains(href, "#") {
				fullURL := "https://your-org.atlassian.net" + href
				if !visited[fullURL] {
					links[fullURL] = true
				}
			}
		}
	}

	fmt.Println("🎉 Finished crawling.")
}

