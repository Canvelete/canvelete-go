package canvelete

import (
	"context"
	"fmt"
	"net/url"
)

// BillingInfo represents billing information
type BillingInfo struct {
	Plan               string `json:"plan"`
	Status             string `json:"status"` // active, cancelled, past_due, trialing
	CreditBalance      int    `json:"creditBalance"`
	CreditLimit        int    `json:"creditLimit"`
	NextBillingDate    string `json:"nextBillingDate"`
	CurrentPeriodStart string `json:"currentPeriodStart"`
	CurrentPeriodEnd   string `json:"currentPeriodEnd"`
	CancelAtPeriodEnd  bool   `json:"cancelAtPeriodEnd"`
}

// Invoice represents an invoice
type Invoice struct {
	ID          string  `json:"id"`
	Date        string  `json:"date"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"` // paid, pending, failed
	Description string  `json:"description"`
	PDFURL      string  `json:"pdfUrl,omitempty"`
}

// BillingSummary represents billing summary
type BillingSummary struct {
	TotalSpent     float64 `json:"totalSpent"`
	CurrentMonth   float64 `json:"currentMonth"`
	PreviousMonth  float64 `json:"previousMonth"`
	AverageMonthly float64 `json:"averageMonthly"`
	Currency       string  `json:"currency"`
}

// Seats represents team seats information
type Seats struct {
	Used      int `json:"used"`
	Total     int `json:"total"`
	Available int `json:"available"`
}

// CreditPurchase represents a credit purchase
type CreditPurchase struct {
	ID         string `json:"id"`
	Amount     float64 `json:"amount"`
	Credits    int    `json:"credits"`
	NewBalance int    `json:"newBalance"`
	Status     string `json:"status"`
}

// BillingService handles billing-related operations
type BillingService struct {
	client *Client
}


// GetInfo retrieves billing information
func (s *BillingService) GetInfo(ctx context.Context) (*BillingInfo, error) {
	var result struct {
		Data BillingInfo `json:"data"`
	}
	
	if err := s.client.request(ctx, "GET", "/api/v1/billing/info", nil, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// InvoiceListOptions are options for listing invoices
type InvoiceListOptions struct {
	Page  int
	Limit int
}

// InvoicesResponse is the response from listing invoices
type InvoicesResponse struct {
	Data       []Invoice         `json:"data"`
	Pagination PaginatedResponse `json:"pagination"`
}

// GetInvoices retrieves invoice history
func (s *BillingService) GetInvoices(ctx context.Context, opts *InvoiceListOptions) (*InvoicesResponse, error) {
	params := url.Values{}
	
	if opts != nil {
		if opts.Page > 0 {
			params.Add("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.Limit > 0 {
			params.Add("limit", fmt.Sprintf("%d", opts.Limit))
		}
	}
	
	endpoint := "/api/v1/billing/invoices"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	
	var result InvoicesResponse
	if err := s.client.request(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// GetSummary retrieves billing summary
func (s *BillingService) GetSummary(ctx context.Context) (*BillingSummary, error) {
	var result struct {
		Data BillingSummary `json:"data"`
	}
	
	if err := s.client.request(ctx, "GET", "/api/v1/billing/summary", nil, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// PurchaseCreditsRequest is the request to purchase credits
type PurchaseCreditsRequest struct {
	CreditAmount    int    `json:"creditAmount"`
	PaymentMethodID string `json:"paymentMethodId,omitempty"`
}

// PurchaseCredits purchases additional credits
func (s *BillingService) PurchaseCredits(ctx context.Context, amount int, paymentMethodID string) (*CreditPurchase, error) {
	req := PurchaseCreditsRequest{
		CreditAmount:    amount,
		PaymentMethodID: paymentMethodID,
	}
	
	var result struct {
		Data CreditPurchase `json:"data"`
	}
	
	if err := s.client.request(ctx, "POST", "/api/v1/billing/credits/purchase", req, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// GetSeats retrieves team seats information
func (s *BillingService) GetSeats(ctx context.Context) (*Seats, error) {
	var result struct {
		Data Seats `json:"data"`
	}
	
	if err := s.client.request(ctx, "GET", "/api/v1/billing/seats", nil, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// SeatsRequest is the request to modify seats
type SeatsRequest struct {
	Count int `json:"count"`
}

// AddSeats adds team seats
func (s *BillingService) AddSeats(ctx context.Context, count int) (*Seats, error) {
	req := SeatsRequest{Count: count}
	
	var result struct {
		Data Seats `json:"data"`
	}
	
	if err := s.client.request(ctx, "POST", "/api/v1/billing/seats/add", req, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// RemoveSeats removes team seats
func (s *BillingService) RemoveSeats(ctx context.Context, count int) (*Seats, error) {
	req := SeatsRequest{Count: count}
	
	var result struct {
		Data Seats `json:"data"`
	}
	
	if err := s.client.request(ctx, "DELETE", "/api/v1/billing/seats/remove", req, &result); err != nil {
		return nil, err
	}
	
	return &result.Data, nil
}

// GetPortalURL retrieves the billing portal URL
func (s *BillingService) GetPortalURL(ctx context.Context) (string, error) {
	var result struct {
		URL string `json:"url"`
	}
	
	if err := s.client.request(ctx, "GET", "/api/billing/portal", nil, &result); err != nil {
		return "", err
	}
	
	return result.URL, nil
}