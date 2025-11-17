package http

import (
	"log"
	"net/http"

	"github.com/QuatroQuatros/go-real-time-chat/internal/dto"
	usecase "github.com/QuatroQuatros/go-real-time-chat/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *usecase.AuthUsecase
}

func NewAuthHandler(auth *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.SignUpDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.auth.Register(req)
	if err != nil {
		log.Println("Erro de registro:", err)
		status, message := MapDomainError(err)
		c.JSON(status, gin.H{"error": message})
		return
	}

	userPresenter := dto.NewUserPresenter(user)

	response := dto.AuthPresenter{
		User:  userPresenter,
		Token: token,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		log.Println("Erro de registro:", err)
		status, message := MapDomainError(err)
		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}
