package user

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	BearerToken = "Bearer"
	TokenExpire = 86400 // 24h
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// AuthenticationCtx validates token and authorizes users
func AuthenticationCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData := c.Request.Header.Get("Authorization")
		clientToken := strings.Replace(tokenData, BearerToken+" ", "", 1)
		if clientToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		claims, err := ValidateToken(clientToken)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user", claims.Email)
		c.Next()
	}
}

// SignedDetails
type SignedDetails struct {
	Email string
	jwt.StandardClaims
}

// GenerateAllTokens generates both the detailed token and refresh token
func GenerateAllTokens(email string) (string, string) {
	claims := &SignedDetails{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * TokenExpire).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return "", ""
	}

	return token, refreshToken
}

//ValidateToken validates the jwt token
func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, fmt.Errorf("the token is invalid")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Errorf("the token has expired")
	}

	return claims, nil
}
