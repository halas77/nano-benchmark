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
		for b.Loop() {
			c := colly.NewCollector()
			c.OnHTML(".quote", func(e *colly.HTMLElement) {

			})
			_ = c.Visit(server.URL)

		}
	})

	b.Run("GoQuery", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			res, err := http.Get(server.URL)
			if err != nil {
				b.Fatal(err)
			}
			doc, err := goquery.NewDocumentFromReader(res.Body)
			_ = res.Body.Close()
			if err != nil {
				b.Fatal(err)
			}

			doc.Find(".quote").Each(func(idx int, s *goquery.Selection) {

			})
		}
	})

	b.Run("NanoScrape", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			doc, err := nano.LoadDocument(server.URL)
			if err != nil {
				b.Fatal(err)
			}

			isCssSearch := false

			if isCssSearch {
				doc.SelectAll(".quote")
			} else {
				name := "div"
				params := []*nano.Attribute{
					{
						Key:   "class",
						Value: ".quote",
					},
				}
				doc.FindAll(name, params)
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
		for b.Loop() {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
			if err != nil {
				b.Fatal(err)
			}

			doc.Find(".quote").Each(func(idx int, s *goquery.Selection) {

			})

		}
	})

	b.Run("NanoScrape", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			doc, err := nano.InitDocument(htmlContent)
			if err != nil {
				b.Fatal(err)
			}

			isCssSearch := false
			if isCssSearch {
				doc.SelectAll(".quote")
			} else {
				name := "div"
				params := []*nano.Attribute{
					{
						Key:   "class",
						Value: ".quote",
					},
				}
				doc.FindAll(name, params)
			}
		}
	})
}
