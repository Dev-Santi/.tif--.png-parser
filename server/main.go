package main

import (
	"archive/zip"
	"bytes"
	"image/png"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/tiff"
)

func main() {
	router := gin.Default()
	router.Use(setCors)

	router.POST("/convert", convertFiles)
	router.Run("localhost:8080")
}

func convertFiles(c *gin.Context) {
	// Multipart form with multiple files, all named "tif"
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to parse multipart form: %v", err)
		return
	}

	files := form.File["tif"]
	if len(files) == 0 {
		c.String(http.StatusBadRequest, "No files uploaded with field 'tif'")
		return
	}

	// Buffer to write ZIP archive into memory
	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to open file %s: %v", file.Filename, err)
			return
		}

		img, err := tiff.Decode(f)
		f.Close()
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid TIFF file %s: %v", file.Filename, err)
			return
		}

		// Prepare filename for PNG
		base := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
		pngName := base + ".png"

		// Create a file inside the ZIP archive
		zipFile, err := zipWriter.Create(pngName)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to create zip entry: %v", err)
			return
		}

		// Encode PNG directly into the zip file
		err = png.Encode(zipFile, img)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to encode PNG %s: %v", pngName, err)
			return
		}
	}

	err = zipWriter.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to close zip archive: %v", err)
		return
	}

	// Send the ZIP file as response
	c.Header("Content-Disposition", "attachment; filename=converted_images.zip")
	c.Data(http.StatusOK, "application/zip", zipBuffer.Bytes())
}

func setCors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
