package appy_middleware

import "errors"

var ErrTokenMalformed = errors.New("token malformed, could be missing 'Bearer' prefix")

var ErrTokenInvalid = errors.New("token invalid")

var ErrTokenMissingClaims = errors.New("token invalid, missing claims")

var ErrInsufficientPermissions = errors.New("insufficient permissions")

var ErrApiKeysDontMatch = errors.New("api keys don't match")

var ErrAuthorizationHeaderMissing = errors.New("authorization header missing")
