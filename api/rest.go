package api

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Code int

const (
	CodeError Code = 1 << iota
)

type ErrorResponse struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

type Handler struct {
	e  *gin.Engine
	ws *Hub
	db DB
}

type Config struct {
	JwtUtil JWTService
	db      DB
}

func NewHandler(conf *Config) *Handler {
	h := &Handler{
		gin.Default(),
		NewHub(),
		conf.db,
	}

	h.e.Use(Cors())
	go h.ws.Listen()

	NewAuthHandler(h.e, conf.db, conf.JwtUtil)
	//NewGuildHandler(h.e, conf.Discord, conf.JwtUtil)

	h.e.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	h.e.GET("/ws", h.ws.Handler())

	return h
}

func (h *Handler) Run(address string) error {
	return h.e.Run(address)
}

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
