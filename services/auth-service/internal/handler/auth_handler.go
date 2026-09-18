package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration requests.
func (h *AuthHandler) Register(c *gin.Context) {

	var request service.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	response, err := h.authService.Register(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles user login requests. It expects a JSON payload with the user's email and password, and returns a JWT token if the credentials are valid.
func (h *AuthHandler) Login(c *gin.Context) {

	var request service.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	response, err := h.authService.Login(request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Refresh handles requests to refresh the access token using a valid refresh token. It expects a JSON payload with the refresh token and returns a new access token if the refresh token is valid.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var request service.RefreshTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	response, err := h.authService.RefreshAccessToken(
		request.RefreshToken,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Me handles requests to retrieve the authenticated user's information. It expects a valid JWT token in the Authorization header and returns the user's ID, email, and role.
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
		"role":    role,
	})
}
