package appy_http

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var validate = validator.New()

type ginHeaderParser struct {
	ctx *gin.Context
}

type ginQueryParser struct {
	ctx *gin.Context
}

type ginPathParser struct {
	ctx *gin.Context
}

type ginBodyParser struct {
	ctx *gin.Context
}

func (g *ginHeaderParser) ExpectSingleString(key string) (string, error) {
	val := g.ctx.GetHeader(key)
	if val == "" {
		return "", errors.New("header '" + key + "' not specified")
	}

	return val, nil
}

func (g *ginQueryParser) GetString(key string) string {
	return g.ctx.Query(key)
}

func (g *ginQueryParser) GetInt(key string) int {
	value := g.ctx.DefaultQuery(key, "0")
	if value == "" {
		return 0
	}

	numericalValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return numericalValue
}

func (g *ginQueryParser) Page() int {
	return g.GetInt("page")
}

func (g *ginQueryParser) ExpectString(key string) (string, error) {
	val := g.ctx.Query(key)
	if val == "" {
		return "", errors.New("query parameter '" + key + "' not specified")
	}

	return val, nil
}

func (g *ginQueryParser) ExpectInt(key string) (int, error) {
	value := g.ctx.DefaultQuery(key, "0")
	if value == "" {
		return 0, errors.New("query parameter '" + key + "' not specified")
	}

	numericalValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return numericalValue, nil
}

func (g *ginQueryParser) ExpectPage() (int, error) {
	return g.ExpectInt("page")
}

func (g *ginPathParser) GetInt(key string) int {
	value := g.ctx.Param(key)
	if value == "" {
		return 0
	}

	numericalValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return numericalValue
}

func (g *ginPathParser) ExpectInt(key string) (int, error) {
	value := g.ctx.Param(key)
	if value == "" {
		return 0, errors.New("path parameter '" + key + "' not specified")
	}

	numericalValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return numericalValue, nil
}

func (g *ginBodyParser) ParseSingle(out any) error {
	body, err := io.ReadAll(g.ctx.Request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &out)
	if err != nil {
		return err
	}

	err = validate.Struct(out)
	if err != nil {
		return err
	}

	return nil
}

func (g *ginBodyParser) ParseArray(out any) error {
	body, err := io.ReadAll(g.ctx.Request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &out)
	if err != nil {
		return err
	}

	s := reflect.ValueOf(out)
	s = s.Elem()

	for i := 0; i < s.Len(); i++ {
		err = validate.Struct(s.Index(i))
		if err != nil {
			return err
		}
	}

	return nil
}
