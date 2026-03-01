package main

import (
	"Confluence_Uploader/confluence"
	"Confluence_Uploader/parser"
	"fmt"
	"log"
	"strings"
)

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Configuration file (config.json) could not be loaded : %v", err)
	}

	// Eğer kullanıcı config içini doldurmayı unuttuysa koruma önlemi:
	if config.BaseURL == "" || config.Email == "" || config.Token == "" || config.PageID == "" || config.ReportsDir == "" {
		log.Fatalf("Configuration file (config.json) is missing required fields. Please fill in all fields before running the program.")
	}

	fmt.Printf("HTML Reports Directory: %s \n", config.ReportsDir)

	reports, err := parser.ParseReports(config.ReportsDir)
	if err != nil {
		log.Fatalf("Error parsing reports: %v", err)
	}

	if len(reports) == 0 {
		fmt.Println("No HTML reports found in the specified directory. Please check the path and ensure it contains .html files.")
		return
	}

	fmt.Printf("[2/3] %d HTML report(s) found. Connecting to Confluence API...\n", len(reports))

	client := confluence.NewClient(config.BaseURL, config.Email, config.Token)

	page, err := client.GetPage(config.PageID)
	if err != nil {
		log.Fatalf("Confluence page could not be retrieved: %v", err)
	}

	fmt.Printf("[3/3] '%s' page found (Version: %d). Adding new content...\n", page.Title, page.Version.Number)

	var newContentBuilder strings.Builder

	newContentBuilder.WriteString(page.Body.Storage.Value)
	newContentBuilder.WriteString("<hr/>")

	for _, rep := range reports {
		newContentBuilder.WriteString(fmt.Sprintf("<h3>%s</h3>", rep.Title))
		newContentBuilder.WriteString(rep.Content)
	}

	finalContent := newContentBuilder.String()

	err = client.UpdatePage(page, finalContent)
	if err != nil {
		log.Fatalf("Page update error: %v", err)
	}

	fmt.Println("Page updated successfully with new report content!")
}
