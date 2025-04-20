package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/luxifa/beauty/config"
	"github.com/luxifa/beauty/models"
)

// ExportService 提供数据导出相关的服务
type ExportService struct{}

// ExportFormat 导出格式
type ExportFormat string

const (
    FormatCSV  ExportFormat = "csv"
    FormatJSON ExportFormat = "json"
)

// ExportUsers 导出用户数据
func (s *ExportService) ExportUsers(format ExportFormat) (string, error) {
    var users []models.User

    // 从数据库获取所有用户
    if err := config.DB.Find(&users).Error; err != nil {
        return "", fmt.Errorf("获取用户数据失败: %w", err)
    }

    // 创建临时目录，如果不存在
    tempDir := "./temp"
    if _, err := os.Stat(tempDir); os.IsNotExist(err) {
        if err := os.Mkdir(tempDir, 0755); err != nil {
            return "", fmt.Errorf("创建临时目录失败: %w", err)
        }
    }

    // 创建唯一文件名
    timestamp := time.Now().Format("20060102_150405")
    var filePath string
    var err error

    switch format {
    case FormatCSV:
        filePath = fmt.Sprintf("%s/users_%s.csv", tempDir, timestamp)
        err = exportToCSV(users, filePath)
    case FormatJSON:
        filePath = fmt.Sprintf("%s/users_%s.json", tempDir, timestamp)
        err = exportToJSON(users, filePath)
    default:
        return "", fmt.Errorf("不支持的导出格式: %s", format)
    }

    if err != nil {
        return "", err
    }

    return filePath, nil
}

// exportToCSV 导出为CSV格式
func exportToCSV(users []models.User, filePath string) error {
    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("创建CSV文件失败: %w", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // 写入CSV头部
    headers := []string{"ID", "用户名", "邮箱", "创建时间", "更新时间"}
    if err := writer.Write(headers); err != nil {
        return fmt.Errorf("写入CSV头部失败: %w", err)
    }

    // 写入数据行
    for _, user := range users {
        row := []string{
            fmt.Sprintf("%d", user.ID),
            user.Username,
            user.Email,
            user.CreatedAt.Format("2006-01-02 15:04:05"),
            user.UpdatedAt.Format("2006-01-02 15:04:05"),
        }
        if err := writer.Write(row); err != nil {
            return fmt.Errorf("写入CSV数据行失败: %w", err)
        }
    }

    return nil
}

// exportToJSON 导出为JSON格式
func exportToJSON(users []models.User, filePath string) error {
    file, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("创建JSON文件失败: %w", err)
    }
    defer file.Close()

    // 清理敏感字段
    exportUsers := make([]map[string]interface{}, len(users))
    for i, user := range users {
        // 不包含密码字段
        exportUsers[i] = map[string]interface{}{
            "id":         user.ID,
            "username":   user.Username,
            "email":      user.Email,
            "created_at": user.CreatedAt,
            "updated_at": user.UpdatedAt,
        }
    }

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ") // 美化JSON输出
    if err := encoder.Encode(exportUsers); err != nil {
        return fmt.Errorf("写入JSON数据失败: %w", err)
    }

    return nil
}