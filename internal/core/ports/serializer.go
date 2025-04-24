package ports

import "github.com/stevensopi/smart_investor/symbol_scraper/internal/core/domain"

type ISymbolSerializer interface {
	Serialize(symbols []domain.Symbol) ([]byte, error)
}
