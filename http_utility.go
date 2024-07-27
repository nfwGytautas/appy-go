package appy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

var validate = validator.New()

// The size of pages for all requests
const PageSize = 20

// Utility struct for getting required handler parameters
type ParamChain struct {
	Context *gin.Context

	currentError error
}

type PagingSettings struct {
	Offset uint64
	Count  uint64
}

func (ps PagingSettings) String() string {
	return fmt.Sprintf("{Offset: %v, Count: %v}", ps.Offset, ps.Count)
}

func NewParamChain(context *gin.Context) *ParamChain {
	return &ParamChain{Context: context, currentError: nil}
}

func (pc *ParamChain) GetUser(outId *uint64, outName *string) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	token, exists := pc.Context.Get("accessToken")
	if !exists {
		panic("accesToken not found in context")
	}

	accessToken := token.(AccessTokenInfo)

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

func (pc *ParamChain) GetPageSized(out *PagingSettings, size uint64) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	pageString := pc.Context.Query("page")
	if pageString == "" {
		*out = PagingSettings{
			Count:  size,
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
		Count:  size,
		Offset: uint64(numericalValue) * size,
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

func (pc *ParamChain) ReadQueryInt(name string, out *uint64) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	valueStr := pc.Context.Query(name)
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

func (pc *ParamChain) ReadQueryString(name string, out *string) *ParamChain {
	if pc.currentError != nil {
		return pc
	}

	valueStr := pc.Context.Query(name)
	if valueStr == "" {
		pc.currentError = errors.New("missing parameter: " + name)
		return pc
	}

	*out = valueStr

	return pc
}

func (p *ParamChain) HasError() bool {
	return p.currentError != nil
}

func (p *ParamChain) Error() error {
	return p.currentError
}

func StoreMultipartFile(c *gin.Context, key string, outDir string) (string, error) {
	file, header, err := c.Request.FormFile(key)
	if err != nil {
		return "", err
	}

	// Create a name uuid
	extension := strings.Split(header.Filename, ".")[1]
	filename := uuid.New().String() + "." + extension

	// Store locally
	out, err := os.Create(outDir + filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	// File route
	return fmt.Sprintf("/images/%v", filename), nil
}

func UnimplementedEndpoint(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "This endpoint is not implemented yet",
	})
}
