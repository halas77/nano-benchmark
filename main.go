package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/halas77/nano-scrape/engine"
)

const targetURL = "https://quotes.toscrape.com/"

func runColly(url string) ([]map[string]string, time.Duration, error) {
	start := time.Now()
	var quotes []map[string]string
	c := colly.NewCollector()

	c.OnHTML(".quote", func(e *colly.HTMLElement) {
		text := e.ChildText("span.text")
		author := e.ChildText("small.author")
		quotes = append(quotes, map[string]string{
			"text":   text,
			"author": author,
		})
	})

	err := c.Visit(url)
	duration := time.Since(start)
	if err != nil {
		return nil, 0, err
	}
	return quotes, duration, nil
}

func runGoQuery(url string) ([]map[string]string, time.Duration, error) {
	start := time.Now()
	res, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, 0, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, 0, err
	}

	var quotes []map[string]string
	doc.Find(".quote").Each(func(i int, s *goquery.Selection) {
		text := s.Find("span.text").Text()
		author := s.Find("small.author").Text()
		quotes = append(quotes, map[string]string{
			"text":   text,
			"author": author,
		})
	})

	duration := time.Since(start)
	return quotes, duration, nil
}

func runNanoScrape(url string) ([]map[string]string, time.Duration, error) {
	start := time.Now()
	doc, err := engine.LoadDocument(url)
	if err != nil {
		return nil, 0, err
	}

	quotes := doc.SelectAll(".quote")
	mapping := map[string]string{
		"text":   "span.text",
		"author": "small.author",
	}
	mappedData := quotes.Map(mapping)
	duration := time.Since(start)
	return mappedData, duration, nil
}

func main() {
	fmt.Println("Starting Scraping Live Test...")
	fmt.Printf("Target URL: %s\n\n", targetURL)

	// 1. Colly
	fmt.Println("--- Running Colly Scraper ---")
	collyQuotes, collyDuration, err := runColly(targetURL)
	if err != nil {
		log.Printf("Colly error: %v\n", err)
	} else {
		fmt.Printf("Colly found %d quotes in %s\n", len(collyQuotes), collyDuration)
		if len(collyQuotes) > 0 {
			fmt.Printf("Sample quote: %q - By %s\n\n", collyQuotes[0]["text"], collyQuotes[0]["author"])
		}
	}

	// 2. GoQuery
	fmt.Println("--- Running GoQuery Scraper ---")
	goQueryQuotes, goQueryDuration, err := runGoQuery(targetURL)
	if err != nil {
		log.Printf("GoQuery error: %v\n", err)
	} else {
		fmt.Printf("GoQuery found %d quotes in %s\n", len(goQueryQuotes), goQueryDuration)
		if len(goQueryQuotes) > 0 {
			fmt.Printf("Sample quote: %q - By %s\n\n", goQueryQuotes[0]["text"], goQueryQuotes[0]["author"])
		}
	}

	// 3. Nano Scrape
	fmt.Println("--- Running Nano Scrape Scraper ---")
	nanoQuotes, nanoDuration, err := runNanoScrape(targetURL)
	if err != nil {
		log.Printf("Nano Scrape error: %v\n", err)
	} else {
		fmt.Printf("Nano Scrape found %d quotes in %s\n", len(nanoQuotes), nanoDuration)
		if len(nanoQuotes) > 0 {
			fmt.Printf("Sample quote: %q - By %s\n\n", nanoQuotes[0]["text"], nanoQuotes[0]["author"])
		}
	}

	fmt.Println("Live scraping test completed.")
}
