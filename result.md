# Benchmark & Performance Analysis: Go Web Scraping Libraries

This document presents a comprehensive performance benchmark comparing three Go HTML scraping and parsing libraries:

1. **Colly** (v2.3.0) — The standard, feature-rich HTTP scraping framework for Go.
2. **GoQuery** (v1.12.0) — The classic, jQuery-like HTML parsing and selector library.
3. **Nano Scrape** (v0.0.0-20260815074329) — A lightweight web scraping library by [@halas77](https://github.com/halas77/nano-scrape).

---

## 📊 Benchmark Results

The benchmarks were run on a **12th Gen Intel(R) Core(TM) i5-12500** CPU under **Linux/amd64** using `go test -bench=. -benchmem -count=5`. All figures are averaged over 5 runs.

### 1. Pure HTML Parsing & CSS Selecting (No Network)

This benchmark measures the pure CPU parsing speed and CSS selector performance. It parses a pre-loaded in-memory HTML document (~11KB from `quotes.toscrape.com`) and extracts the text and author of all 10 quotes.

| Scraper Library | Speed (ns/op)  | Memory (B/op) | Allocations (allocs/op) | Performance Rank                     |
| :-------------- | :------------- | :------------ | :---------------------- | :----------------------------------- |
| **Nano Scrape** | **257,726 ns** | 95,304 B      | **1,464 allocs**        | 🏆 **1st (Fastest & Fewest Allocs)** |
| **GoQuery**     | 303,405 ns     | **94,240 B**  | 1,528 allocs            | 🥈 2nd (Lowest Memory)               |

> Colly uses GoQuery internally to parse HTML, so its parsing performance is equivalent to GoQuery with added framework overhead — it is excluded from this benchmark.

> **Nano Scrape outperforms GoQuery** in pure parsing speed by approximately **15%** and uses **64 fewer allocations** per operation.

---

### 2. Localhost HTTP Request + Parsing Lifecycle

This benchmark measures the entire lifecycle: executing an HTTP GET request (to a zero-latency `httptest.NewServer` serving the 11KB quotes HTML), loading the response, parsing it, and selecting the 10 quotes.

| Scraper Library | Speed (ns/op)  | Memory (B/op) | Allocations (allocs/op) | Performance Rank                     |
| :-------------- | :------------- | :------------ | :---------------------- | :----------------------------------- |
| **Nano Scrape** | **483,119 ns** | 136,903 B     | **1,571 allocs**        | 🏆 **1st (Fastest & Fewest Allocs)** |
| **GoQuery**     | 583,903 ns     | **100,943 B** | 1,604 allocs            | 🥈 2nd (Lowest Memory)               |
| **Colly**       | 812,200 ns     | 146,502 B     | 1,827 allocs            | 🥉 3rd                               |

> **Nano Scrape is the fastest end-to-end**, outpacing GoQuery by ~**17%** and Colly by ~**40%**. GoQuery retains the lowest memory footprint due to streaming the response body directly without buffering.

> [!WARNING]
> Colly pulls in a large transitive dependency tree (protobuf, appengine, robotstxt, etc.) which increases binary size significantly compared to GoQuery and Nano Scrape.
