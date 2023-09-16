package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

// captcha store
var store = base64Captcha.DefaultMemStore

func NewCaptcha(c *gin.Context) {
	captcha := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, store)
	id, b64s, err := captcha.Generate()
	if err != nil {
		zap.S().Errorf("generate captcha error: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})

}
