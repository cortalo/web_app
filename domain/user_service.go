package domain

import "go.uber.org/zap"

// UserService 接口，handler 依赖这个
type UserService interface {
	Register(name, email string) error
	Login(username string, password string) error
}

// UserRepository 接口
type UserRepository interface {
	Save(user *User) error
	ExistsByUsername(username string) (bool, error)
	FindByUsername(username string) (*User, error)
}

type Snowflake interface {
	NextID() int64
}

// userService 小写，外部不能直接用具体类型
type userService struct {
	userRepo  UserRepository
	snowflake Snowflake
}

// 返回接口，隐藏具体实现
func NewUserService(userRepo UserRepository, snowflake Snowflake) UserService {
	return &userService{
		userRepo:  userRepo,
		snowflake: snowflake,
	}
}

func (s *userService) Register(username, password string) error {
	exist, err := s.userRepo.ExistsByUsername(username)
	if err != nil {
		return err
	}
	if exist {
		return ErrUsernameExist
	}
	userId := s.snowflake.NextID()
	user := &User{
		UserID:   userId,
		Username: username,
		Password: password,
	}
	return s.userRepo.Save(user)
}

func (s *userService) Login(username string, password string) error {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		zap.L().Info("user not found", zap.String("username", username))
		return err
	}
	if password != user.Password {
		zap.L().Info("password is wrong", zap.String("username", username))
		return ErrWrongPassword
	}
	return nil
}
