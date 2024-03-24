package appy_middleware

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	appy_http "github.com/nfwGytautas/appy/http"
)

// Struct for containing token info
type AccessTokenInfo struct {
	ID       uint
	Username string
	Role     string
}

type RefreshTokenInfo struct {
	ID uint
}

type JwtAuth struct {
	secret string
}

func NewJwtAuth(secret string) JwtAuth {
	return JwtAuth{
		secret: secret,
	}
}

func (j JwtAuth) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := j.ParseAccessToken(c)
		if err != nil {
			c.Abort()
			appy_http.Get().HandleError(c, err)
			return
		}

		c.Set("accessToken", info)

		c.Next()
	}
}

func (j JwtAuth) Authorization(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authenticate
		info, err := j.ParseAccessToken(c)
		if err != nil {
			c.Abort()
			appy_http.Get().HandleError(c, err)
			return
		}

		// Authorize
		if !isElementInArray(roles, info.Role) {
			c.Abort()
			appy_http.Get().HandleError(c, ErrInsufficientPermissions)
			return
		}

		c.Set("accessToken", info)

		c.Next()
	}
}

func (j JwtAuth) Generate(id uint, name, role string) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = id
	claims["name"] = name
	claims["role"] = role
	// claims["exp"] = time.Now().Add(5 * time.Minute).Unix()
	claims["exp"] = time.Now().Add(3 * 24 * time.Hour).Unix()

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

func (j JwtAuth) ParseAccessToken(c *gin.Context) (AccessTokenInfo, error) {
	result := AccessTokenInfo{}

	// Token empty check if it is inside Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return result, ErrAuthorizationHeaderMissing
	}

	// Since this is bearer token we need to parse the token out
	if len(strings.Split(tokenString, " ")) == 2 {
		tokenString = strings.Split(tokenString, " ")[1]
	} else {
		return result, ErrTokenMalformed
	}

	_, claims, err := j.parseToken(tokenString)
	if err != nil {
		return result, err
	}

	// User id
	uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["sub"]), 10, 32)
	if err != nil {
		return result, ErrTokenInvalid
	}

	result.ID = uint(uid)

	nameClaim := claims["name"]
	roleClaim := claims["role"]

	if nameClaim == nil || roleClaim == nil {
		return result, ErrTokenMissingClaims
	}

	result.Username = nameClaim.(string)
	result.Role = roleClaim.(string)

	return result, nil
}

func (j JwtAuth) ParseRefreshToken(c *gin.Context, token string) (RefreshTokenInfo, error) {
	result := RefreshTokenInfo{}

	_, claims, err := j.parseToken(token)
	if err != nil {
		return result, err
	}

	// User id
	uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["sub"]), 10, 32)
	if err != nil {
		return result, err
	}

	result.ID = uint(uid)

	return result, nil
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
