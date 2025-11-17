package usecase

import (
	"github.com/QuatroQuatros/go-real-time-chat/infra/repository"
	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
	"github.com/QuatroQuatros/go-real-time-chat/internal/dto"
	"github.com/QuatroQuatros/go-real-time-chat/internal/token"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	repo repository.UserRepository
}

func NewAuthUsecase(repo repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func (u *AuthUsecase) Register(req dto.SignUpDTO) (string, *domain.User, error) {
	var existingUser *domain.User
	var err error

	if req.Username != "" && req.Password != "" {
		existingUser, err = u.repo.GetByUsername(req.Username)
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return "", nil, err
	}

	if existingUser != nil {
		jwt, err := token.GenerateJWT(existingUser.ID, existingUser.Guest)
		if err != nil {
			return "", nil, err
		}
		return jwt, existingUser, nil
	}

	user, err := domain.NewUserFomInput(
		req.Username,
		req.Password,
		req.Guest,
	)

	if err := u.repo.Create(user); err != nil {
		return "", nil, err
	}

	jwt, err := token.GenerateJWT(user.ID, user.Guest)
	return jwt, user, nil

}

func (u *AuthUsecase) Login(email, password string) (string, *domain.User, error) {
	// Lógica de login de usuário
	return "", nil, nil
}
