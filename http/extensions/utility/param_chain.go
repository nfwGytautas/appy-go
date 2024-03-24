package utility

import (
	"github.com/nfwGytautas/appy/driver"
	appy_http "github.com/nfwGytautas/appy/http"
	appy_middleware "github.com/nfwGytautas/appy/http/extensions/middleware"
)

// The size of pages for all requests
const PageSize = 20

// Utility struct for getting required handler parameters
type ParamChain struct {
	Context *appy_http.HttpContext

	currentError error
}

func NewParamChain(context *appy_http.HttpContext) *ParamChain {
	return &ParamChain{Context: context, currentError: nil}
}

func (pc *ParamChain) OpenTransaction(out **driver.Tx) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	tx, err := driver.StartTransaction()
	if err != nil {
		pc.currentError = err
	}

	*out = tx

	return pc
}

func (pc *ParamChain) GetUserID(out *uint64) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	token, err := pc.Context.Get("accessToken")
	if err != nil {
		pc.currentError = err
		return pc
	}

	accessToken := token.(appy_middleware.AccessTokenInfo)

	*out = uint64(accessToken.ID)

	pc.Context.Tracker.SetUser(uint64(accessToken.ID), accessToken.Username)

	return pc
}

func (pc *ParamChain) GetPage(out *PagingSettings) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	*out = PagingSettings{
		Count:  PageSize,
		Offset: uint64(pc.Context.Query.Page()) * PageSize,
	}

	return pc
}

func (pc *ParamChain) ReadBodySingle(out any) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	err := pc.Context.Body.ParseSingle(out)
	if err != nil {
		pc.currentError = err
	}

	return pc
}

func (pc *ParamChain) ReadBodyArray(out any) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	err := pc.Context.Body.ParseArray(out)
	if err != nil {
		pc.currentError = err
	}

	return pc
}

func (pc *ParamChain) ReadPathInt(name string, out *uint64) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	value, err := pc.Context.Path.ExpectInt(name)
	if err != nil {
		pc.currentError = err
		return pc
	}

	*out = uint64(value)

	return pc
}

func (p *ParamChain) HasError() bool {
	return p.currentError != nil
}

func (p *ParamChain) Error() error {
	return p.currentError
}
