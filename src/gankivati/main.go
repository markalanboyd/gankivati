package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetchDocument(
	client *http.Client,
	method string,
	targetURL string,
	data url.Values,
) (*goquery.Document, error) {
	var req *http.Request
	var err error

	if method == "POST" {
		req, err = http.NewRequest("POST", targetURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	} else {
		req, err = http.NewRequest("GET", targetURL, nil)
	}

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func saveHTMLToFile(doc *goquery.Document, fileName string) error {
	html, err := doc.Html()
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(html)
	if err != nil {
		return err
	}

	fmt.Printf("Saved HTML to %s\n", fileName)
	return nil
}

func main() {
	searchURL := "https://hjp.znanje.hr/index.php?show=search"

	searchTerm := "imati"

	formData := url.Values{
		"word": {searchTerm},
	}

	client := &http.Client{}

	doc, err := fetchDocument(client, "POST", searchURL, formData)
	if err != nil {
		log.Fatal(err)
	}

	detaljnijeURL, exists := doc.Find("a:contains('detaljnije')").First().Attr("href")
	if exists {
		resolvedDetaljnijeURL, err := url.Parse(detaljnijeURL)
		if err != nil {
			log.Fatal(err)
		}
		baseURL, err := url.Parse(searchURL)
		if err != nil {
			log.Fatal(err)
		}
		resolvedDetaljnijeURL = baseURL.ResolveReference(resolvedDetaljnijeURL)
	}

	resolvedDetaljnijeURL, err := url.Parse(detaljnijeURL)
	if err != nil {
		log.Fatal(err)
	}
	baseURL, err := url.Parse(searchURL)
	if err != nil {
		log.Fatal(err)
	}
	resolvedDetaljnijeURL = baseURL.ResolveReference(resolvedDetaljnijeURL)

	detaljnijeDoc, err := fetchDocument(client, "GET", resolvedDetaljnijeURL.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	err = saveHTMLToFile(detaljnijeDoc, "debug.html")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
