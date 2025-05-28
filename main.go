package main

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// 中文大写数字和单位
var (
	digits = []string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
	units  = []string{"", "拾", "佰", "仟", "万", "拾", "佰", "仟", "亿"}
)

// ConvertToChinese 金额转大写
func ConvertToChinese(amount string) (string, error) {
	// 验证金额格式
	re := regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	if !re.MatchString(amount) {
		return "", http.ErrBodyNotAllowed
	}

	// 分割整数和小数部分
	var integer, decimal string
	if strings.Contains(amount, ".") {
		parts := strings.Split(amount, ".")
		integer = parts[0]
		decimal = parts[1]
		// 补零到两位小数
		if len(decimal) == 1 {
			decimal += "0"
		} else if len(decimal) > 2 {
			decimal = decimal[:2]
		}
	} else {
		integer = amount
		decimal = "00"
	}

	// 处理整数部分
	var result strings.Builder
	integerLen := len(integer)
	if integerLen > 12 { // 最大支持到千亿
		return "", http.ErrBodyNotAllowed
	}

	for i, r := range integer {
		digit := int(r - '0')
		pos := integerLen - i - 1

		if digit != 0 {
			result.WriteString(digits[digit])
			result.WriteString(units[pos])
		} else {
			// 处理零的情况
			if pos == 4 || pos == 8 { // 万位和亿位
				result.WriteString(units[pos])
			}
			// 避免连续的零
			if i < integerLen-1 && int(integer[i+1]-'0') != 0 && !strings.HasSuffix(result.String(), "零") {
				result.WriteString("零")
			}
		}
	}

	// 处理小数部分
	if decimal != "00" {
		result.WriteString("点")
		for _, r := range decimal {
			result.WriteString(digits[int(r-'0')])
		}
	} else {
		result.WriteString("整")
	}

	return result.String(), nil
}

func main() {
	r := gin.Default()

	// 定义静态文件目录
	r.Static("/static", "./static")

	// API 路由
	r.POST("/api/convert", func(c *gin.Context) {
		var request struct {
			Amount string `json:"amount"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		chinese, err := ConvertToChinese(request.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"chinese": chinese})
	})

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// 启动服务器
	r.Run(":8080")
}