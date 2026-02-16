package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		// TODO: answer here
		cookie, err := ctx.Request.Cookie("session_token")
		if err != nil {
			if ctx.GetHeader("Content-Type") == "application/json" {
				ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse("error unauthorized user id"))
			} else {
				ctx.Redirect(http.StatusSeeOther, "/page/login")
			}
			return
		}

		claims := &model.Claims{}

		txn, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(err.Error()))
				return
			}
			ctx.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
			return
		}

		if !txn.Valid {
			ctx.JSON(http.StatusUnauthorized, model.NewErrorResponse(err.Error()))
			return
		}

		ctx.Set("id", claims.ID)
		ctx.Set("email", claims.Email)
		ctx.Next()
	})
}
