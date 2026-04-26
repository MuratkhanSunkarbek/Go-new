package service

import (
	"errors"
	"testing"

	"practice8/repository"
)



func TestRegisterUser_UserExists(t *testing.T) {
	mock := &repository.MockUserRepository{
		GetByEmailFunc: func(email string) (*repository.User, error) {
			return &repository.User{}, nil
		},
	}

	s := NewUserService(mock)
	err := s.RegisterUser(&repository.User{}, "a")

	if err == nil {
		t.Errorf("expected error when user exists")
	}
}

func TestRegisterUser_Success(t *testing.T) {
	mock := &repository.MockUserRepository{
		GetByEmailFunc: func(email string) (*repository.User, error) {
			return nil, nil
		},
		CreateUserFunc: func(user *repository.User) error {
			return nil
		},
	}

	s := NewUserService(mock)
	err := s.RegisterUser(&repository.User{}, "a")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegisterUser_RepoError(t *testing.T) {
	mock := &repository.MockUserRepository{
		GetByEmailFunc: func(email string) (*repository.User, error) {
			return nil, errors.New("db error")
		},
	}

	s := NewUserService(mock)
	err := s.RegisterUser(&repository.User{}, "a")

	if err == nil {
		t.Errorf("expected repo error")
	}
}



func TestUpdateUserName_Empty(t *testing.T) {
	mock := &repository.MockUserRepository{}

	s := NewUserService(mock)
	err := s.UpdateUserName(1, "")

	if err == nil {
		t.Errorf("expected error for empty name")
	}
}

func TestUpdateUserName_GetError(t *testing.T) {
	mock := &repository.MockUserRepository{
		GetUserByIDFunc: func(id int) (*repository.User, error) {
			return nil, errors.New("not found")
		},
	}

	s := NewUserService(mock)
	err := s.UpdateUserName(1, "test")

	if err == nil {
		t.Errorf("expected error")
	}
}

func TestUpdateUserName_Success(t *testing.T) {
	mock := &repository.MockUserRepository{
		GetUserByIDFunc: func(id int) (*repository.User, error) {
			return &repository.User{ID: id, Name: "old"}, nil
		},
		UpdateUserFunc: func(user *repository.User) error {
			if user.Name != "new" {
				t.Errorf("name not updated")
			}
			return nil
		},
	}

	s := NewUserService(mock)
	err := s.UpdateUserName(1, "new")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUpdateUserName_UpdateFails(t *testing.T) {
	mock := &repository.MockUserRepository{
		GetUserByIDFunc: func(id int) (*repository.User, error) {
			return &repository.User{ID: id, Name: "old"}, nil
		},
		UpdateUserFunc: func(user *repository.User) error {
			return errors.New("update failed")
		},
	}

	s := NewUserService(mock)
	err := s.UpdateUserName(1, "new")

	if err == nil {
		t.Errorf("expected error")
	}
}



func TestDeleteUser_Admin(t *testing.T) {
	mock := &repository.MockUserRepository{}

	s := NewUserService(mock)
	err := s.DeleteUser(1)

	if err == nil {
		t.Errorf("should not delete admin")
	}
}

func TestDeleteUser_Success(t *testing.T) {
	mock := &repository.MockUserRepository{
		DeleteUserFunc: func(id int) error {
			return nil
		},
	}

	s := NewUserService(mock)
	err := s.DeleteUser(2)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDeleteUser_Error(t *testing.T) {
	mock := &repository.MockUserRepository{
		DeleteUserFunc: func(id int) error {
			return errors.New("db error")
		},
	}

	s := NewUserService(mock)
	err := s.DeleteUser(2)

	if err == nil {
		t.Errorf("expected error")
	}
}