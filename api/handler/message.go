package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/intrntsrfr/vue-ws-test"
	"github.com/intrntsrfr/vue-ws-test/database"
)

type MessageHandler struct {
	r   *gin.Engine
	db  database.DB
	jwt api.JWTService
	ws  *Hub
}

func NewMessageHandler(r *gin.Engine, db database.DB, jwtService api.JWTService, hub *Hub) {
	h := &MessageHandler{r, db, jwtService, hub}

	g := h.r.Group("/api/messages")
	g.GET("/", h.jwt.IsAuthorized(), h.postMessage())
	g.POST("/", h.getMessages())

	g.POST("/:id/reactions", h.jwt.IsAuthorized(), h.postReaction())
	g.DELETE("/:id/reactions", h.jwt.IsAuthorized(), h.deleteReaction())
}

func (h *MessageHandler) postMessage() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (h *MessageHandler) getMessages() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (h *MessageHandler) postReaction() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (h *MessageHandler) deleteReaction() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
