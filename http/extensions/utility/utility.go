package utility

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
