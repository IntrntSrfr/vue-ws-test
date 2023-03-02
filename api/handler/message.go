package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/intrntsrfr/vue-ws-test"
	"github.com/intrntsrfr/vue-ws-test/database"
	"github.com/intrntsrfr/vue-ws-test/structs"
	"net/http"
	"strings"
	"time"
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
	g.POST("/", h.jwt.IsAuthorized(), h.postMessage())
	g.GET("/", h.getMessages())

	g.POST("/:id/reactions", h.jwt.IsAuthorized(), h.postReaction())
	g.DELETE("/:id/reactions", h.jwt.IsAuthorized(), h.deleteReaction())
}

func (h *MessageHandler) postMessage() gin.HandlerFunc {
	type PostMessageBody struct {
		Content string `json:"content"`
	}

	return func(c *gin.Context) {
		var postMessageBody PostMessageBody
		if err := c.BindJSON(&postMessageBody); err != nil || strings.TrimSpace(postMessageBody.Content) == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{CodeError, "bad request"})
			return
		}

		claims, ok := c.MustGet("claims").(*api.UserClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, ErrorResponse{CodeError, "internal server error"})
			return
		}

		user := h.db.FindUserByID(claims.Subject)
		if user == nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{CodeError, "user does not exist"})
			return
		}
		user.Password = ""

		msg, _ := h.db.CreateMessage(&structs.Message{
			ID:        uuid.New(),
			Author:    user,
			Content:   strings.TrimSpace(postMessageBody.Content),
			Timestamp: time.Now(),
			Reactions: []*structs.Reaction{},
		})

		c.JSON(http.StatusOK, msg)
		//h.ws.userMessage(msg)
		_ = h.ws.dispatchEvent(ActionUserMessage, nil, msg)
	}
}

func (h *MessageHandler) getMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		//msgs := h.Messages[len(h.Messages)-util.Min(len(h.Messages), 50):]

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
