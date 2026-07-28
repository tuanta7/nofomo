package bot

import "context"

type Engine interface {
	PlaceOrder(ctx context.Context)
	CancelOrder(ctx context.Context)
	GetPortfolio(ctx context.Context)
}
