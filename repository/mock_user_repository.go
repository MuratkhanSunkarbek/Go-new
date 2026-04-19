
package repository

type MockUserRepository struct {
    GetUserByIDFunc func(id int) (*User, error)
    CreateUserFunc func(user *User) error
    GetByEmailFunc func(email string) (*User, error)
    UpdateUserFunc func(user *User) error
    DeleteUserFunc func(id int) error
}

func (m *MockUserRepository) GetUserByID(id int) (*User, error) {
    return m.GetUserByIDFunc(id)
}

func (m *MockUserRepository) CreateUser(user *User) error {
    return m.CreateUserFunc(user)
}

func (m *MockUserRepository) GetByEmail(email string) (*User, error) {
    return m.GetByEmailFunc(email)
}

func (m *MockUserRepository) UpdateUser(user *User) error {
    return m.UpdateUserFunc(user)
}

func (m *MockUserRepository) DeleteUser(id int) error {
    return m.DeleteUserFunc(id)
}
