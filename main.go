package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/stevensopi/smart_investor/symbol_scraper/internal/adapters/config"
	"github.com/stevensopi/smart_investor/symbol_scraper/internal/adapters/publishers"
	"github.com/stevensopi/smart_investor/symbol_scraper/internal/adapters/scrapers"
	"github.com/stevensopi/smart_investor/symbol_scraper/internal/core/domain"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("---> Could not load config: %s\n", err)
	}

	scaper := scrapers.NewAlphaVantageScraper(config, []domain.Symbol{
		{
			Ticker: "NVDA",
			Active: true,
		},
	})
	results := make(chan domain.ScrapeResult, 100)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		scaper.Scrape(results)
		defer wg.Done()
	}()

	publisher, err := publishers.NewKafkaPublisher(config)
	if err != nil {
		log.Fatalf("---> Could not create kafka publisher: %s\n", err)
	}

	go func(ctx context.Context) {
		for {
			select {
			case msg := <-results:
				data, err := json.Marshal(msg)
				if err != nil {
					log.Printf("---> Could not encode data: %s\n", err)
				}
				publisher.Publish(msg.Ticker, data, config.SymbolTopic)
			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}(context.Background())
}
