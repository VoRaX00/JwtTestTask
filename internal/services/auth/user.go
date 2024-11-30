package auth

import (
	authrepo "JwtTestTask/internal/storage/auth"
	"JwtTestTask/models"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"os"
)

type UserService struct {
	repo authrepo.IUserRepository
}

func NewUserService(repo authrepo.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(user models.User) (string, error) {
	user.Id = uuid.New().String()
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.Create(user)
}

func (s *UserService) generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))
}
