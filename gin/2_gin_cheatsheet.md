# Шпаргалка по командам Gin

Краткий справочник для быстрого доступа к основным функциям, методам и структурам Gin Framework.

---

## Установка

```bash
# Установка Gin
go get -u github.com/gin-gonic/gin

# Установка с конкретной версией
go get github.com/gin-gonic/gin@v1.9.1
```

---

## Создание роутера

```go
import "github.com/gin-gonic/gin"

// Роутер с Logger и Recovery middleware
r := gin.Default()

// Пустой роутер без middleware
r := gin.New()

// Роутер с выборочными middleware
r := gin.New()
r.Use(gin.Logger())
r.Use(gin.Recovery())
```

---

## HTTP методы

```go
// GET запрос
r.GET("/path", handler)

// POST запрос
r.POST("/path", handler)

// PUT запрос
r.PUT("/path", handler)

// DELETE запрос
r.DELETE("/path", handler)

// PATCH запрос
r.PATCH("/path", handler)

// HEAD запрос
r.HEAD("/path", handler)

// OPTIONS запрос
r.OPTIONS("/path", handler)

// Любой метод
r.Any("/path", handler)

// Несколько методов
r.Match([]string{"GET", "POST"}, "/path", handler)
```

---

## Параметры запроса

### Path параметры

```go
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")  // string
})

// Wildcard
r.GET("/files/*filepath", func(c *gin.Context) {
    path := c.Param("filepath")  // /docs/file.pdf
})
```

### Query параметры

```go
// /search?q=golang&page=1
r.GET("/search", func(c *gin.Context) {
    query := c.Query("q")                    // "golang" или ""
    page := c.DefaultQuery("page", "1")      // "1" если нет параметра
    exists := c.QueryArray("tag")            // []string
})

// Получить все параметры
params := c.Request.URL.Query()  // url.Values
```

### POST данные

```go
// application/x-www-form-urlencoded
r.POST("/form", func(c *gin.Context) {
    name := c.PostForm("name")
    age := c.DefaultPostForm("age", "18")
})
```

---

## Ответы

### JSON

```go
// Map
c.JSON(200, gin.H{
    "message": "success",
    "data": data,
})

// Struct
type Response struct {
    Message string `json:"message"`
    Data    any    `json:"data"`
}
c.JSON(200, Response{Message: "success", Data: data})

// IndentedJSON (с отступами для debug)
c.IndentedJSON(200, data)

// SecureJSON (предотвращает JSON hijacking)
c.SecureJSON(200, data)

// JSONP (с callback)
c.JSONP(200, data)

// PureJSON (не экранирует HTML)
c.PureJSON(200, data)
```

### XML

```go
c.XML(200, gin.H{"message": "success"})
```

### YAML

```go
c.YAML(200, gin.H{"message": "success"})
```

### String

```go
c.String(200, "Hello %s", name)
```

### HTML

```go
// Загрузка шаблонов
r.LoadHTMLGlob("templates/*")

// Рендеринг
c.HTML(200, "index.html", gin.H{
    "title": "Main",
})
```

### Redirect

```go
// HTTP redirect
c.Redirect(301, "https://google.com")

// Route redirect
c.Request.URL.Path = "/new-path"
r.HandleContext(c)
```

### File

```go
// Отправка файла
c.File("/path/to/file.pdf")

// С attachment (скачивание)
c.FileAttachment("/path/to/file.pdf", "download.pdf")

// Из файловой системы
c.FileFromFS("file.pdf", http.Dir("/static"))
```

### Data

```go
// Сырые данные
c.Data(200, "application/octet-stream", []byte("data"))

// Stream
c.Stream(func(w io.Writer) bool {
    w.Write([]byte("chunk"))
    return true  // false для остановки
})
```

### Status

```go
// Только статус
c.Status(204)

// Abort с статусом
c.AbortWithStatus(404)

// Abort с JSON
c.AbortWithStatusJSON(400, gin.H{"error": "bad request"})
```

