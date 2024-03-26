package utility

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/nfwGytautas/appy/driver"
	appy_middleware "github.com/nfwGytautas/appy/http/extensions/middleware"
)

var validate = validator.New()

// The size of pages for all requests
const PageSize = 20

// Utility struct for getting required handler parameters
type ParamChain struct {
	Context *gin.Context

	currentError error
}

func NewParamChain(context *gin.Context) *ParamChain {
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

func (pc *ParamChain) GetUser(outId *uint64, outName *string) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	token, exists := pc.Context.Get("accessToken")
	if !exists {
		panic("accesToken not found in context")
	}

	accessToken := token.(appy_middleware.AccessTokenInfo)

	*outId = uint64(accessToken.ID)
	*outName = accessToken.Username

	return pc
}

func (pc *ParamChain) GetPage(out *PagingSettings) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	pageString := pc.Context.Query("page")
	if pageString == "" {
		*out = PagingSettings{
			Count:  PageSize,
			Offset: 0,
		}
		return pc
	}

	numericalValue, err := strconv.Atoi(pageString)
	if err != nil {
		pc.currentError = err
		return pc
	}

	*out = PagingSettings{
		Count:  PageSize,
		Offset: uint64(numericalValue) * PageSize,
	}

	return pc
}

func (pc *ParamChain) ReadBodySingle(out any) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	body, err := io.ReadAll(pc.Context.Request.Body)
	if err != nil {
		pc.currentError = err
		return pc
	}

	err = json.Unmarshal(body, &out)
	if err != nil {
		pc.currentError = err
		return pc
	}

	err = validate.Struct(out)
	if err != nil {
		pc.currentError = err
		return pc
	}

	return pc
}

func (pc *ParamChain) ReadBodyArray(out any) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	body, err := io.ReadAll(pc.Context.Request.Body)
	if err != nil {
		pc.currentError = err
		return pc
	}

	err = json.Unmarshal(body, &out)
	if err != nil {
		pc.currentError = err
		return pc
	}

	s := reflect.ValueOf(out)
	s = s.Elem()

	for i := 0; i < s.Len(); i++ {
		err = validate.Struct(s.Index(i))
		if err != nil {
			pc.currentError = err
			return pc
		}
	}

	return pc
}

func (pc *ParamChain) ReadPathInt(name string, out *uint64) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	valueStr := pc.Context.Param(name)
	if valueStr == "" {
		pc.currentError = errors.New("missing parameter: " + name)
		return pc
	}

	numericalValue, err := strconv.Atoi(valueStr)
	if err != nil {
		pc.currentError = err
		return pc
	}

	*out = uint64(numericalValue)

	return pc
}

func (p *ParamChain) HasError() bool {
	return p.currentError != nil
}

func (p *ParamChain) Error() error {
	return p.currentError
}
