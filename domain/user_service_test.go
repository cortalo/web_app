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
	findByIDFunc         func(id string) (User, error)
}

func (m *mockUserRepository) ExistsByUsername(username string) (bool, error) {
	return m.existsByUsernameFunc(username)
}

func (m *mockUserRepository) Save(user *User) error {
	return m.saveFunc(user)
}

func (m *mockUserRepository) FindByID(id string) (User, error) {
	return m.findByIDFunc(id)
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
