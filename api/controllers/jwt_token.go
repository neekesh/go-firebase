package controllers

import (
	"fmt"
	"letschat/api/helper"
	"letschat/api/services"
	"letschat/api/validators"
	"letschat/dtos"
	"letschat/errors"
	"letschat/infrastructure"
	"letschat/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JwtAuthController -> struct
type JwtAuthController struct {
	logger      infrastructure.Logger
	userService services.UserService
	jwtService  services.JWTAuthService
	env         infrastructure.Env
	validator   validators.UserValidator
}

// NewJwtAuthController -> constructor
func NewJwtAuthController(
	logger infrastructure.Logger,
	userService services.UserService,
	jwtService services.JWTAuthService,
	env infrastructure.Env,
	validator validators.UserValidator,
) JwtAuthController {
	return JwtAuthController{
		logger:      logger,
		userService: userService,
		jwtService:  jwtService,
		env:         env,
		validator:   validator,
	}
}

func (cc JwtAuthController) ObtainJwtToken(c *gin.Context) {
	reqData := dtos.JWTLoginRequestData{}

	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Zap.Error("Error [ShouldBindJSON] : ", err.Error())
		err := errors.BadRequest.Wrap(err, "Failed to bind request data")
		responses.HandleError(c, err)
		return
	}
	// validating using custom validator
	if validationErr := cc.validator.Validate.Struct(reqData); validationErr != nil {
		cc.logger.Zap.Error("[Validate Struct] Validation error: ", validationErr.Error())
		err := errors.BadRequest.Wrap(validationErr, "Validation error")
		err = errors.SetCustomMessage(err, "Invalid input information")
		err = errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		responses.HandleError(c, err)
		return
	}

	user, exits, err := cc.userService.CheckUserWithPhone(reqData.Phone)
	if err != nil {
		cc.logger.Zap.Error("Something went wrong: ", err.Error())
		responses.HandleError(c, err)
		return
	}
	if !exits {
		responses.ErrorJSON(c, http.StatusBadRequest, "Phone number or Password does not match")
		return
	}

	isValidPassword := helper.CompareHashAndPlainPassword(user.Password, reqData.Password)
	if !isValidPassword {
		cc.logger.Zap.Error("[CompareHashAndPassword] hash and plain password doesnot match")
		responses.ErrorJSON(c, http.StatusBadRequest, "Invalid user credentials")
		return
	}

	// Create a new JWT access claims object
	accessClaims := services.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(cc.env.JwtAccessTokenExpiresAt)).Unix(),
			Id:        fmt.Sprintf("%v", user.UserId),
		},
	}
	// Create a new JWT Access token using the claims and the secret key
	accessToken, tokenErr := cc.jwtService.GenerateToken(accessClaims, cc.env.JwtAccessSecret)
	if tokenErr != nil {
		cc.logger.Zap.Error("[SignedString] Error getting token: ", tokenErr.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, tokenErr.Error())
		return
	}
	// Create a new JWT refresh claims object
	refreshClaims := services.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(cc.env.JwtRefreshTokenExpiresAt)).Unix(),
			Id:        fmt.Sprintf("%v", user.UserId),
		},
	}
	// Create a new JWT Refresh token using the claims and the secret key
	refreshToken, refreshTokenErr := cc.jwtService.GenerateToken(refreshClaims, cc.env.JwtRefreshSecret)
	if refreshTokenErr != nil {
		cc.logger.Zap.Error("[SignedString] Error getting token: ", refreshTokenErr.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, refreshTokenErr.Error())
		return
	}
	data := map[string]interface{}{
		"user":               user.ToMap(),
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
		"access_expires_at":  accessClaims.ExpiresAt,
		"refresh_expires_at": refreshClaims.ExpiresAt,
	}
	responses.SuccessJSON(c, http.StatusOK, data)
	return
}

func (cc JwtAuthController) RefreshJwtToken(c *gin.Context) {
	tokenString, err := cc.jwtService.GetTokenFromHeader(c)
	if err != nil {
		cc.logger.Zap.Error("Error getting token from header: ", err.Error())
		err = errors.Unauthorized.Wrap(err, "Something went wrong")
		responses.HandleError(c, err)
		return
	}
	parsedToken, parseErr := cc.jwtService.ParseToken(tokenString, cc.env.JwtRefreshSecret)
	if parseErr != nil {
		cc.logger.Zap.Error("Error parsing token: ", parseErr.Error())
		err = errors.Unauthorized.Wrap(parseErr, "Something went wrong")
		responses.HandleError(c, err)
		return
	}
	claims, verifyErr := cc.jwtService.VerifyToken(parsedToken)
	if verifyErr != nil {
		cc.logger.Zap.Error("Error veriefying token: ", verifyErr.Error())
		err = errors.Unauthorized.Wrap(verifyErr, "Something went wrong")
		responses.HandleError(c, err)
		return
	}
	// Create a new JWT Access claims
	accessClaims := services.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(cc.env.JwtAccessTokenExpiresAt)).Unix(),
			Id:        fmt.Sprintf("%v", claims.Id),
		},
		// Add other claims
	}
	// Create a new JWT token using the claims and the secret key
	accessToken, tokenErr := cc.jwtService.GenerateToken(accessClaims, cc.env.JwtAccessSecret)
	if tokenErr != nil {
		cc.logger.Zap.Error("[SignedString] Error getting token: ", tokenErr.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, tokenErr.Error())
		return
	}
	data := map[string]interface{}{
		"access_token": accessToken,
		"expires_at":   accessClaims.ExpiresAt,
	}
	responses.SuccessJSON(c, http.StatusOK, data)
	return

}
