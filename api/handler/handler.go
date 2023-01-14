package handler

import (
	"context"
	"fmt"
	"github.com/intrntsrfr/vue-ws-test"
	"github.com/intrntsrfr/vue-ws-test/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	db database.DB
}

type Config struct {
	JwtUtil api.JWTService
	DB      database.DB
}

func NewHandler(conf *Config) *Handler {
	h := &Handler{
		gin.Default(),
		NewHub(conf.DB),
		conf.DB,
	}

	h.e.Use(Cors())

	NewAuthHandler(h.e, conf.DB, conf.JwtUtil)
	NewMessageHandler(h.e, conf.DB, conf.JwtUtil, h.ws)

	h.e.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	h.e.GET("/ws", h.ws.Handler())

	return h
}

func (h *Handler) Run(address string) error {
	srv := &http.Server{
		Addr:    address,
		Handler: h.e,
	}

	go h.ws.Listen()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("ERROR:", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	return srv.Shutdown(context.Background())
}

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
