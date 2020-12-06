package main

import (
    "context"
    "fmt"
    _ "github.com/akurczyk/golang_echo_blogging_platform/001_postgres_and_redis/app/src/docs"
    "github.com/go-playground/validator"
    "github.com/go-redis/redis/v8"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo-contrib/prometheus"
    "github.com/labstack/echo/v4/middleware"
    "github.com/swaggo/echo-swagger"
    "gorm.io/driver/postgres"
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
    sqlClient   *gorm.DB
    redisCtx    context.Context
    redisClient *redis.Client
)

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func setupSql() {
    var err error

    host := os.Getenv("POSTGRES_HOST")
    port := os.Getenv("POSTGRES_PORT")
    user := os.Getenv("POSTGRES_USER")
    password := os.Getenv("POSTGRES_PASSWORD")
    dbname := os.Getenv("POSTGRES_DB")

    template := "host=%v port=%v user=%v password=%v dbname=%v sslmode=disable"
    dsn := fmt.Sprintf(template, host, port, user, password, dbname)
    sqlClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Could not connect to database.")
    }

    err = sqlClient.AutoMigrate(&Post{})
    if err != nil {
        panic("Could not migrate posts.")
    }

    err = sqlClient.AutoMigrate(&Comment{})
    if err != nil {
        panic("Could not migrate comments.")
    }

    err = sqlClient.AutoMigrate(&User{})
    if err != nil {
        panic("Could not migrate users.")
    }
}

func setupRedis() {
    host := os.Getenv("REDIS_HOST")
    port := os.Getenv("REDIS_PORT")
    password := os.Getenv("REDIS_PASSWORD")
    db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

    redisCtx = context.Background()

    redisClient = redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%v:%v", host, port),
        Password: password,
        DB:       db,
    })
}

// @title Simple blogging platform API based on PostgreSQL and Redis databases
// @version 1.0
// @description Simple blogging platform API created in Golang with the use of Echo framework, GORM ORM library,
// @description PostgreSQL database for storing objects, and Redis for storing session data with integrated Prometheus
// @description and Swagger, Dockerized.

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
