package handler

import (
	"JwtTestTask/models"
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

	id, err := h.service.Create(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) RefreshTokens(c *gin.Context) {
	token := c.Query("refresh_token")

	ip := c.ClientIP()
	tokens, err := h.service.RefreshTokens(token, ip)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"AccessToken":  tokens["access_token"],
		"RefreshToken": tokens["refresh_token"],
	})
}
