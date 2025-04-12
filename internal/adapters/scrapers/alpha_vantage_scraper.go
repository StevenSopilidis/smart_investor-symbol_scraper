package scrapers

import (
	"encoding/json"
	"fmt"
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
}

func NewAlphaVantageScraper(config config.Config, symbols []domain.Symbol) *AlphaVantageScraper {
	return &AlphaVantageScraper{
		apiKey:         config.AlphaVantageApiKey,
		symbols:        symbols,
		running:        true,
		scrapeInterval: config.ScrapeInterval,
	}
}

func (s *AlphaVantageScraper) fetchQuote(
	symbol string,
	wg *sync.WaitGroup,
	results chan<- GlobalQuoteResponse,
) {
	defer wg.Done()
	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		symbol, s.apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("---> Could not fetch data for symbol %s\n", symbol)
		return
	}
	defer resp.Body.Close()

	var quote GlobalQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		fmt.Printf("Error decoding response for symbol %s: %v\n", symbol, err)
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
				fmt.Printf("Error parsing price for %s: %v\n", ticker, err)
				continue
			}

			fmt.Printf("---> Scrapred stick: %s with price: %s", ticker, priceStr)

			results <- domain.ScrapeResult{
				Ticker:       ticker,
				CurrentPrice: price,
			}
		}

		time.Sleep(s.scrapeInterval)
	}
}
