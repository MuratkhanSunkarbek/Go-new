package service

import (
	"errors"
	"fmt"
	"practice8/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) RegisterUser(user *repository.User, email string) error {
	existing, err := s.repo.GetByEmail(email)

	if existing != nil {
		return fmt.Errorf("user exists")
	}

	if err != nil {
		return err
	}

	return s.repo.CreateUser(user)
}

func (s *UserService) UpdateUserName(id int, name string) error {
	if name == "" {
		return fmt.Errorf("empty")
	}

	u, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	u.Name = name
	return s.repo.UpdateUser(u)
}

func (s *UserService) DeleteUser(id int) error {
	if id == 1 {
		return errors.New("admin")
	}

	return s.repo.DeleteUser(id)
}