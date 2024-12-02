package handler

import (
	"JwtTestTask/internal/domain/models"
	"JwtTestTask/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary GetPair
// @Tags auth
// @Description Creating a pair of tokens
// @ID get-pair
// @Accept json
// @Produce json
// @Param id query string false "User id"
// @Success 200 {object} services.Tokens
// @Failure 500 {object} ErrorResponse
// @Router /manager/getPair [post]
func (h *Handler) GetPair(c *gin.Context) {
	userId := c.Query("id")

	ip := c.ClientIP()
	tokens, err := h.service.GenerateTokens(userId, ip)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// @Summary SignUp
// @Tags auth
// @Description User registration
// @ID signUp
// @Accept json
// @Produce json
// @Param input body services.RegisterUser true "User data"
// @Success 200 {object} SuccessID
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /manager/signUp [post]
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

	c.JSON(http.StatusOK, SuccessID{
		ID: id,
	})
}

// @Summary RefreshTokens
// @Tags auth
// @Description Refresh tokens
// @ID refresh
// @Accept json
// @Produce json
// @Param input body services.Tokens true "User data"
// @Success 200 {object} services.Tokens
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /manager/refreshTokens [post]
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

	c.JSON(http.StatusOK, tokens)
}
