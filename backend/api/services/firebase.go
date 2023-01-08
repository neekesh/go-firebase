package services

import (
	"context"
	"letschat/infrastructure"
	"letschat/models"

	"firebase.google.com/go/auth"
)

type FirebaseService struct {
	fbAuth *auth.Client
	logger infrastructure.Logger
}

func NewFirebaseService(
	fbAuth *auth.Client,
	logger infrastructure.Logger,
) FirebaseService {
	return FirebaseService{
		fbAuth: fbAuth,
		logger: logger,
	}
}

func (fb *FirebaseService) VerifyToken(idToken string) (*auth.Token, error) {
	token, err := fb.fbAuth.VerifyIDToken(context.Background(), idToken)
	return token, err
}

func (fb *FirebaseService) CreateCustomToken(uid string) (string, error) {
	token, err := fb.fbAuth.CustomToken(context.Background(), uid)
	return token, err
}

func (fb *FirebaseService) GetUser(uid string) (*auth.UserRecord, error) {
	user, err := fb.fbAuth.GetUser(context.Background(), uid)
	return user, err
}

func (fb *FirebaseService) CreateUser(newUser models.FirebaseAuthUser) (string, error) {
	fb.logger.Zap.Error("am i here")
	params := (&auth.UserToCreate{}).
		Email(newUser.Email).
		Password(newUser.Password).
		DisplayName(newUser.DisplayName)
	fb.logger.Zap.Info("password", newUser.Password)
	if newUser.Enabled == 0 {
		params = params.Disabled(true)
	}
	if newUser.Enabled == 1 {
		params = params.Disabled(false)
	}
	u, err := fb.fbAuth.CreateUser(context.Background(), params)
	if err != nil {
		return "", err
	}
	claims := map[string]interface{}{
		"role":   newUser.Role,
		"fb_uid": u.UID,
		"id":     newUser.UserId,
	}
	err = fb.fbAuth.SetCustomUserClaims(context.Background(), u.UID, claims)
	if err != nil {
		return "Internal Server Error", err
	}
	return u.DisplayName, err
}

func (fb *FirebaseService) GetUserByEmail(email string) string {
	user, _ := fb.fbAuth.GetUserByEmail(context.Background(), email)
	if user != nil {
		return user.UID
	}
	return ""
}

func (fb *FirebaseService) LoginUser(user models.FirebaseAuthUser) (string, error) {
	uid := fb.GetUserByEmail(user.Email)
	if len(uid) == 0 {
		return "Invalid Email", nil
	}
	tokens, err := fb.CreateCustomToken(uid)
	if err != nil {
		return "Error Getting tokens", nil
	}
	return tokens, err
}
