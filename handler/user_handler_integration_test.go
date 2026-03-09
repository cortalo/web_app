// handler/user_handler_integration_test.go
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"web_app/domain"
	"web_app/infrastructure"
	"web_app/pkg/snowflake"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ===== 测试套件 =====

type UserHandlerIntegrationSuite struct {
	suite.Suite
	db      *sqlx.DB
	router  *gin.Engine
	handler *UserHandler
}

// 整个套件跑一次
func (s *UserHandlerIntegrationSuite) SetupSuite() {
	// 连接测试数据库
	db, err := sqlx.Connect("mysql", "root:root@tcp(localhost:3306)/bluebell_test?charset=utf8mb4&parseTime=True")
	s.Require().NoError(err)
	s.db = db

	// 初始化 snowflake
	err = snowflake.Init("2020-01-01", 1)
	s.Require().NoError(err)

	// 组装
	snowflakeNode := snowflake.NewSnowflakeNode()
	userRepo := infrastructure.NewMySQLUserRepository(db)
	userService := domain.NewUserService(userRepo, snowflakeNode)
	s.handler = NewUserHandler(userService)

	// 起 gin
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/users", s.handler.Register)
	s.router = r
}

// 每个测试前清空数据
func (s *UserHandlerIntegrationSuite) SetupTest() {
	s.db.Exec("DELETE FROM user")
}

// 整个套件结束后关闭 DB
func (s *UserHandlerIntegrationSuite) TearDownSuite() {
	s.db.Close()
}

// ===== 测试 =====

func (s *UserHandlerIntegrationSuite) TestRegister_Success() {
	body := ParamRegister{
		Username:   "alice",
		Password:   "password123",
		RePassword: "password123",
	}

	w := s.sendRequest(body)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	// 验证数据库里真的写进去了
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", "alice").Scan(&count)
	assert.Equal(s.T(), 1, count)
}

func (s *UserHandlerIntegrationSuite) TestRegister_DuplicateUsername() {
	body := ParamRegister{
		Username:   "alice",
		Password:   "password123",
		RePassword: "password123",
	}

	// 第一次注册
	s.sendRequest(body)

	// 第二次注册同名用户
	w := s.sendRequest(body)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *UserHandlerIntegrationSuite) TestRegister_InvalidParam_MissingUsername() {
	body := map[string]string{
		"password":    "password123",
		"re_password": "password123",
		// username 缺失
	}

	w := s.sendRequest(body)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *UserHandlerIntegrationSuite) TestRegister_InvalidParam_PasswordMismatch() {
	body := ParamRegister{
		Username:   "alice",
		Password:   "password123",
		RePassword: "different", // 和 password 不一致
	}

	w := s.sendRequest(body)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

// ===== 工具方法 =====

func (s *UserHandlerIntegrationSuite) sendRequest(body any) *httptest.ResponseRecorder {
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	return w
}

// ===== 入口 =====

func TestUserHandlerIntegration(t *testing.T) {
	suite.Run(t, new(UserHandlerIntegrationSuite))
}
