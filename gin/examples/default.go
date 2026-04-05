package main

import "github.com/gin-gonic/gin"

func minimalGin() {
	/*
		Default создает роутер с двумя middleware из коробки:
		Logger — логирует каждый запрос (метод, путь, статус, время)
		Recovery — ловит panic и возвращает 500 вместо краша
	*/
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Run(":8080")

}
