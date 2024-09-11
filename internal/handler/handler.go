package handler

import (
	"JwtTestTask/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/getPair", h.GetPair)
		auth.POST("/signUp", h.SignUp)
		auth.POST("/refreshTokens", h.RefreshTokens)
	}

	return router
}
