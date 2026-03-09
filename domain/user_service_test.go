package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ===== Mock 实现 =====

type mockUserRepository struct {
	existsByUsernameFunc func(username string) (bool, error)
	saveFunc             func(user *User) error
	findByUsernameFunc   func(username string) (*User, error)
}

func (m *mockUserRepository) ExistsByUsername(username string) (bool, error) {
	return m.existsByUsernameFunc(username)
}

func (m *mockUserRepository) Save(user *User) error {
	return m.saveFunc(user)
}

func (m *mockUserRepository) FindByUsername(username string) (*User, error) {
	return m.findByUsernameFunc(username)
}

type mockSnowflake struct {
	id int64
}

func (m *mockSnowflake) NextID() int64 {
	return m.id
}

// ===== 测试 =====

func TestRegister_Success(t *testing.T) {
	repo := &mockUserRepository{
		existsByUsernameFunc: func(username string) (bool, error) {
			return false, nil // 用户不存在
		},
		saveFunc: func(user *User) error {
			return nil // 保存成功
		},
	}
	sf := &mockSnowflake{id: 123456}
	svc := NewUserService(repo, sf)

	err := svc.Register("alice", "password123")

	assert.NoError(t, err)
}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	repo := &mockUserRepository{
		existsByUsernameFunc: func(username string) (bool, error) {
			return true, nil // 用户已存在
		},
	}
	sf := &mockSnowflake{id: 123456}
	svc := NewUserService(repo, sf)

	err := svc.Register("alice", "password123")

	assert.Error(t, err)
	assert.EqualError(t, err, "username exist")
}

func TestRegister_ExistsByUsernameError(t *testing.T) {
	repo := &mockUserRepository{
		existsByUsernameFunc: func(username string) (bool, error) {
			return false, errors.New("db error") // 数据库报错
		},
	}
	sf := &mockSnowflake{id: 123456}
	svc := NewUserService(repo, sf)

	err := svc.Register("alice", "password123")

	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
}

func TestRegister_SaveError(t *testing.T) {
	repo := &mockUserRepository{
		existsByUsernameFunc: func(username string) (bool, error) {
			return false, nil
		},
		saveFunc: func(user *User) error {
			return errors.New("save failed") // 保存失败
		},
	}
	sf := &mockSnowflake{id: 123456}
	svc := NewUserService(repo, sf)

	err := svc.Register("alice", "password123")

	assert.Error(t, err)
	assert.EqualError(t, err, "save failed")
}

func TestRegister_UserFieldsCorrect(t *testing.T) {
	var savedUser *User

	repo := &mockUserRepository{
		existsByUsernameFunc: func(username string) (bool, error) {
			return false, nil
		},
		saveFunc: func(user *User) error {
			savedUser = user // 捕获保存的 user
			return nil
		},
	}
	sf := &mockSnowflake{id: 999}
	svc := NewUserService(repo, sf)

	err := svc.Register("alice", "password123")

	assert.NoError(t, err)
	assert.Equal(t, int64(999), savedUser.UserID)
	assert.Equal(t, "alice", savedUser.Username)
	assert.Equal(t, "password123", savedUser.Password)
}

func TestLogin_Success(t *testing.T) {
	var user = User{
		Username: "alice",
		Password: "password123",
	}
	repo := &mockUserRepository{
		findByUsernameFunc: func(username string) (*User, error) {
			return &user, nil
		},
	}
	sf := &mockSnowflake{id: 999}
	svc := NewUserService(repo, sf)

	err := svc.Login("alice", "password123")
	assert.NoError(t, err)
}

func TestLogin_DBError(t *testing.T) {
	repo := &mockUserRepository{
		findByUsernameFunc: func(username string) (*User, error) {
			return nil, errors.New("db error")
		},
	}
	sf := &mockSnowflake{id: 999}
	svc := NewUserService(repo, sf)

	err := svc.Login("alice", "password123")
	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
}

func TestLogin_WrongPassword(t *testing.T) {
	var user = User{
		Username: "alice",
		Password: "password123",
	}
	repo := &mockUserRepository{
		findByUsernameFunc: func(username string) (*User, error) {
			return &user, nil
		},
	}
	sf := &mockSnowflake{id: 999}
	svc := NewUserService(repo, sf)

	err := svc.Login("alice", "password")
	assert.Error(t, err)
	assert.EqualError(t, err, "wrong password")
}
