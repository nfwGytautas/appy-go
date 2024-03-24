package appy_http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

func (c *HttpContext) Set(key string, value any) {
	if c.tempStorage == nil {
		c.tempStorage = make(map[string]any)
	}

	c.tempStorage[key] = value
}

func (c *HttpContext) Get(key string) (any, error) {
	if c.tempStorage == nil {
		return nil, errors.New("no values in temporary storage")
	}

	value, ok := c.tempStorage[key]
	if !ok {
		return nil, errors.New("key '" + key + "' not found in temporary storage")
	}

	return value, nil
}

func (c *HttpContext) StoreMultipartFile(key string, outDir string) (string, HttpResult) {
	file, header, err := c.Request.FormFile(key)
	if err != nil {
		return "", c.Error(err)
	}

	// Create a name uuid
	extension := strings.Split(header.Filename, ".")[1]
	filename := uuid.New().String() + "." + extension

	// Store locally
	out, err := os.Create(outDir + filename)
	if err != nil {
		return "", c.Error(err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", c.Error(err)
	}

	// File route
	return fmt.Sprintf("/images/%v", filename), c.Nil()
}

func (c *HttpContext) Nil() HttpResult {
	return HttpResult{
		failed: false,
	}
}

func (c *HttpContext) Ok(statusCode int, body interface{}) HttpResult {
	return HttpResult{
		StatusCode: statusCode,
		Body:       body,
		failed:     false,
	}
}

func (c *HttpContext) NotFound() HttpResult {
	return HttpResult{
		StatusCode: http.StatusNotFound,
		failed:     true,
	}
}

func (c *HttpContext) BadRequest(body interface{}) HttpResult {
	return HttpResult{
		StatusCode: http.StatusBadRequest,
		Body:       body,
		failed:     true,
	}
}

func (c *HttpContext) Fail(statusCode int, body interface{}) HttpResult {
	return HttpResult{
		StatusCode: statusCode,
		Body:       body,
		failed:     true,
	}
}

func (c *HttpContext) Error(err error) HttpResult {
	return HttpResult{
		StatusCode: http.StatusInternalServerError,
		Error:      err,
		failed:     true,
		Tracker:    c.getTrackerInfo(),
	}
}

func (hr HttpResult) IsFailed() bool {
	return hr.failed || hr.Error != nil
}

func (hr HttpResult) HasError() bool {
	return hr.Error != nil
}

func (c *HttpContext) getTrackerInfo() HttpResultTrackerInfo {
	_, file, line, _ := runtime.Caller(2)
	return HttpResultTrackerInfo{
		At: fmt.Sprintf("%v:%v", file, line),
	}
}
