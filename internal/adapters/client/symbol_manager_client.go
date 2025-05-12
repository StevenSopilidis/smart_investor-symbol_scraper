package client

import (
	"context"

	"github.com/stevensopi/smart_investor/symbol_scraper/internal/adapters/config"
	"github.com/stevensopi/smart_investor/symbol_scraper/internal/core/domain"
	"github.com/stevensopi/smart_investor/symbol_scraper/internal/core/ports"
	"github.com/stevensopi/smart_investor/symbol_scraper/internal/infrastructure/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SymbolManagerClient struct {
	client grpc_client.ApiClient
}

func NewSymbolManagerClient(config config.Config) (ports.ISymbolManagerClient, error) {
	conn, err := grpc.NewClient(
		config.SymbolManagerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	client := grpc_client.NewApiClient(conn)
	return &SymbolManagerClient{
		client: client,
	}, nil
}

func (s *SymbolManagerClient) GetActiveSymbols(ctx context.Context) ([]domain.Symbol, error) {
	res, err := s.client.GetActiveSymbols(ctx, &grpc_client.GetActiveSymbolsRequest{})

	if err != nil {
		return nil, err
	}

	data := make([]domain.Symbol, 0)
	for _, symbol := range res.Symbols {
		data = append(data, domain.Symbol{
			Ticker:   symbol.Ticker,
			Exchange: symbol.Exchange,
		})
	}

	return data, nil
}
