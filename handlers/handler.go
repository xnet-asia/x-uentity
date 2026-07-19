package handlers

import "github.com/gin-gonic/gin"

// Handler defines the interface for HTTP handlers
type Handler interface {
	RegisterRoutes(router *gin.Engine)
}

// BaseHandler provides common handler utilities
type BaseHandler struct {
	engine *gin.Engine
}

func NewBaseHandler(engine *gin.Engine) *BaseHandler {
	return &BaseHandler{engine: engine}
}

func (h *BaseHandler) RegisterRoutes(router *gin.Engine) {
	// Base implementation
}
