package main

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

func main() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch Playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // GUI mode
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}

	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	_, err = page.Goto("https://wiki.example.com")
	if err != nil {
		log.Fatalf("could not go to login page: %v", err)
	}

	// Wait for user to manually login
	page.WaitForTimeout(60000) // 1 minutes to log in
	//page.WaitForTimeout(120000) // 2 minutes to log in

	_, err = context.StorageState("auth.json")
	if err != nil {
		log.Fatalf("could not save storage state: %v", err)
	}
}
