package rest

import (
	"context"
	"fmt"

	"trading-bot/internal/config"
)

// Account is a single trading sub-account under the authenticated investor.
type Account struct {
	ID                string `json:"id"`
	DealAccount       bool   `json:"dealAccount"`
	DerivativeAccount bool   `json:"derivativeAccount"`
}

// AccountsResponse is the payload returned by GET /accounts.
type AccountsResponse struct {
	Accounts    []Account `json:"accounts"`
	CustodyCode string    `json:"custodyCode"`
	InvestorID  string    `json:"investorId"`
	Name        string    `json:"name"`
}

// GetAccounts retrieves the investor's identity and trading sub-accounts.
// https://developers.dnse.com.vn/docs/dnse/get-accounts
func (c *Client) GetAccounts(ctx context.Context) (*AccountsResponse, error) {
	var result AccountsResponse
	if err := c.get(ctx, config.Accounts, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// StockBalance is the stock sub-account balance breakdown.
type StockBalance struct {
	AvailableCash         int64 `json:"availableCash"`
	CashDividendReceiving int64 `json:"cashDividendReceiving"`
	DepositFeeAmount      int64 `json:"depositFeeAmount"`
	DepositInterest       int64 `json:"depositInterest"`
	TotalCash             int64 `json:"totalCash"`
	TotalDebt             int64 `json:"totalDebt"`
	WithdrawableCash      int64 `json:"withdrawableCash"`
}

// AccountBalanceResponse is the payload returned by GET /accounts/{accountNo}/balances.
type AccountBalanceResponse struct {
	Stock *StockBalance `json:"stock"`
}

// GetBalance retrieves the balance of a trading sub-account.
// https://developers.dnse.com.vn/docs/dnse/get-balance
func (c *Client) GetBalance(ctx context.Context, accountNo string) (*AccountBalanceResponse, error) {
	path := fmt.Sprintf(config.AccountBalances, accountNo)
	var result AccountBalanceResponse
	if err := c.get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