---

## Binding и валидация

### JSON Binding

```go
type User struct {
    Name string `json:"name" binding:"required"`
    Age  int    `json:"age" binding:"gte=0"`
}

var user User

// ShouldBindJSON - не прерывает при ошибке
if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}

// BindJSON - возвращает 400 автоматически
if err := c.BindJSON(&user); err != nil {
    return  // уже отправлен ответ
}
```

### Query Binding

```go
type Pagination struct {
    Page  int `form:"page" binding:"required,gte=1"`
    Limit int `form:"limit" binding:"required,gte=1,lte=100"`
}

var page Pagination
if err := c.ShouldBindQuery(&page); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### URI Binding

```go
type GetUser struct {
    ID int `uri:"id" binding:"required,gte=1"`
}

var user GetUser
if err := c.ShouldBindUri(&user); err != nil {
    c.JSON(400, gin.H{"error": "invalid id"})
    return
}
```

### Header Binding

```go
type Headers struct {
    Token string `header:"Authorization" binding:"required"`
}

var h Headers
if err := c.ShouldBindHeader(&h); err != nil {
    c.JSON(401, gin.H{"error": "no token"})
    return
}
```

### Form Binding

```go
type LoginForm struct {
    User string `form:"user" binding:"required"`
    Pass string `form:"password" binding:"required"`
}

var form LoginForm
if err := c.ShouldBind(&form); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### Multipart Form

```go
// Один файл
file, _ := c.FormFile("file")
c.SaveUploadedFile(file, "./uploads/"+file.Filename)

// Несколько файлов
form, _ := c.MultipartForm()
files := form.File["upload"]
for _, file := range files {
    c.SaveUploadedFile(file, "./uploads/"+file.Filename)
}
```

---

## Validation Tags

```go
type Example struct {
    // Обязательное поле
    Required string `binding:"required"`
    
    // Длина строки
    MinStr   string `binding:"min=3"`
    MaxStr   string `binding:"max=10"`
    LenStr   string `binding:"len=5"`
    
    // Числовые значения
    MinNum   int     `binding:"min=1"`
    MaxNum   int     `binding:"max=100"`
    GTE      int     `binding:"gte=18"`  // >=
    LTE      int     `binding:"lte=100"` // <=
    GT       int     `binding:"gt=0"`    // >
    LT       int     `binding:"lt=100"`  // <
    
    // Email и URL
    Email    string  `binding:"email"`
    URL      string  `binding:"url"`
    
    // Одно из значений
    Color    string  `binding:"oneof=red green blue"`
    
    // UUID
    UUID     string  `binding:"uuid"`
    
    // Опциональное (если есть - валидировать)
    Optional string  `binding:"omitempty,email"`
    
    // Массив
    Tags     []string `binding:"required,min=1"`
    
    // Вложенная структура
    Address  Address  `binding:"required"`
}
```

---

## Middleware

### Использование middleware

```go
// Глобально
r.Use(middleware1())
r.Use(middleware2())

// Для конкретного роута
r.GET("/path", middleware1(), middleware2(), handler)

// Для группы
admin := r.Group("/admin")
admin.Use(authMiddleware())
```

### Создание middleware

```go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Код ДО обработчика
        
        c.Next()  // Вызов следующего middleware/handler
        
        // Код ПОСЛЕ обработчика
    }
}

// С параметрами
func AuthMiddleware(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Проверка роли
        c.Next()
    }
}
```

### Прерывание цепочки

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !isAuthorized(c) {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return  // Прерывание
        }
        c.Next()
    }
}
```

### Встроенные middleware

```go
// Logger
r.Use(gin.Logger())

// Recovery (от panic)
r.Use(gin.Recovery())

// BasicAuth
authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
    "user": "password",
}))
```

---

## Группы роутов

```go
// Простая группа
api := r.Group("/api")
{
    api.GET("/users", getUsers)
    api.POST("/users", createUser)
}

