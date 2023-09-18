package api

import (
    "context"
    "fmt"
    "math/rand"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
    sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
    "go.uber.org/zap"

    "github.com/xiangnan0811/mxshop-api/user-web/forms"
    "github.com/xiangnan0811/mxshop-api/user-web/global"
)

func GenerateSmsCode(witdh int) string {
    numeric := [10]byte{0,1,2,3,4,5,6,7,8,9}
    r := len(numeric)
    rand.New(rand.NewSource(time.Now().UnixNano()))

    var sb strings.Builder
    for i := 0; i < witdh; i++ {
        fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
    }
    return sb.String()
}

func SendSms(c *gin.Context) {
    // register form
    smsForm := forms.SendSmsForm{}
    if err := c.ShouldBind(&smsForm); err != nil {
        zap.S().Errorw("发送验证码失败", "msg", err.Error())
        HandleValidateError(c, err)
        return
    }

    // 实例化一个认证对象，入参需要传入腾讯云账户密钥对secretId，secretKey
    credential := common.NewCredential(
        global.ServerConfig.TencentSmsInfo.SecretId,
        global.ServerConfig.TencentSmsInfo.SecretKey,
    )
    // 实例化一个客户端配置对象
    cpf := profile.NewClientProfile()
    // 实例化 sms 的 client 对象
    client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)
    // 实例化 sms 发送请求对象
    request := sms.NewSendSmsRequest()
    // 填充请求参数
    request.PhoneNumberSet = common.StringPtrs([]string{smsForm.Mobile})
    request.SmsSdkAppId = common.StringPtr(global.ServerConfig.TencentSmsInfo.SmsSdkAppId)
    request.TemplateId = common.StringPtr(global.ServerConfig.TencentSmsInfo.TemplateId)
    request.SignName = common.StringPtr(global.ServerConfig.TencentSmsInfo.SignName)

    // 生成验证码
    code := GenerateSmsCode(global.ServerConfig.SmsInfo.Length)
    // 验证码过期时间 默认 5 分钟
    codeExpire := strconv.Itoa(global.ServerConfig.SmsInfo.Expires)
    request.TemplateParamSet = common.StringPtrs([]string{code, codeExpire})
    response, err := client.SendSms(request)
    if _, ok := err.(*errors.TencentCloudSDKError); ok {
        zap.S().Errorf("An API error has returned: %s", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "msg": "发送失败",
        })
        return
    }
    if err != nil {
        zap.S().Errorf("Unknown error has returned: %s", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "msg": "发送失败",
        })
        return
    }

    // 将验证码存入 redis
    rdb := redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
        Password: global.ServerConfig.RedisInfo.Password,
        DB: global.ServerConfig.SmsInfo.SmsRedisDB,
    })
    rdb.Set(context.Background(), smsForm.Mobile, code, time.Duration(global.ServerConfig.SmsInfo.Expires))
    // 发送成功
    zap.S().Infof("response: %s", response.ToJsonString())
    c.JSON(http.StatusOK, gin.H{
        "msg": "发送成功",
    })
    c.Next()
}
