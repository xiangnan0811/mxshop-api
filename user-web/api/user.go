package api

import (
    "context"
    "fmt"
    "net/http"
    "strconv"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "go.uber.org/zap"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/status"

    "github.com/xiangnan0811/mxshop-api/user-web/forms"
    "github.com/xiangnan0811/mxshop-api/user-web/global"
    "github.com/xiangnan0811/mxshop-api/user-web/global/response"
    "github.com/xiangnan0811/mxshop-api/user-web/middlewares"
    "github.com/xiangnan0811/mxshop-api/user-web/models"
    "github.com/xiangnan0811/mxshop-api/user-web/proto"
    "github.com/xiangnan0811/mxshop-api/user-web/utils"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
    // 将grpc的code转换为http的状态码
    if e, ok := status.FromError(err); ok {
        switch e.Code() {
        case codes.NotFound:
            c.JSON(http.StatusNotFound, gin.H{
                "msg": e.Message(),
            })
        case codes.Internal:
            c.JSON(http.StatusInternalServerError, gin.H{
                "msg": "内部错误" + e.Message(),
            })
        case codes.InvalidArgument:
            c.JSON(http.StatusBadRequest, gin.H{
                "msg": "参数错误",
            })
        case codes.Unavailable:
            c.JSON(http.StatusInternalServerError, gin.H{
                "msg": "服务不可用",
            })
        default:
            c.JSON(http.StatusInternalServerError, gin.H{
                "msg": "其他错误",
            })
            return
        }
    }
}

func GetUserList(ctx *gin.Context) {
    // 拨号连接用户grpc服务器
    userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConfig.Host,
        global.ServerConfig.UserSrvConfig.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        zap.S().Errorw(
            "[GetUserList] 连接 【用户服务】 失败",
            "msg", err.Error(),
        )
    }
    // 生成grpc的client并调用接口
    userSrvClient := proto.NewUserClient(userConn)

    pn := ctx.DefaultQuery("pn", "0")
    pnInt, _ := strconv.Atoi(pn)
    pSize := ctx.DefaultQuery("psize", "10")
    pSizeInt, _ := strconv.Atoi(pSize)

    rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
        Pn:    uint32(pnInt),
        PSize: uint32(pSizeInt),
    })
    if err != nil {
        zap.S().Errorw("[GetUserList] 查询 【用户列表】 失败")
        HandleGrpcErrorToHttp(err, ctx)
        return
    }
    result := make([]interface{}, 0)
    for _, value := range rsp.Data {

        user := response.UserResponse{
            Id:       value.Id,
            NickName: value.NickName,
            Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
            Gender:   value.Gender,
            Mobile:   value.Mobile,
        }

        result = append(result, user)
    }
    ctx.JSON(http.StatusOK, result)
}

func PassWordLogin(c *gin.Context) {
    // 参数校验
    passwordLoginForm := forms.PassWordLoginForm{}
    if err := c.ShouldBindJSON(&passwordLoginForm); err != nil {
        HandleValidateError(c, err)
        return
    }

    // 验证码校验
    ok := store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{
            "msg": "验证码错误",
        })
        return
    }

    // 拨号连接用户grpc服务器
    userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConfig.Host,
        global.ServerConfig.UserSrvConfig.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        zap.S().Errorw(
            "[GetUserList] 连接 【用户服务】 失败",
            "msg", err.Error(),
        )
    }
    // 生成grpc的client并调用接口
    userSrvClient := proto.NewUserClient(userConn)

    // 登录逻辑
    userRsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
        Mobile: passwordLoginForm.Mobile,
    })
    if err != nil {
        zap.S().Infoln("[PassWordLogin] 查询 【用户】 失败")
        if e, ok := status.FromError(err); ok {
            switch e.Code() {
            case codes.NotFound:
                c.JSON(http.StatusBadRequest, gin.H{
                    "msg": "用户不存在",
                })
            default:
                c.JSON(http.StatusInternalServerError, gin.H{
                    "msg": "内部错误",
                })
            }
        }
        return
    }
    // 校验密码
    passRsp, passErr := userSrvClient.CheckPassWord(
        context.Background(), 
        &proto.PassWordCheckRequest{
            Password: passwordLoginForm.Password,
            EncryptedPassword: userRsp.PassWord,
    })
    if passErr != nil {
        zap.S().Errorw("[PassWordLogin] 查询 【用户】 失败")
        c.JSON(http.StatusInternalServerError, gin.H{
            "msg": "登录失败",
        })
        return
    }
    if  !passRsp.Success {
        c.JSON(http.StatusBadRequest, gin.H{
            "msg": "密码错误",
        })
        return
    }
    // 生成token
    j := middlewares.NewJWT()
    claims := models.CustomClaims{
        ID:          uint(userRsp.Id),
        NickName:    userRsp.NickName,
        AuthorityId: uint(userRsp.Role),
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Duration(global.ServerConfig.JWTInfo.Expires) * time.Second).Unix(),
            Issuer:    global.ServerConfig.JWTInfo.Issuer,
        },
    }
    token, err := j.CreateToken(claims)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "msg": "内部错误",
        })
        return
    }

    // 返回结果
    c.JSON(http.StatusOK, gin.H{
        "id": userRsp.Id,
        "nick_name": userRsp.NickName,
        "token": token,
        "expired_at":  int64(claims.ExpiresAt) * 1000,
    })
}

func HandleValidateError(c *gin.Context, err error) {
    // 处理参数校验错误
    errs, ok := err.(validator.ValidationErrors)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    c.JSON(http.StatusBadRequest, gin.H{
        "error": utils.RemoveTopStruct(errs.Translate(global.Trans)),
    })
}
