package appy_middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nfwGytautas/appy"
)

// Struct for containing token info
type AccessTokenInfo struct {
	ID   uint
	Role string
}

type RefreshTokenInfo struct {
	ID uint
}

type JwtAuth struct {
	secret string
}

type tokenInfo struct {
	jwtToken *jwt.Token
	claims   jwt.MapClaims
	expired  bool
}

func NewJwtAuth(secret string) JwtAuth {
	return JwtAuth{
		secret: secret,
	}
}

func (j JwtAuth) Authentication() appy.HttpMiddleware {
	return func(c *appy.HttpContext) appy.HttpResult {
		info, res := j.ParseAccessToken(c)
		if res.IsFailed() {
			return res
		}

		c.Set("accessToken", info)

		return c.Nil()
	}
}

func (j JwtAuth) Authorization(roles []string) appy.HttpMiddleware {
	return func(c *appy.HttpContext) appy.HttpResult {
		// Authenticate
		info, res := j.ParseAccessToken(c)
		if res.IsFailed() {
			return res
		}

		// Authorize
		if !isElementInArray(roles, info.Role) {
			return c.Fail(http.StatusUnauthorized, "Access denied, insufficient permissions")
		}

		c.Set("accessToken", info)

		return c.Nil()
	}
}

func (j JwtAuth) Generate(id uint, role string) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(5 * time.Minute).Unix()

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS512)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = id
	rtClaims["exp"] = time.Now().Add(30 * time.Hour * 24).Unix()

	refreshTokenString, err := refreshToken.SignedString([]byte(j.secret))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func (j JwtAuth) ParseAccessToken(c *appy.HttpContext) (AccessTokenInfo, appy.HttpResult) {
	result := AccessTokenInfo{}

	// Token empty check if it is inside Authorization header
	tokenString, err := c.Header.ExpectSingleString("Authorization")
	if err != nil {
		return result, c.Error(err)
	}

	// Since this is bearer token we need to parse the token out
	if len(strings.Split(tokenString, " ")) == 2 {
		tokenString = strings.Split(tokenString, " ")[1]
	} else {
		return result, c.Fail(http.StatusUnauthorized, "Invalid token")
	}

	_, claims, err := j.parseToken(tokenString)
	if err != nil {
		return result, c.Fail(http.StatusUnauthorized, err.Error())
	}

	// User id
	uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["sub"]), 10, 32)
	if err != nil {
		return result, c.Fail(http.StatusUnauthorized, "Invalid token")
	}

	result.ID = uint(uid)
	result.Role = claims["role"].(string)

	return result, c.Nil()
}

func (j JwtAuth) ParseRefreshToken(c *appy.HttpContext, token string) (RefreshTokenInfo, appy.HttpResult) {
	result := RefreshTokenInfo{}

	_, claims, err := j.parseToken(token)
	if err != nil {
		return result, c.Error(err)
	}

	// User id
	uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["sub"]), 10, 32)
	if err != nil {
		return result, c.Error(err)
	}

	result.ID = uint(uid)

	return result, c.Nil()
}

func (j JwtAuth) parseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if !jwtToken.Valid {
		return nil, nil, errors.New("invalid token")
	}

	// Token valid fill token information
	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	if !ok {
		return nil, nil, errors.New("failed to get claims")
	}

	// Expiration
	timeStamp := claims["exp"]
	validity, ok := timeStamp.(float64)
	if !ok {
		return nil, nil, errors.New("invalid token")
	}

	tm := time.Unix(int64(validity), 0)
	remainer := time.Until(tm)
	if remainer <= 0 {
		return nil, nil, errors.New("token expired")
	}

	return jwtToken, claims, nil
}

func isElementInArray[T comparable](arr []T, val T) bool {
	for _, element := range arr {
		if element == val {
			return true
		}
	}

	return false
}
