# Benchmark & Performance Analysis: Go Web Scraping Libraries

This document presents a comprehensive performance benchmark comparing three Go HTML scraping and parsing libraries:

1. **Colly** (v2.3.0) — The standard, feature-rich HTTP scraping framework for Go.
2. **GoQuery** (v1.12.0) — The classic, jQuery-like HTML parsing and selector library.
3. **Nano Scrape** (v0.0.0, Local) — A custom lightweight web scraping engine developed by us.

---

## 📊 Benchmark Results

The benchmarks were run on a **12th Gen Intel(R) Core(TM) i5-12500** CPU under **Linux/amd64** using `go test -bench=. -benchmem`.

### 1. Pure HTML Parsing & CSS Selecting (No Network)

This benchmark measures the pure CPU parsing speed and CSS selector performance. It parses a pre-loaded in-memory HTML document (~11KB from `quotes.toscrape.com`) and extracts the text and author of all 10 quotes.

| Scraper Library          | Speed (ns/op)  | Memory (B/op) | Allocations (allocs/op) | Performance Rank                     |
| :----------------------- | :------------- | :------------ | :---------------------- | :----------------------------------- |
| **Nano Scrape (Engine)** | **233,318 ns** | 95,304 B      | **1,464 allocs**        | 🏆 **1st (Fastest & Fewest Allocs)** |
| **GoQuery**              | 256,706 ns     | **94,240 B**  | 1,528 allocs            | 🥈 2nd                               |

> Colly uses goquery to parse the HTML document, so it has similar benchmark results to goquery, but it is slightly slower due to the overhead of the Colly framework.

> **Nano Scrape outperforms GoQuery** in pure parsing and query selection! It is approximately **9% faster** and uses **64 fewer allocations** per document parse and query operation.

---

### 2. Localhost HTTP Request + Parsing Lifecycle

This benchmark measures the entire lifecycle: executing an HTTP GET request (to a zero-latency `httptest.NewServer` serving the 11KB quotes HTML), loading the response, parsing it, and selecting the 10 quotes.

| Scraper Library | Speed (ns/op)  | Memory (B/op) | Allocations (allocs/op) | Performance Rank                     |
| :-------------- | :------------- | :------------ | :---------------------- | :----------------------------------- |
| **Nano Scrape** | **515,724 ns** | 136,811 B     | **1,570 allocs**        | 🏆 **1st (Fastest & Fewest Allocs)** |
| **GoQuery**     | 566,315 ns     | **100,907 B** | 1,604 allocs            | 🥈 2nd (Lowest Memory)               |
| **Colly**       | 740,527 ns     | 146,208 B     | 1,827 allocs            | 🥉 3rd                               |

> **Nano Scrape now outperforms GoQuery in speed** in the HTTP Request lifecycle benchmark as well, while maintaining fewer allocations. However, GoQuery still has a lower memory footprint (`B/op`) due to streaming response bodies directly.
