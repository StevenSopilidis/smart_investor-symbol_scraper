package ports

import "github.com/stevensopi/smart_investor/symbol_scraper/internal/core/domain"

type SymbolScraper interface {
	Scrape(symbols []domain.Symbol) []domain.ScrapeResult
}
