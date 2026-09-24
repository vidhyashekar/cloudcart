package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ProductResponse represents the response structure for a product.
type productClient struct {
	baseURL string
	client  *http.Client
}

// NewProductClient creates a new instance of ProductClient with the given base URL.
func NewProductClient(baseURL string) ProductClient {
	return &productClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// GetProduct retrieves product details from the product service by its ID.
func (p *productClient) GetProduct(id uint) (*ProductResponse, error) {
	url := fmt.Sprintf("%s/api/v1/products/%d", p.baseURL, id)

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product not found")
	}

	var product ProductResponse

	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

// DecreaseStock decreases the stock quantity of a product by its ID.
func (p *productClient) DecreaseStock(id uint, quantity int) error {
	return p.updateStock(id, quantity, "decrease")
}

// IncreaseStock increases the stock quantity of a product by its ID.
func (p *productClient) IncreaseStock(id uint, quantity int) error {
	return p.updateStock(id, quantity, "increase")
}

// updateStock is a helper method to update the stock quantity of a product by its ID. It sends a PATCH request to the product service with the specified action (decrease or increase).
func (p *productClient) updateStock(id uint, quantity int, action string) error {
	url := fmt.Sprintf(
		"%s/api/v1/products/%d/stock/%s",
		p.baseURL,
		id,
		action,
	)

	body := fmt.Sprintf(`{"quantity":%d}`, quantity)

	req, err := http.NewRequest(
		http.MethodPatch,
		url,
		strings.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to %s stock for product %d",
			action,
			id,
		)
	}

	return nil
}
