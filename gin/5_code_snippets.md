# Готовые snippets кода для Gin

## Auth Middleware (JWT)
```go
import "github.com/golang-jwt/jwt/v5"

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte("secret"), nil
        })
        if !token.Valid {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## CORS
```go
import "github.com/gin-contrib/cors"
r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"https://yoursite.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
}))
```

## Pagination
```go
type Pagination struct {
    Page  int `form:"page" binding:"gte=1"`
    Limit int `form:"limit" binding:"gte=1,lte=100"`
}
```

## Error Handling
```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) > 0 {
            c.JSON(500, gin.H{"error": c.Errors.Last().Error()})
        }
    }
}
```

## File Upload
```go
func uploadFile(c *gin.Context) {
    file, _ := c.FormFile("file")
    c.SaveUploadedFile(file, "./uploads/"+file.Filename)
    c.JSON(200, gin.H{"filename": file.Filename})
}
```
