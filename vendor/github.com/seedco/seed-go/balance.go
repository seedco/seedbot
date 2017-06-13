package seed

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Balance struct {
	CheckingAccountID string `json:"checking_account_id"`
	TotalAvailable    int64  `json:"total_available"`
	Settled           int64  `json:"settled"`
	PendingCredits    uint64 `json:"pending_credits"`
	PendingDebits     uint64 `json:"pending_debits"`
	ScheduledDebits   uint64 `json:"scheduled_debits"`
	Accessible        int64  `json:"accessible"`
	Lockbox           uint64 `json:"lockbox"`
}

type BalanceRequest struct {
	CheckingAccountID string
	Client            *Client
}

type BalanceResponse struct {
	Errors  ErrorList `json:"errors"`
	Results []Balance `json:"results"`
}

func (c *Client) NewBalanceRequest() *BalanceRequest {
	return &BalanceRequest{Client: c}
}

func (b *BalanceRequest) Get() (Balance, error) {
	var response BalanceResponse
	var req *http.Request
	var err error

	if req, err = http.NewRequest("GET", fmt.Sprintf("%s/%s", ApiBase, "balance"), nil); err != nil {
		return Balance{}, err
	}
	var resp *http.Response

	if resp, err = b.Client.do(req); err != nil {
		return Balance{}, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Balance{}, err
	}

	balances := response.Results

	if len(balances) == 0 {
		return Balance{}, errors.New("no balance found")
	}

	return balances[0], nil
}
