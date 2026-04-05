package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func customMW() {
	r := gin.Default()

	// Глобально для всех роутов
	r.Use(MyMiddleware())

	// Для конкретного роута
	r.GET("/greeting", GreetingMiddleware(), func(c *gin.Context) {
		fmt.Println("I'm greeting middleware!")
		c.JSON(200, gin.H{"message": "hello"})
	})

	r.GET("/ping", func(c *gin.Context) {
		fmt.Println("I'm ping/pong middleware!")
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Run(":8080")

}

// ///////////////////////////////////////////////////////////////////////////

/*
	Best practices:

	Порядок имеет значение: Logger → Recovery → Auth → Your Logic
	Используйте Abort() для прерывания цепочки
	Не забывайте вызывать Next() если хотите продолжить
	Держите middleware простыми — одна ответственность
	Используйте группы для организации — не применяйте middleware глобально если не нужно
*/

func MyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// КОД ДО основного обработчика
		fmt.Println("MyMiddleware Before")

		c.Next() // Вызов следующего middleware или обработчика

		// КОД ПОСЛЕ основного обработчика
		fmt.Println("MyMiddleware After")
	}
}

func GreetingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// КОД ДО основного обработчика
		fmt.Println("GreetingMiddleware Before")

		c.Next() // Вызов следующего middleware или обработчика

		// КОД ПОСЛЕ основного обработчика
		fmt.Println("GreetingMiddleware After")
	}
}

// //////////////////////////////////////////// Bonus
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		// Проверяем наличие токена
		if token == "" {
			c.JSON(401, gin.H{"error": "No authorization header"})
			c.Abort() // Прерываем цепочку middleware
			return
		}

		// Проверяем валидность токена (упрощенно)
		if token != "Bearer secret-token" {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Можем добавить данные в контекст
		c.Set("user_id", 123)

		c.Next() // Продолжаем цепочку
	}
}
