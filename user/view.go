package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginPayload struct is used to store incoming login requests.
type LoginPayload struct {
	ClientID string `json:"client_id" form:"client_id"`
	Secret   string `json:"client_secret" form:"client_secret"`
}

// TokenResponse is the authorization server response
type TokenResponse struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"` // bearer
	ExpiresIn    int64  `json:"expires_in"` // secs
}

// VerifyPassword mock
func VerifyPassword(userPassword string, providedPassword string) error {
	if userPassword != providedPassword {
		return fmt.Errorf("login or passowrd is incorrect")
	}
	return nil
}

// Login user action
func Login(c *gin.Context) {
	var loginPayload *LoginPayload
	var foundUser *User
	var err error

	if err := c.ShouldBind(&loginPayload); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundUser, err = DBGetUserByEmail(loginPayload.ClientID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "login or passowrd is incorrect"})
		return
	}

	err = VerifyPassword(loginPayload.Secret, foundUser.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	token, refreshToken := GenerateAllTokens(foundUser.Email)
	foundUser.Token = token
	foundUser.RefreshToken = refreshToken
	_, err = DBUpdateUser(foundUser.ID, foundUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
		TokenType:    BearerToken,
		ExpiresIn:    TokenExpire,
	})
}
