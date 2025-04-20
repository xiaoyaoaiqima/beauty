package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luxifa/beauty/services"
)

type ExportHandler struct {
    exportService *services.ExportService
}

func NewExportHandler() *ExportHandler {
    return &ExportHandler{
        exportService: &services.ExportService{},
    }
}

// ExportUsers 导出用户数据
func (h *ExportHandler) ExportUsers(c *gin.Context) {
    // 获取格式参数，默认为CSV
    format := c.DefaultQuery("format", "csv")
    var exportFormat services.ExportFormat

    switch format {
    case "json":
        exportFormat = services.FormatJSON
    default:
        exportFormat = services.FormatCSV
    }

    // 调用服务导出数据
    filePath, err := h.exportService.ExportUsers(exportFormat)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "导出数据失败: " + err.Error()})
        return
    }

    // 获取文件名
    fileName := filepath.Base(filePath)

    // 设置响应头，通知浏览器这是一个需要下载的文件
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename="+fileName)

    // 根据文件类型设置Content-Type
    if exportFormat == services.FormatCSV {
        c.Header("Content-Type", "text/csv")
    } else {
        c.Header("Content-Type", "application/json")
    }

    // 提供文件下载
    c.File(filePath)

    // 文件发送后删除临时文件
    go func() {
        // 等待一段时间以确保文件已被发送
        time.Sleep(5 * time.Second)
        os.Remove(filePath)
    }()
}