package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	jwtToken "github.com/go-micro/plugins/v4/auth/jwt/token"
	"go-micro.dev/v4/auth"
	"strconv"
	"strings"
	"web/code"
	"web/pkg/response"
)

const (
	authorization = "Authorization"
)

func GinJwt(jwt jwtToken.Provider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorization)
		if !strings.HasPrefix(authHeader, auth.BearerScheme) {
			response.Response(ctx, code.NewCode(code.AuthFailed))
			return
		}
		token := strings.TrimPrefix(authHeader, auth.BearerScheme)
		if token == "" {
			response.Failed(ctx, errors.New("token为空"))
			ctx.Abort()
			return
		}
		l.Info("token信息:", token)
		account, err := jwt.Inspect(token)
		if err != nil {
			response.Failed(ctx, errors.Wrap(err, "token验证失败"))
			ctx.Abort()
			return
		}

		uid, _ := strconv.Atoi(account.ID)
		/*		_, err = api.CheckToken(ctx, &user.CheckTokenReq{
					Uids:   int32(uid),
					Token: token,
				})
				if err != nil {
					response.Failed(ctx, err)
					ctx.Abort()
					return
				}*/
		l.Infof("token解析后数据:%+v", account)
		ctx.Set("uid", uid)
		ctx.Set("gatewayId", account.Metadata["gatewayId"])
		ctx.Next()
	}
}
