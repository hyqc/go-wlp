package middleware

import (
	"fmt"
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

func JWTCheck(jwt jwtToken.Provider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorization)
		if !strings.HasPrefix(authHeader, auth.BearerScheme) {
			response.Response(ctx, code.NewCode(code.AuthTokenInvalid))
			return
		}
		token := strings.TrimPrefix(authHeader, auth.BearerScheme)
		if token == "" {
			response.Response(ctx, code.NewCode(code.AuthTokenInvalid))
			return
		}
		// TODO
		fmt.Println("token: ", token)
		account, err := jwt.Inspect(token)
		if err != nil {
			response.Response(ctx, code.NewCode(code.AuthTokenInspectInvalid))
			return
		}

		uid, err := strconv.Atoi(account.ID)
		if err != nil {
			response.Response(ctx, code.NewCode(code.AuthTokenInfoInvalid))
			return
		}
		// TODO
		fmt.Println("token decode: ", account)
		ctx.Set("uid", uid)
		ctx.Next()
	}
}
