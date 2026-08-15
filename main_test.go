package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/halas77/nano-scrape/nano"
)

var htmlContent string

func init() {
	// Read the target page HTML cached locally
	data, err := os.ReadFile("testdata/quotes.html")
	if err != nil {
		panic(err)
	}
	htmlContent = string(data)
}

// BenchmarkHTTPScraping measures the complete workflow:
// HTTP Request (on localhost) -> Load HTML -> Selector Parsing -> Extracting fields.
func BenchmarkHTTPScraping(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, htmlContent)
	}))
	defer server.Close()

	b.ResetTimer()

	b.Run("Colly", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
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
			_ = c.Visit(server.URL)
			if len(quotes) != 10 {
				b.Fatalf("Expected 10 quotes, got %d", len(quotes))
			}
		}
	})

	b.Run("GoQuery", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res, err := http.Get(server.URL)
			if err != nil {
				b.Fatal(err)
			}
			doc, err := goquery.NewDocumentFromReader(res.Body)
			_ = res.Body.Close()
			if err != nil {
				b.Fatal(err)
			}

			var quotes []map[string]string
			doc.Find(".quote").Each(func(idx int, s *goquery.Selection) {
				text := s.Find("span.text").Text()
				author := s.Find("small.author").Text()
				quotes = append(quotes, map[string]string{
					"text":   text,
					"author": author,
				})
			})
			if len(quotes) != 10 {
				b.Fatalf("Expected 10 quotes, got %d", len(quotes))
			}
		}
	})

	b.Run("NanoScrape", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			doc, err := nano.LoadDocument(server.URL)
			if err != nil {
				b.Fatal(err)
			}
			quotes := doc.SelectAll(".quote")
			mapping := map[string]string{
				"text":   "span.text",
				"author": "small.author",
			}
			mappedData := quotes.Map(mapping)
			if len(mappedData) != 10 {
				b.Fatalf("Expected 10 quotes, got %d", len(mappedData))
			}
		}
	})
}

// BenchmarkParsingOnly measures pure CPU/memory efficiency by
// parsing and extracting data from the in-memory HTML string.
func BenchmarkParsingOnly(b *testing.B) {
	b.ResetTimer()

	b.Run("GoQuery", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
			if err != nil {
				b.Fatal(err)
			}

			var quotes []map[string]string
			doc.Find(".quote").Each(func(idx int, s *goquery.Selection) {
				text := s.Find("span.text").Text()
				author := s.Find("small.author").Text()
				quotes = append(quotes, map[string]string{
					"text":   text,
					"author": author,
				})
			})
			if len(quotes) != 10 {
				b.Fatalf("Expected 10 quotes, got %d", len(quotes))
			}
		}
	})

	b.Run("NanoScrape", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			doc, err := nano.InitDocument(htmlContent)
			if err != nil {
				b.Fatal(err)
			}
			quotes := doc.SelectAll(".quote")
			mapping := map[string]string{
				"text":   "span.text",
				"author": "small.author",
			}
			mappedData := quotes.Map(mapping)
			if len(mappedData) != 10 {
				b.Fatalf("Expected 10 quotes, got %d", len(mappedData))
			}
		}
	})
}
