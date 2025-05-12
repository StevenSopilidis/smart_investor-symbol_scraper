package main

import (
	"context"
	"log"
	"sync"

	"github.com/stevensopi/smart_investor/symbol_scraper/internal/adapters/client"
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

	// get active symbols to scraper
	client, err := client.NewSymbolManagerClient(config)
	if err != nil {
		log.Fatalf("Could not create client to symbol manager service %s\n", err)
	}

	symbols, err := client.GetActiveSymbols(context.Background())
	if err != nil {
		log.Fatal("Could not get active symbols %s\n", err)
	}

	scaper := scrapers.NewAlphaVantageScraper(config, symbols)
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
