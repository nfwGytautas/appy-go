package appy_config

import "context"

type HttpErrorMapper interface {
	Map(context.Context, error) (int, any)
}
