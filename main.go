package main

import "github.com/gin-gonic/gin"
import "embed"
import "fmt"

//go:embed templates/*
var assets embed.FS

const listenPort = ":9091"

func main() {
        gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
        err := loadTemplate(r, assets)
	if err != nil {
		panic(err)
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "templates/index.tmpl", gin.H{
			"title": "Test",
		})
	})

	r.GET("/_data", func(c *gin.Context) {
		data, err := MakeDep2()
		if err != nil {
			panic(err)
		}
		c.JSON(200, data)
	})

	fmt.Printf("listen: %v\n", listenPort)
	r.Run(listenPort)
}
