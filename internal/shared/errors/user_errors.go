package errors

import "errors"

var (
	ErrUserNotFound   = errors.New("usuário não encontrado")
	ErrUsernameLength = errors.New("o username deve ter pelo menos 3 caracteres")
	ErrPasswordLength = errors.New("a senha deve ter pelo menos 8 caracteres")
)
