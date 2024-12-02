package handler

import (
	"JwtTestTask/internal/domain/models"
	"JwtTestTask/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetPair(c *gin.Context) {
	userId := c.Query("id")

	ip := c.ClientIP()
	tokens, err := h.service.GenerateTokens(userId, ip)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"AccessToken":  tokens["access_token"],
		"RefreshToken": tokens["refresh_token"],
	})
}

func (h *Handler) SignUp(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) RefreshTokens(c *gin.Context) {
	var input services.Tokens
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ip := c.ClientIP()

	tokens, err := h.service.RefreshTokens(input, ip)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"AccessToken":  tokens["access_token"],
		"RefreshToken": tokens["refresh_token"],
	})
}
