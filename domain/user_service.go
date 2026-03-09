package domain

// UserService 接口，handler 依赖这个
type UserService interface {
	Register(name, email string) (User, error)
	GetUser(id string) (User, error)
}

// UserRepository 接口
type UserRepository interface {
	Save(user User) error
	FindByID(id string) (User, error)
	ExistsByEmail(email string) (bool, error)
}

// userService 小写，外部不能直接用具体类型
type userService struct {
	userRepo UserRepository
}

// 返回接口，隐藏具体实现
func NewUserService(userRepo UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// 方法先空着
func (s *userService) Register(name, email string) (User, error) {
	return User{}, nil
}

func (s *userService) GetUser(id string) (User, error) {
	return User{}, nil
}
