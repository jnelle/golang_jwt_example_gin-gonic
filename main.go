package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"example.com/config"
	"example.com/mongo"
	"example.com/token"
	"github.com/gin-gonic/gin"
)

type bindFile struct {
	Token string                `form:"token" binding:"required"`
	File  *multipart.FileHeader `form:"file" binding:"required"`
}

func readConfig() *config.Config {
	c, err := config.ReadConf()
	if err != nil {
		panic(err)
	}
	return c
}

func initMongo(c string) *mongo.MongoDB_Client {
	db, err := mongo.NewMongoDB(c)
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	conf := readConfig()
	db := initMongo(conf.MongoDBURL)
	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20

	f, err := os.Create("backend_" + time.Now().Format("2006-Jan-02") + ".log")
	if err != nil {
		panic(err)
	}

	gin.DefaultWriter = io.MultiWriter(f)

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.Use(gin.Recovery())
	router.POST("/upload", func(c *gin.Context) {
		var bindedFile bindFile

		if err := c.ShouldBind(&bindedFile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Missing Values"})
			return
		}
		currentToken, isValid := token.ValidateToken(conf.SecretKey, bindedFile.Token)
		if isValid {
			file := bindedFile.File
			dst := filepath.Base(file.Filename)
			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Upload Failed"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "SUCCESS"})
			return
		} else {
			newTokenByte, newTokenString, err := token.GenerateToken(conf.SecretKey)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Token Generation Failed"})
				return
			}
			db.RefreshToken(currentToken, newTokenString)

			c.JSON(http.StatusOK, gin.H{"token": newTokenByte})
			return

		}

	})
	router.Run("0.0.0.0:3006")
}
