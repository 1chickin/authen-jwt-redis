package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/1chickin/authen-jwt-redis/config"
	"github.com/1chickin/authen-jwt-redis/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// get token from cookie
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// validate token
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["id"], claims["expire"])
		// check expiration
		if float64(time.Now().Unix()) > claims["expire"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// check user exist
		var user model.User
		config.DB.First(&user, claims["id"])
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	} else {
		fmt.Println(err)
	}

	//continue
	c.Next()
}
