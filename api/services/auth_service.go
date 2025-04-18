package services

import (
	"errors"

	"github.com/luxifa/beauty/config"
	"github.com/luxifa/beauty/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 提供身份验证相关的服务
type AuthService struct{}

// Register 注册新用户
func (s *AuthService) Register(username, password, email string) (*models.User, error) {
    // 检查用户是否已存在
    var existingUser models.User
    if err := config.DB.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
        return nil, errors.New("用户名或电子邮箱已被使用")
    }

    // 对密码进行哈希处理
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // 创建新用户
    newUser := models.User{
        Username: username,
        Password: string(hashedPassword),
        Email:    email,
    }

    if err := config.DB.Create(&newUser).Error; err != nil {
        return nil, err
    }

    // 不返回密码
    newUser.Password = ""
    return &newUser, nil
}

// Login 用户登录
func (s *AuthService) Login(username, password string) (*models.User, error) {
    var user models.User
    
    // 查找用户
    if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
        return nil, errors.New("用户不存在")
    }

    // 验证密码
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("密码错误")
    }

    // 不返回密码
    user.Password = ""
    return &user, nil
}

// GetUserByID 通过ID获取用户信息
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
    var user models.User
    
    if err := config.DB.First(&user, id).Error; err != nil {
        return nil, errors.New("用户不存在")
    }
    
    // 不返回密码
    user.Password = ""
    return &user, nil
}