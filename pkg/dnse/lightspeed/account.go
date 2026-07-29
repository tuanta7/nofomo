package lightspeed

import (
	"context"
	"fmt"
)

// Account is a single trading subaccount under the authenticated investor.
type Account struct {
	ID                string `json:"id"`
	DealAccount       bool   `json:"dealAccount"`
	DerivativeAccount bool   `json:"derivativeAccount"`
}

// AccountsResponse is the payload returned by GET /accounts.
type AccountsResponse struct {
	Name        string    `json:"name"`
	CustodyCode string    `json:"custodyCode"`
	InvestorID  string    `json:"investorId"`
	Accounts    []Account `json:"accounts"`
}

// GetAccounts retrieves the investor's identity and trading sub-accounts.
// https://developers.dnse.com.vn/docs/dnse/get-accounts
func (c *Client) GetAccounts(ctx context.Context) (*AccountsResponse, error) {
	var result AccountsResponse
	if err := c.get(ctx, pathAccounts, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// StockBalance is the stock subaccount balance breakdown.
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

// GetBalance retrieves the balance of a trading subaccount.
// https://developers.dnse.com.vn/docs/dnse/get-account-balances
func (c *Client) GetBalance(ctx context.Context, accountNo string) (*AccountBalanceResponse, error) {
	path := fmt.Sprintf(pathAccountBalances, accountNo)
	var result AccountBalanceResponse
	if err := c.get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
