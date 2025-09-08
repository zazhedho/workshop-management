package service

import (
	"time"
	"workshop-management/internal/domain"
	"workshop-management/internal/dto"
	"workshop-management/internal/repository/auth"
	"workshop-management/internal/repository/users"
	"workshop-management/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo      users.Users
	BlacklistRepo auth.Blacklist
}

func NewUserService(userRepo users.Users, blacklistRepo auth.Blacklist) *UserService {
	return &UserService{
		UserRepo:      userRepo,
		BlacklistRepo: blacklistRepo,
	}
}

func (s *UserService) RegisterUser(req dto.UserRegister) (domain.Users, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.Users{}, err
	}

	user := domain.Users{
		Id:        utils.CreateUUID(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPwd),
		Role:      utils.RoleMember,
		CreatedAt: time.Now(),
	}

	if err = s.UserRepo.Store(user); err != nil {
		return domain.Users{}, err
	}

	return user, nil
}

func (s *UserService) LoginUser(req dto.Login, logId string) (string, error) {
	user, err := s.UserRepo.GetByEmail(req.Email)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", err
	}

	token, err := utils.GenerateJwt(&user, logId)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) LogoutUser(token string) error {
	blacklist := domain.Blacklist{
		ID:        utils.CreateUUID(),
		Token:     token,
		CreatedAt: time.Now(),
	}

	err := s.BlacklistRepo.Store(blacklist)
	if err != nil {
		return err
	}

	return nil
}
