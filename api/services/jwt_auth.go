package services

import (
	"letschat/errors"
	"letschat/infrastructure"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTAuthService struct {
	logger infrastructure.Logger
	env    infrastructure.Env
}

func NewJWTAuthService(
	logger infrastructure.Logger,
	env infrastructure.Env,
) JWTAuthService {
	return JWTAuthService{
		logger: logger,
		env:    env,
	}
}

type JWTClaims struct {
	jwt.StandardClaims
}

func (m JWTAuthService) GetTokenFromHeader(c *gin.Context) (string, error) {
	// Get the token from the request header
	header := c.GetHeader("Authorization")
	if header == "" {
		err := errors.BadRequest.New("Authorization token is required in header")
		err = errors.SetCustomMessage(err, "Authorization token is required in header")
		m.logger.Zap.Error("[GetHeader]: ", err.Error())
		return "", err
	}

	if !strings.Contains(header, "Bearer") {
		err := errors.BadRequest.New("Token type is required")
		m.logger.Zap.Error("Missing token type: ", err.Error())
		return "", err
	}
	tokenString := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
	return tokenString, nil

}

func (m JWTAuthService) ParseToken(tokenString, secret string) (*jwt.Token, error) {
	// Parse the token using the secret key
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if !strings.Contains(err.Error(), "expired") {
			m.logger.Zap.Error("Invalid token[ParseWithClaims] :", err.Error())
			err := errors.BadRequest.New("Invalid ID token")
			return nil, err
		}
		m.logger.Zap.Error("Invalid token[ParseWithClaims] :", err.Error())
		return nil, err
	}
	return token, nil
}

func (m JWTAuthService) VerifyToken(token *jwt.Token) (*JWTClaims, error) {
	// Verfify token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		err := errors.BadRequest.New("Invalid token")
		err = errors.SetCustomMessage(err, "Invalid token")
		m.logger.Zap.Error("Invalid token [token.Valid]: ", err.Error())
		return nil, err
	}
	return claims, nil

}

func (m JWTAuthService) GenerateToken(claims JWTClaims, secret string) (string, error) {
	// Create a new JWT token using the claims and the secret key
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, tokenErr := tokenClaim.SignedString([]byte(secret))
	if tokenErr != nil {
		return "", tokenErr
	}
	return token, nil
}
