package appy_http

import "errors"

var ErrAuthorizationHeaderMissing = errors.New("authorization header missing")
var ErrApiKeysDontMatch = errors.New("api keys don't match")

var ErrInsufficientPermissions = errors.New("insufficient permissions")

var ErrTokenMalformed = errors.New("token malformed, could be missing 'Bearer' prefix")
var ErrTokenExpired = errors.New("token is expired")
var ErrTokenInvalid = errors.New("token invalid")
var ErrTokenMissingClaims = errors.New("token invalid, missing claims")
