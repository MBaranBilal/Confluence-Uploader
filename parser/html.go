package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Report struct {
	Title   string
	Content string
}

func ParseReports(dir string) ([]Report, error) {
	var reports []Report

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Error reading directory: %w", err)
	}

	//range kullanımında index değeri kullanılmazsa "_" ile belirtilir
	for _, file := range files {
		//	Dizinleri ve .html uzantısı olmayan dosyaları atla
		if file.IsDir() || filepath.Ext(file.Name()) != ".html" {
			continue
		}

		filePath := filepath.Join(dir, file.Name())

		baseName := strings.TrimSuffix(file.Name(), ".html")
		title := strings.Title(strings.ReplaceAll(baseName, "-", " "))

		content, err := parseHTMLFiles(filePath)
		if err != nil {
			fmt.Printf("Error parsing file %s: %v\n", file.Name(), err)
			continue
		}

		reports = append(reports, Report{
			Title:   title,
			Content: content,
		})
	}

	return reports, nil

}

func parseHTMLFiles(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Error opening file: %w", err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return "", fmt.Errorf("Error parsing HTML: %w", err)
	}

	var tablesHTML string

	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		tableHTML, err := goquery.OuterHtml(s)
		if err == nil {
			tablesHTML += tableHTML + "<br>" // Tabloyu ayırmak için <br> eklenebilir
		}
	})

	if tablesHTML == "" {
		return "", fmt.Errorf("No tables found in the HTML file")
	}

	return tablesHTML, nil
}
