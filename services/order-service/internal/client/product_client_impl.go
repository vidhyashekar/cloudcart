package client

import (
	"encoding/json"
	"fmt"
	"net/http"
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
