package main

import (
	"fmt"
	"net/url"
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
	baseURL := os.Getenv("CONFLUENCE_BASE_URL") // e.g. https://your-org.atlassian.net
	startURL := os.Getenv("SPACE_URL")          // e.g. https://your-org.atlassian.net/wiki/spaces/XYZ/overview

	if baseURL == "" || startURL == "" {
		panic("Missing CONFLUENCE_BASE_URL or SPACE_URL environment variables")
	}

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
		StorageStatePath: playwright.String("/auth/auth.json"),
	})
	if err != nil {
		panic(err)
	}

	page, err := context.NewPage()
	if err != nil {
		panic(err)
	}

	converter := markdown.NewConverter("", true, nil)
	visited := make(map[string]bool)
	toVisit := []string{startURL}

	fmt.Println("📄 Visiting base page:", startURL)

	// Step 1: Visit the base page and collect depth-1 links
	_, err = page.Goto(startURL)
	if err != nil {
		panic(err)
	}

	err = page.WaitForSelector(".wiki-content", playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(15000),
	})
	if err != nil {
		panic("Base page has no .wiki-content")
	}

	anchors, _ := page.QuerySelectorAll("a")
	for _, a := range anchors {
		href, _ := a.GetAttribute("href")
		if href == "" {
			continue
		}

		// Convert relative to absolute
		var full string
		if strings.HasPrefix(href, "/") {
			full = baseURL + href
		} else if strings.HasPrefix(href, baseURL) {
			full = href
		} else {
			continue // skip external
		}

		if !visited[full] {
			toVisit = append(toVisit, full)
		}
	}

	// Step 2: Visit all collected depth-1 links
	for _, link := range toVisit {
		if visited[link] {
			continue
		}
		visited[link] = true

		fmt.Println("🔗 Crawling:", link)
		_, err := page.Goto(link, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
			Timeout:   playwright.Float(20000),
		})
		if err != nil {
			fmt.Println("❌ Failed to load:", link)
			continue
		}

		err = page.WaitForSelector(".wiki-content", playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(10000),
		})
		if err != nil {
			fmt.Println("⚠️  No wiki-content found:", link)
			continue
		}

		html, _ := page.InnerHTML(".wiki-content")
		md, _ := converter.ConvertString(html)

		title, _ := page.Title()
		if title == "" {
			u, _ := url.Parse(link)
			title = u.Path
		}

		filename := sanitizeFileName(title) + ".md"
		filePath := filepath.Join("/data", filename)
		err = os.WriteFile(filePath, []byte(md), 0644)
		if err != nil {
			fmt.Println("❌ Failed to save:", filePath)
			continue
		}

		fmt.Println("✅ Saved:", filePath)
	}

	fmt.Println("✅ Done. Pages scraped:", len(visited))
}

