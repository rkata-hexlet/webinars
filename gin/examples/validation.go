package main

import "github.com/gin-gonic/gin"

// Gin — встроенная валидация
type User struct {
	Name  string `json:"name" binding:"required,min=3"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"gte=18,lte=100"`
}

type SearchQuery struct {
	Query string `form:"q" binding:"required,min=3"`
	Page  int    `form:"page" binding:"gte=1"`
}

// Валидация query параметров:
func search(c *gin.Context) {
	var query SearchQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// query.Query и query.Page валидны
}

// Валидация URI параметров:
type GetUserURI struct {
	ID int `uri:"id" binding:"required,gte=1"`
}

func getUser(c *gin.Context) {
	var uri GetUserURI

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	// uri.ID валиден
}

/*
Best practices:

Всегда валидируйте входящие данные
Используйте omitempty для опциональных полей
Валидируйте на уровне приложения, не только в БД
Возвращайте понятные сообщения об ошибках
Не полагайтесь только на front-end валидацию
*/
