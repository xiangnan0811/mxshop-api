package middlewares

import (
    "errors"
    "net/http"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "github.com/xiangnan0811/mxshop-api/user-web/global"
    "github.com/xiangnan0811/mxshop-api/user-web/models"
)

var (
    TokenMalformed          = errors.New("Token is malformed")
    TokenUnverifiable       = errors.New("Token could not be verified because of signing problems")
    TokenSignatureInvalid   = errors.New("Token signature is invalid")
    TokenErrorAudience      = errors.New("Token audience mismatch")
    TokenErrorExpired       = errors.New("Token is expired")
    TokenErrorIssuedAt      = errors.New("Token time is incorrect")
    TokenErrorIssuer        = errors.New("Token issuer mismatch")
    TokenErrorNotValidYet   = errors.New("Token not active yet")
    TokenErrorId            = errors.New("Token id is invalid")
    TokenErrorClaimsInvalid = errors.New("Token is invalid")
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
        token := c.Request.Header.Get("x-token")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "msg": "请登录",
            })
            c.Abort()
            return
        }
        j := NewJWT()
        // parseToken 解析token包含的信息
        claims, err := j.ParseToken(token)
        if err != nil {
            switch err {
            case TokenErrorExpired:
                c.JSON(http.StatusUnauthorized, gin.H{
                    "msg": "token过期",
                })
                c.Abort()
            default:
                zap.S().Infoln("非法 token: ", err.Error())
                c.JSON(http.StatusUnauthorized, gin.H{
                    "msg": "非法token",
                })
                c.Abort()
            }
            return
        }
        c.Set("claims", claims)
        c.Set("userId", claims.ID)
        c.Next()
    }
}

type JWT struct {
    SigningKey []byte
    ExpiresAt  int64
    Issuer     string
}

func NewJWT() *JWT {
    return &JWT{
        []byte(global.ServerConfig.JWTInfo.SigningKey),
        time.Now().Add(time.Duration(global.ServerConfig.JWTInfo.Expires) * time.Second).Unix(),
        global.ServerConfig.JWTInfo.Issuer,
    }
}

// 创建一个token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.SigningKey)
}

// 自定义检验
func (j *JWT) customClaimsValid(t *jwt.Token) (i interface{}, e error) {
    claims, ok := t.Claims.(*models.CustomClaims)
    if !ok {
        zap.S().Infof("customClaimsValid: %v", ok)
        return nil, jwt.NewValidationError(TokenMalformed.Error(), jwt.ValidationErrorMalformed)
    }
    if issOk := claims.VerifyIssuer(global.ServerConfig.JWTInfo.Issuer, true); !issOk {
        return nil, jwt.NewValidationError(TokenErrorIssuer.Error(), jwt.ValidationErrorIssuer)
    }
    return j.SigningKey, nil
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, j.customClaimsValid)
    if err != nil {
        if ve, ok := err.(*jwt.ValidationError); ok {
            if ve.Errors&jwt.ValidationErrorMalformed != 0 {
                return nil, TokenMalformed
            } else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
                return nil, TokenUnverifiable
            } else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
                return nil, TokenSignatureInvalid
            } else if ve.Errors&jwt.ValidationErrorAudience != 0 {
                return nil, TokenErrorAudience
            } else if ve.Errors&jwt.ValidationErrorExpired != 0 {
                return nil, TokenErrorExpired
            } else if ve.Errors&jwt.ValidationErrorIssuedAt != 0 {
                return nil, TokenErrorIssuedAt
            } else if ve.Errors&jwt.ValidationErrorIssuer != 0 {
                return nil, TokenErrorIssuer
            } else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
                return nil, TokenErrorNotValidYet
            } else if ve.Errors&jwt.ValidationErrorId != 0 {
                return nil, TokenErrorId
            } else {
                return nil, TokenErrorClaimsInvalid
            }
        }
    }
    if token != nil {
        if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
            return claims, nil
        }
    }
    return nil, TokenErrorClaimsInvalid
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
    jwt.TimeFunc = func() time.Time {
        return time.Unix(0, 0)
    }
    token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, j.customClaimsValid)
    if err != nil {
        return "", err
    }
    if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
        jwt.TimeFunc = time.Now
        claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
        return j.CreateToken(*claims)
    }
    return "", TokenErrorClaimsInvalid
}
