package main

import (
    "context"
    "fmt"
    _ "github.com/akurczyk/golang_echo_blogging_platform/002_mysql_and_redis/app/src/docs"
    "github.com/go-playground/validator"
    "github.com/go-redis/redis/v8"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo-contrib/prometheus"
    "github.com/labstack/echo/v4/middleware"
    "github.com/swaggo/echo-swagger"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "os"
    "strconv"
)

type (
    CustomValidator struct {
        validator *validator.Validate
    }
)

var (
    sql_db   *gorm.DB
    ctx      context.Context
    redis_db *redis.Client
)

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func setupSql() {
    var err error

    host := os.Getenv("MYSQL_HOST")
    port := os.Getenv("MYSQL_PORT")
    user := os.Getenv("MYSQL_USER")
    password := os.Getenv("MYSQL_PASSWORD")
    dbname := os.Getenv("MYSQL_DATABASE")

    template := "%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local"
    dsn := fmt.Sprintf(template, user, password, host, port, dbname)
    sql_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Could not connect to database.")
    }

    err = sql_db.AutoMigrate(&Post{})
    if err != nil {
        panic("Could not migrate posts.")
    }

    err = sql_db.AutoMigrate(&Comment{})
    if err != nil {
        panic("Could not migrate comments.")
    }

    err = sql_db.AutoMigrate(&User{})
    if err != nil {
        panic("Could not migrate users.")
    }
}

func setupRedis() {
    host := os.Getenv("REDIS_HOST")
    port := os.Getenv("REDIS_PORT")
    password := os.Getenv("REDIS_PASSWORD")
    db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

    ctx = context.Background()

    redis_db = redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%v:%v", host, port),
        Password: password,
        DB:       db,
    })
}

// @title Simple blogging platform API based on MySQL and Redis databases
// @version 1.0
// @description Simple blogging platform API created in Golang with the use of Echo framework, GORM ORM library,
// @description MySQL database for storing objects, and Redis for storing session data with integrated Prometheus
// @description and Swagger dockerized.

// @contact.name Aleksander Kurczyk
// @contact.url http://github.com/akurczyk
// @contact.email akurczyk000@gmail.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
    setupSql()
    setupRedis()

    e := echo.New()

    // Validator
    e.Validator = &CustomValidator{validator: validator.New()}

    // Prometheus
    p := prometheus.NewPrometheus("echo", nil)
    p.Use(e)

    // Swagger
    e.GET("/swagger/*", echoSwagger.WrapHandler)

    e.GET("/users", listUserAccounts)
    e.POST("/users", createUserAccount)
    e.GET("/users/:id", retrieveUserAccount)
    e.PUT("/users", updateUserAccount, middleware.KeyAuth(checkAuthToken))
    e.DELETE("/users", deleteUserAccount, middleware.KeyAuth(checkAuthToken))

    e.POST("/token", issueAuthToken)
    e.DELETE("/token", revokeAuthToken, middleware.KeyAuth(checkAuthToken))

    e.GET("/posts", listPosts)
    e.POST("/posts", createPost, middleware.KeyAuth(checkAuthToken))
    e.GET("/posts/:id", retrievePost)
    e.PUT("/posts/:id", updatePost, middleware.KeyAuth(checkAuthToken))
    e.DELETE("/posts/:id", deletePost, middleware.KeyAuth(checkAuthToken))

    e.GET("/comments", listComments)
    e.POST("/comments", createComment, middleware.KeyAuth(checkAuthToken))
    e.GET("/comments/:id", retrieveComment)
    e.PUT("/comments/:id", updateComment, middleware.KeyAuth(checkAuthToken))
    e.DELETE("/comments/:id", deleteComment, middleware.KeyAuth(checkAuthToken))

    e.Logger.Fatal(e.Start(":1323"))
}

// TODO: Tests
