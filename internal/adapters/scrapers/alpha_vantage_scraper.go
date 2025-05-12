package scrapers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/stevensopi/smart_investor/symbol_scraper/internal/adapters/config"
	"github.com/stevensopi/smart_investor/symbol_scraper/internal/core/domain"
)

type GlobalQuoteResponse struct {
	GlobalQuote map[string]string `json:"Global Quote"`
}

type AlphaVantageScraper struct {
	apiKey         string
	symbols        []domain.Symbol
	running        bool
	scrapeInterval time.Duration
	scrapeEndpoint string
}

func NewAlphaVantageScraper(config config.Config, symbols []domain.Symbol) *AlphaVantageScraper {
	return &AlphaVantageScraper{
		apiKey:         config.AlphaVantageApiKey,
		symbols:        symbols,
		running:        true,
		scrapeInterval: config.ScrapeInterval,
		scrapeEndpoint: config.ScrapeEndpoint,
	}
}

func (s *AlphaVantageScraper) fetchQuote(
	symbol string,
	wg *sync.WaitGroup,
	results chan<- GlobalQuoteResponse,
) {
	defer wg.Done()
	url := fmt.Sprintf(
		"http://%s/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		s.scrapeEndpoint, symbol, s.apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("---> Could not fetch data for symbol %s\n", symbol)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("---> Error reading response body for symbol %s: %v\n", symbol, err)
		return
	}

	var quote GlobalQuoteResponse
	if err := json.Unmarshal(body, &quote); err != nil {
		log.Printf("---> Error decoding response for symbol %s: %v\n", symbol, err)
		return
	}

	results <- quote
}

func (s *AlphaVantageScraper) Scrape(results chan<- domain.ScrapeResult) {
	for s.running {
		var wg sync.WaitGroup
		quoteResp := make(chan GlobalQuoteResponse, len(s.symbols))
		for _, symbol := range s.symbols {
			wg.Add(1)
			go s.fetchQuote(symbol.Ticker, &wg, quoteResp)
		}

		wg.Wait()
		close(quoteResp)

		for quote := range quoteResp {
			ticker := quote.GlobalQuote["01. symbol"]
			priceStr := quote.GlobalQuote["05. price"]

			price, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				log.Printf("---> Error parsing price for %s: %v\n", ticker, err)
				continue
			}

			results <- domain.ScrapeResult{
				Ticker:       ticker,
				CurrentPrice: price,
			}
		}

		time.Sleep(s.scrapeInterval)
	}
}

func (s *AlphaVantageScraper) Shutdown() {
	s.running = false
}
