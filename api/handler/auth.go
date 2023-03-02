package handler

import (
	"github.com/intrntsrfr/vue-ws-test/structs"
	"net/http"
	"time"

	"github.com/google/uuid"
	api "github.com/intrntsrfr/vue-ws-test"
	"github.com/intrntsrfr/vue-ws-test/database"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	r   *gin.Engine
	db  database.DB
	jwt api.JWTService
}

func NewAuthHandler(r *gin.Engine, db database.DB, jwtService api.JWTService) {
	h := &AuthHandler{r, db, jwtService}

	g := h.r.Group("/api/auth")
	g.POST("/login", h.login())
	g.POST("/register", h.register())
}

func (h *AuthHandler) login() gin.HandlerFunc {
	type LoginBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(c *gin.Context) {
		var loginBody LoginBody
		if err := c.BindJSON(&loginBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{CodeError, "bad request"})
			return
		}

		user := h.db.FindUserByUsername(loginBody.Username)
		if user == nil || loginBody.Password != user.Password {
			c.JSON(http.StatusUnauthorized, ErrorResponse{CodeError, "invalid username or password"})
			return
		}

		token, err := h.jwt.GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{CodeError, "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

func (h *AuthHandler) register() gin.HandlerFunc {
	type RegisterBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(c *gin.Context) {
		var registerBody RegisterBody
		if err := c.BindJSON(&registerBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{CodeError, "bad request"})
			return
		}

		user := h.db.FindUserByUsername(registerBody.Username)
		if user != nil {
			c.JSON(http.StatusConflict, ErrorResponse{CodeError, "username already taken"})
			return
		}

		user, err := h.db.CreateUser(&structs.User{
			ID:       uuid.New(),
			Username: registerBody.Username,
			Password: registerBody.Password,
			Created:  time.Now(),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{CodeError, "internal server error"})
			return
		}

		token, err := h.jwt.GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{CodeError, "internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
