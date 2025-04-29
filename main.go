package main

import (
	"context"
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

	go publishers.RunPublisher(context.Background(), results, publisher, config.SymbolTopic, &wg)

	wg.Wait()
}
