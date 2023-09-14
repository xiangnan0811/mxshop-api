package initialize

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"

	"github.com/xiangnan0811/mxshop-api/user-web/global"
)

func InitTrans(locale string) (err error) {
    // 修改 gin 的 Validator 引擎属性，实现自定义翻译器
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        // 注册一个获取 json tag 的自定义方法
        v.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
            if name == "-" {
                return ""
            }
            return name
        })
        zhT := zh.New() // 中文翻译器
        enT := en.New() // 英文翻译器
        // 第一个参数是备用（fallback）的语言环境, 后面参数是应该支持的语言环境（支持多个）
        uni := ut.New(enT, enT, zhT)
        global.Trans, ok = uni.GetTranslator(locale)
        if !ok {
            zap.S().Errorf("uni.GetTranslator(%s)", locale)
            return fmt.Errorf("uni.GetTranslator(%s)", locale)
        }
        switch locale {
        case "en":
            en_translations.RegisterDefaultTranslations(v, global.Trans)
        case "zh":
            zh_translations.RegisterDefaultTranslations(v, global.Trans)
        default:
            en_translations.RegisterDefaultTranslations(v, global.Trans)
        }
        return
    }
    return
}