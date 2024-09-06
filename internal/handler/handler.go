package handler

import (
	"JwtTestTask/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/SignIn", h.SignIn)
		auth.POST("/SignUp", h.SignUp)
		auth.POST("/RefreshToken", h.RefreshToken)
	}

	return router
}
