package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luxifa/beauty/services"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{
        authService: &services.AuthService{},
    }
}

// 注册请求
type RegisterRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
}

// 登录请求
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// Register 处理用户注册
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    user, err := h.authService.Register(req.Username, req.Password, req.Email)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "注册成功",
        "user":    user,
    })
}

// Login 处理用户登录
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
        return
    }

    user, err := h.authService.Login(req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // 这里可以生成JWT token，但为了保持简洁，我们直接返回用户信息
    c.JSON(http.StatusOK, gin.H{
        "message": "登录成功",
        "user":    user,
    })
}

// GetUserInfo 获取用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
    // 这里应该从JWT中获取用户ID，但为简洁起见，我们从URL参数中获取
    userID := c.Param("id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不能为空"})
        return
    }

    var id uint
    _, err := fmt.Sscanf(userID, "%d", &id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
        return
    }

    user, err := h.authService.GetUserByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"user": user})
}