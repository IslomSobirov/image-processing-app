package router

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

var (
	PathOfOriginalPhoto string
	ResultPathOfPhoto   string
	width               int
	height              int
)

func init() {
	flag.StringVar(&PathOfOriginalPhoto, "path-orig", "./tmp/img_orig/", "Папка оригинальных изображений")
	flag.StringVar(&ResultPathOfPhoto, "path-res", "./tmp/img_res/", "Папка изображений с измененным размером")
	flag.IntVar(&width, "width", 200, "Ширина изображений с измененным размером")
	flag.IntVar(&height, "height", 200, "Высота для изображений с измененным размером")
	flag.Parse()

	_ = os.MkdirAll(PathOfOriginalPhoto, os.ModePerm)
	_ = os.MkdirAll(ResultPathOfPhoto, os.ModePerm)
}

func RunRouter() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"*"})

	router.POST("/upload", func(c *gin.Context) {
		start := time.Now()

		var request struct {
			Image string `json:"image"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат запроса"})
			return
		}

		imgData, err := base64.StdEncoding.DecodeString(request.Image)
		if err != nil || len(imgData) < 1024 || len(imgData) > 10*1024*1024 {
			c.JSON(400, gin.H{"error": "Недопустимое или слишком большое изображение"})
			return
		}

		// Хешируем изображение для кэша
		hash := sha256.Sum256(imgData)
		hashStr := hex.EncodeToString(hash[:])

		origFile := filepath.Join(PathOfOriginalPhoto, fmt.Sprintf("%s.jpg", hashStr))

		resFile := filepath.Join(ResultPathOfPhoto, fmt.Sprintf("%s_%dx%d.jpg", hashStr, width, height))

		// Проверяем, существует ли уже измененное изображение
		if _, err := os.Stat(resFile); err == nil {
			elapsed := time.Since(start).Milliseconds()
			c.JSON(200, gin.H{"time": elapsed, "cached": true})
			return
		}

		// Сохраните исходное изображение
		if _, err := os.Stat(origFile); os.IsNotExist(err) {
			if err := os.WriteFile(origFile, imgData, 0644); err != nil {
				c.JSON(500, gin.H{"error": "Не удалось сохранить исходное изображение"})
				return
			}
		}

		img, _, err := image.Decode(strings.NewReader(string(imgData)))
		if err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат изображения"})
			return
		}

		bounds := img.Bounds()
		originalWidth := bounds.Dx()
		originalHeight := bounds.Dy()

		if originalHeight == 0 || originalWidth == 0 {
			c.JSON(400, gin.H{"error": "Высота и ширина изображения должны превышать 0"})
			return
		}

		resized := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

		out, err := os.Create(resFile)
		if err != nil {
			c.JSON(500, gin.H{"error": "Не удалось сохранить изображение"})
			return
		}
		defer out.Close()
		jpeg.Encode(out, resized, &jpeg.Options{Quality: 80})

		elapsed := time.Since(start).Milliseconds()
		c.JSON(200, gin.H{"time": elapsed, "cached": false})
	})

	router.Run(":8085")
}
