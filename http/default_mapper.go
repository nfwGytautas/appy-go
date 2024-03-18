package appy_driver_http

import "github.com/nfwGytautas/appy"

type defaultErrorMapper struct{}

func (dem *defaultErrorMapper) Map(res *appy.HttpResult) {}
