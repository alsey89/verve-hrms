package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"verve-hrms/internal/common"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (ah *AuthHandler) Signup(c echo.Context) error {
	creds := new(Credentials)
	err := c.Bind(creds)
	if err != nil {
		log.Printf("error binding credentials: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	username := creds.Username
	if username == "" {
		username = "New User" // default username
	}
	email := creds.Email
	if !common.EmailValidator(email) {
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid email",
			Data:    nil,
		})
	}

	password := creds.Password

	newUser, err := ah.authService.Signup(email, password, username)
	if err != nil {
		if err == ErrEmailNotAvailable {
			return c.JSON(http.StatusConflict, common.APIResponse{
				Message: "email is already in use",
				Data:    nil,
			})
		}
		log.Printf("error signing up: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	claims := Claims{
		ID:      newUser.ID, // Store the ObjectId
		IsAdmin: newUser.IsAdmin,
		Email:   newUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		log.Printf("Error signing jwt with claims: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = t
	cookie.HttpOnly = true
	cookie.Secure = viper.GetBool("IS_PRODUCTION")
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(time.Hour * 72)

	c.SetCookie(cookie)

	//! send event to worker thread
	// event := footprint.Event{
	// 	Name:      "_signedUp",
	// 	UserID:    newUserID,
	// 	TimeStamp: common.GetCurrentDateTime(),
	// }
	// worker.SendEvent(event)

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user has been signed up and signed in",
		Data:    newUser,
	})
}

func (ah *AuthHandler) Signin(c echo.Context) error {
	creds := new(Credentials)
	err := c.Bind(creds)
	if err != nil {
		log.Printf("error binding credentials: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	email := creds.Email
	password := creds.Password

	existingUser, err := ah.authService.Signin(email, password)
	if err != nil {
		if err == ErrInvalidCredentials {
			return c.JSON(http.StatusBadRequest, common.APIResponse{
				Message: "invalid credentials",
				Data:    nil,
			})
		}
		log.Printf("error signing in: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	claims := Claims{
		ID:      existingUser.ID, // Store the ObjectId
		IsAdmin: existingUser.IsAdmin,
		Email:   existingUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		log.Printf("error signing jwt with claims: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = t
	cookie.HttpOnly = true
	cookie.Secure = viper.GetBool("IS_PRODUCTION")
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(time.Hour * 72)

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user has been signed in",
		Data:    existingUser,
	})
}

func (ah *AuthHandler) Signout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Secure = viper.GetBool("IS_PRODUCTION")
	cookie.Path = "/"
	cookie.Expires = time.Unix(0, 0) //* set the cookie to expire immediately

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "user has been signed out",
		Data:    nil,
	})
}

func (ah *AuthHandler) CheckAuth(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token) //echo handles missing/malformed token response
	if !ok {
		log.Printf("error asserting token")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("error asserting claims: %v", user.Claims)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid claims data",
			Data:    nil,
		})
	}

	isAdmin, ok := claims["isAdmin"].(bool)
	if !ok {
		log.Printf("error asserting isAdmin: %v", claims["isAdmin"])
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "admin status not found",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "success",
		Data: echo.Map{
			"Authenticated": true,
			"IsAdmin":       isAdmin,
		},
	})
}
