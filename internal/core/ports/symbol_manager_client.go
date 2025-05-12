package ports

import (
	"context"

	"github.com/stevensopi/smart_investor/symbol_scraper/internal/core/domain"
)

type ISymbolManagerClient interface {
	GetActiveSymbols(ctx context.Context) ([]domain.Symbol, error)
}
