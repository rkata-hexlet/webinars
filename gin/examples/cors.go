package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/*
	CORS (Cross-Origin Resource Sharing):
	Проблема:
	У вас frontend на https://myapp.com, API на https://api.myapp.com.
 	Браузер блокирует запросы между разными доменами из соображений безопасности.

	Решение — CORS middleware:
*/

func corsMW() {
	r := gin.Default()
	r.Use(cors.Default()) // Разрешает все origins

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Run(":8080")

}

/*    Для production — настройте правильно:
config := cors.Config{
    AllowOrigins:     []string{"https://myapp.com"},  // ТОЛЬКО ваш домен!
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}

r.Use(cors.New(config))


КРИТИЧЕСКАЯ ошибка:
❌ AllowOrigins: []string{"*"} — НИКОГДА в production!
✅ AllowOrigins: []string{"https://myapp.com"} — конкретные домены
*/

/*

Security Headers:
r.Use(secureHeaders())

func secureHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Защита от clickjacking
        c.Header("X-Frame-Options", "DENY")

        // Защита от MIME-type sniffing
        c.Header("X-Content-Type-Options", "nosniff")

        // XSS защита
        c.Header("X-XSS-Protection", "1; mode=block")

        // HSTS для HTTPS
        c.Header("Strict-Transport-Security", "max-age=31536000")

        c.Next()
    }
}
*/
