package main

import (
	"github.com/playwright-community/playwright-go"
)

func main() {
	pw, _ := playwright.Run()
	browser, _ := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // GUI mode
	})
	context, _ := browser.NewContext()
	page, _ := context.NewPage()

	page.Goto("https://your-org.atlassian.net/wiki/login")
	// Wait for user to manually login
	page.WaitForTimeout(120000) // 2 minutes to log in

	context.StorageState(path("auth.json")) // Save session
}

