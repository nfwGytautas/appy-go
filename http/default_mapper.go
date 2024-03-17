package appy_driver_http

import "github.com/nfwGytautas/appy"

type defaultErrorMapper struct{}

func (dem *defaultErrorMapper) Map(err error) appy.HttpResult {
	return appy.HttpResult{
		StatusCode: 500,
		Body:       "Internal Server Error",
	}
}
