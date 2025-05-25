package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"os"
)

func main() {
	pw, err := playwright.Run()
	if err != nil {
		panic(err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	page, err := browser.NewPage()

	// Navigate to login page
	_, err = page.Goto("https://your-org.atlassian.net/wiki/login", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})

	// SSO Login - adjust selectors for your SSO form
	page.Fill("#username", os.Getenv("JIRA_USER"))
	page.Click("#login-submit")
	page.Fill("#password", os.Getenv("JIRA_PASS"))
	page.Click("#login-submit")

	// Wait for dashboard
	page.WaitForSelector(".dashboard")

	// Now navigate to a known page
	_, err = page.Goto("https://your-org.atlassian.net/wiki/spaces/DOC/pages/123456789/Page+Title")
	page.WaitForSelector(".wiki-content")

	content, _ := page.InnerHTML(".wiki-content")

	// Save HTML to file
	os.WriteFile("page.html", []byte(content), 0644)
	fmt.Println("Page HTML saved!")
}