// Вложенные группы
api := r.Group("/api")
{
    v1 := api.Group("/v1")
    {
        v1.GET("/users", getUsersV1)
    }
    
    v2 := api.Group("/v2")
    {
        v2.GET("/users", getUsersV2)
    }
}

// С middleware
admin := r.Group("/admin")
admin.Use(AuthMiddleware())
{
    admin.GET("/users", adminGetUsers)
}
```

---

## Context методы

### Данные запроса

```go
// Headers
token := c.GetHeader("Authorization")
c.Header("X-Custom", "value")  // Установить header ответа

// Cookies
cookie, _ := c.Cookie("name")
c.SetCookie("name", "value", 3600, "/", "domain", false, true)

// IP адрес
ip := c.ClientIP()

// User Agent
ua := c.Request.UserAgent()

// Метод и путь
method := c.Request.Method
path := c.Request.URL.Path
```

### Хранение данных

```go
// Установить значение
c.Set("key", "value")

// Получить значение
value, exists := c.Get("key")
if exists {
    str := value.(string)
}

// MustGet (паникует если нет)
value := c.MustGet("key")
```

### Контроль выполнения

```go
// Прервать с ошибкой
c.Abort()
c.AbortWithStatus(404)
c.AbortWithStatusJSON(400, gin.H{"error": "bad request"})
c.AbortWithError(500, err)

// Проверить прерывание
if c.IsAborted() {
    return
}

// Вызвать следующий handler
c.Next()
```

---

## Запуск сервера

```go
// На порту 8080
r.Run()

// На конкретном порту
r.Run(":3000")

// На конкретном адресе
r.Run("192.168.1.1:8080")

// HTTPS
r.RunTLS(":443", "cert.pem", "key.pem")

// Unix socket
r.RunUnix("/tmp/gin.sock")

// Custom HTTP server
srv := &http.Server{
    Addr:    ":8080",
    Handler: r,
}
srv.ListenAndServe()
```

---

## Режимы работы

```go
// Установить режим
gin.SetMode(gin.ReleaseMode)  // Production
gin.SetMode(gin.DebugMode)    // Development (default)
gin.SetMode(gin.TestMode)     // Testing

// Через environment variable
// GIN_MODE=release go run main.go
```

---

## Логирование

```go
// Логирование в файл
f, _ := os.Create("gin.log")
gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

// Disable console color
gin.DisableConsoleColor()

// Force console color
gin.ForceConsoleColor()

// Кастомный формат логов
r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
    return fmt.Sprintf("[%s] %s %s %d\n",
        param.TimeStamp.Format("2006/01/02 - 15:04:05"),
        param.Method,
        param.Path,
        param.StatusCode,
    )
}))
```

---

## Тестирование

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    // Test
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ping", nil)
    r.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "pong")
}
```

---

## Полезные паттерны

### Error Handling

```go
func handler(c *gin.Context) {
    if err := doSomething(); err != nil {
        c.Error(err)  // Сохранить ошибку
        c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
        return
    }
}

// Middleware для обработки ошибок
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            // Обработать все ошибки
            for _, e := range c.Errors {
                log.Println(e.Err)
            }
        }
    }
}
```

### Response Wrapper

```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(200, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(code, Response{
        Code:    code,
        Message: message,
    })
}
```

### Dependency Injection

```go
type Handler struct {
    db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
    return &Handler{db: db}
}

func (h *Handler) GetUsers(c *gin.Context) {
    var users []User
    h.db.Find(&users)
    c.JSON(200, users)
}

// Использование
h := NewHandler(db)
r.GET("/users", h.GetUsers)
```

---

## Полезные ссылки

- Документация: https://gin-gonic.com/docs/
- GitHub: https://github.com/gin-gonic/gin
- Examples: https://github.com/gin-gonic/examples
- Middleware: https://github.com/gin-contrib
