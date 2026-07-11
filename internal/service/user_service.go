package service

import (
	"errors"
	"nailly-back-end/internal/model"
	"nailly-back-end/internal/repository"
	"nailly-back-end/pkg/utils"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers(filter repository.UserFilter, pagination utils.Pagination) ([]model.User, int64, error) {
	return s.repo.FindAll(filter, pagination)
}

func (s *UserService) GetUserByID(id string) (model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) GetUserByEmail(email string) (model.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *UserService) GetUsersOlderThan(age int) ([]model.User, error) {
	return s.repo.FindOlderThan(age)
}

func (s *UserService) CreateUser(input model.User) (model.User, error) {
	if input.Name == "" {
		return model.User{}, errors.New("name is required")
	}
	if input.Email == "" {
		return model.User{}, errors.New("email is required")
	}
	if input.Age <= 0 {
		return model.User{}, errors.New("age is required and must be greater than 0")
	}

	if err := s.repo.Create(&input); err != nil {
		return model.User{}, err
	}

	return input, nil
}

func (s *UserService) UpdateUser(id string, input model.User) (model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return model.User{}, err
	}

	if err := s.repo.Update(&user, input); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(&user)
}
