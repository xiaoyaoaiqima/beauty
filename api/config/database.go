package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/luxifa/beauty/models"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
    // 加载.env文件
    err := godotenv.Load()
    if err != nil {
        log.Println("未找到.env文件，使用环境变量")
    }

    // 获取数据库配置信息
    dbUser := getEnv("DB_USER", "root")
    dbPass := getEnv("DB_PASSWORD", "")
    dbHost := getEnv("DB_HOST", "localhost")
    dbPort := getEnv("DB_PORT", "3306")
    dbName := getEnv("DB_NAME", "beauty")

    // 首先尝试连接MySQL服务器（不指定数据库）
    rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/", 
        dbUser, dbPass, dbHost, dbPort)
    
    rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{})
    if err != nil {
        log.Fatalf("连接MySQL服务器失败: %v", err)
    }
    
    // 检查数据库是否存在，如果不存在则创建
    var count int64
    rootDB.Raw(fmt.Sprintf("SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = '%s'", dbName)).Scan(&count)
    
    if count == 0 {
        log.Printf("数据库 %s 不存在，正在创建...", dbName)
        if err := rootDB.Exec(fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", dbName)).Error; err != nil {
            log.Fatalf("创建数据库失败: %v", err)
        }
        log.Printf("数据库 %s 已成功创建", dbName)
    }

    // 构建完整DSN连接到指定数据库
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        dbUser, dbPass, dbHost, dbPort, dbName)

    // 连接到特定数据库
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("连接数据库失败: %v", err)
    }

    // 检查用户表是否存在，如果不存在则创建
    if !db.Migrator().HasTable(&models.User{}) {
        log.Println("用户表不存在，正在创建...")
        if err := db.AutoMigrate(&models.User{}); err != nil {
            log.Fatalf("创建用户表失败: %v", err)
        }
        log.Println("用户表创建成功")
        
        // 可以在这里添加一些初始用户数据
        // 例如：创建管理员账户
        // createAdminUser(db)
    } else {
        log.Println("用户表已存在")
    }

    DB = db
    log.Println("数据库连接成功")
}

// 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

// 可选：创建管理员账户
/*
func createAdminUser(db *gorm.DB) {
    // 这里可以添加代码创建初始管理员用户
    // 记得对密码进行哈希处理
}
*/