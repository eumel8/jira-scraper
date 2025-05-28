# jira-scraper

Scrape your confluence wiki data and store them as Markdown

This project is based on [Playwright](https://playwright.dev/), an end-to-end testing for modern web apps.

Target of this project is to build a RAG, similar to [Ollama Chatbot](https://github.com/eumel8/ollama-chatbot) with data from a Confluence Wiki. The Wiki is hosted and secured by a 2FA login mechanism. That's the reason why Playwright comes into. Let's start!

## Generate Login Token

The idea is to manual login into the Confluence with credentials and 2FA to generate an `auth.json` file for later use. Some preparation are required on your Mac or Linux machine, like install GoLang, git, kubectl.

local install Playwright for 2FA Auth login

```bash
git clone https://github.com/eumel8/jira-scraper.git
go install github.com/playwright-community/playwright-go/cmd/playwright@latest
playwright install
cd auth
go mod tidy
go run state_save.go
```

## Generate Auth Secret

hint: The following steps are compute on a Kubernetes cluster.

```bash
cd kubernetes
kubectl create ns jira-scraper
./create-secret.sh
```

## Create PVC

We need some storage to save the Markdown files. Adjust your StorageClass and size in `kubernetes/pvc.yaml` and make a:

```bash
kubectl -n jira-scrapper apply -f pvc.yaml
```

## Start Scrape Job/Cronjob

We can use a Kubernetes Job for one time usage or a CronJob. Adjust the variables for your Confluence Wiki Space in the manifest and make a:

```bash
kubectl -n jira-scrapper apply -f job.yaml
```
or

```bash
kubectl -n jira-scrapper apply -f cronjob.yaml
```

Follow the output on the Pod logs.

Enjoy

## credits

Frank Kloeker f.kloeker@telekom.de

Life is for sharing. If you have an issue with the code or want to improve it, feel free to open an issue or an pull request.

