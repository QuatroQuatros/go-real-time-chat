package http

import (
	"net/http"

	userErrors "github.com/QuatroQuatros/go-real-time-chat/internal/shared/errors"
)

// MapDomainError mapeia erros de domínio para códigos de status HTTP.
func MapDomainError(err error) (int, string) {
	switch err {
	// Erros de validação
	case userErrors.ErrUsernameLength, userErrors.ErrPasswordLength:
		return http.StatusUnprocessableEntity, err.Error()
	default:
		// Para erros inesperados, retorna um erro genérico.
		return http.StatusInternalServerError, "um erro inesperado ocorreu"
	}
}
