package handler

import (
	"JwtTestTask/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) SignIn(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ip := c.ClientIP()
	tokens, err := h.service.GenerateTokens(input, ip)
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

func (h *Handler) RefreshToken(c *gin.Context) {

}
