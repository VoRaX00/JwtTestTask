package handler

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/manager")
	{
		auth.POST("/getPair", h.GetPair)
		auth.POST("/signUp", h.SignUp)
		auth.POST("/refreshTokens", h.RefreshTokens)
	}

	return router
}
