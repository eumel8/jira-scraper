package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	markdown "github.com/JohannesKaufmann/html-to-markdown"
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

type QueueItem struct {
	Url   string
	Depth int
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
	queue := []QueueItem{{Url: startURL, Depth: 0}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.Url] {
			continue
		}
		visited[current.Url] = true

		fmt.Printf("🔗 Crawling (depth %d): %s\n", current.Depth, current.Url)

		_, err := page.Goto(current.Url, playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
			Timeout:   playwright.Float(20000),
		})
		if err != nil {
			fmt.Println("❌ Failed to load:", current.Url)
			continue
		}

		_,err = page.WaitForSelector(".wiki-content", playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(10000),
		})
		if err != nil {
			fmt.Println("⚠️  No wiki-content found:", current.Url)
			continue
		}

		html, _ := page.InnerHTML(".wiki-content")
		md, _ := converter.ConvertString(html)

		title, _ := page.Title()
		if title == "" {
			u, _ := url.Parse(current.Url)
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

		// Only crawl further if we haven't hit max depth
		if current.Depth < 2 {
			anchors, _ := page.QuerySelectorAll("a")
			for _, a := range anchors {
				href, _ := a.GetAttribute("href")
				if href == "" {
					continue
				}

				var full string
				if strings.HasPrefix(href, "/wiki/") {
					full = baseURL + href
				} else if strings.HasPrefix(href, baseURL+"/wiki/") {
					full = href
				} else {
					continue // skip external links
				}

				if !visited[full] {
					queue = append(queue, QueueItem{Url: full, Depth: current.Depth + 1})
				}
			}
		}
	}

	fmt.Println("🎉 Done. Pages saved:", len(visited))
}

