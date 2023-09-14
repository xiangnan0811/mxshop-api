package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	"github.com/xiangnan0811/mxshop-api/user-web/proto"
)

func removeTopStruct(fields map[string]string) map[string]string {
    r := make(map[string]string, len(fields))
    for field, val := range fields {
        r[field[strings.Index(field, ".")+1:]] = val
    }
    return r
}

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
        errs, ok := err.(validator.ValidationErrors)
        if !ok {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{
            "error": removeTopStruct(errs.Translate(global.Trans)),
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "status": "you are logged in",
    })
}