package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidhyashekar/cloudcart/services/product-service/internal/service"
)

// ProductHandler handles HTTP requests related to products.
type ProductHandler struct {
	productService service.ProductService
}

// NewProductHandler creates a new instance of ProductHandler with the provided ProductService.
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// Create handles the HTTP request for creating a new product.
func (h *ProductHandler) Create(c *gin.Context) {
	var request service.CreateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	product, err := h.productService.CreateProduct(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetAll handles the HTTP request for retrieving all products.
func (h *ProductHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	categoryID, _ := strconv.ParseUint(
		c.DefaultQuery("category_id", "0"),
		10,
		64,
	)

	query := service.ProductQuery{
		Search:     c.Query("search"),
		CategoryID: uint(categoryID),
		Page:       page,
		Limit:      limit,
	}

	products, err := h.productService.GetProducts(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetByID handles the HTTP request for retrieving a product by its ID.
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Update handles the HTTP request for updating an existing product by its ID.
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	var request service.UpdateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	product, err := h.productService.UpdateProduct(
		uint(id),
		request,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Delete handles the HTTP request for deleting a product by its ID.
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id",
		})
		return
	}

	if err := h.productService.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
