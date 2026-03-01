package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL    string
	Email      string
	Token      string
	HTTPClient *http.Client
}

func NewClient(baseURL, email, token string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Email:      email,
		Token:      token,
		HTTPClient: &http.Client{},
	}
}

// Confluence sayfasından dönen sayfa bilgisinin, sadece bizim ihtiyacımız olan kısmını temsil eden yapı
type Page struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Version struct {
		Number int `json:"number"`
	} `json:"version"`
	Body struct {
		Storage struct {
			Value          string `json:"value"`
			Representation string `json:"representation"`
		} `json:"storage"`
	} `json:"body"`
}

func (c *Client) GetPage(pageID string) (*Page, error) {
	url := fmt.Sprintf("%s/rest/api/content/%s?expand=body.storage,version", c.BaseURL, pageID)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}

	req.SetBasicAuth(c.Email, c.Token)
	// Confluence API'sinden JSON formatında veri almak istediğimizi belirtir
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Error response from Confluence API: %s - %s", resp.Status, string(bodyBytes))
	}

	var page Page

	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("Error decoding response: %w", err)
	}

	return &page, nil
}

func (c *Client) UpdatePage(page *Page, newContent string) error {
	url := fmt.Sprintf("%s/rest/api/content/%s", c.BaseURL, page.ID)

	// JSON verisi oluşturmak için kullanılan yapı
	payload := map[string]interface{}{
		"id":    page.ID,
		"type":  page.Type,
		"title": page.Title,
		"version": map[string]interface{}{
			"number": page.Version.Number + 1,
		},
		"body": map[string]interface{}{
			"storage": map[string]interface{}{
				"value":          newContent,
				"representation": "storage",
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}

	req.SetBasicAuth(c.Email, c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sayfa güncellenemedi. Hata: %s", string(bodyBytes))
	}
	return nil

}
