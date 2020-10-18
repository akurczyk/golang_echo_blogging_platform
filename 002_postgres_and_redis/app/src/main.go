package main

import (
    "fmt"
    _ "github.com/akurczyk/golang_echo_blogging_platform/002_postgres_and_redis/app/src/docs"
    "github.com/go-playground/validator"
    "github.com/labstack/echo"
    "github.com/labstack/echo-contrib/prometheus"
    "github.com/labstack/echo/middleware"
    "github.com/swaggo/echo-swagger"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "os"
)

type (
    CustomValidator struct {
        validator *validator.Validate
    }
)

var (
    db *gorm.DB
)

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func setupDb() {
    var err error

    host := os.Getenv("POSTGRES_HOST")
    port := os.Getenv("POSTGRES_PORT")
    user := os.Getenv("POSTGRES_USER")
    password := os.Getenv("POSTGRES_PASSWORD")
    dbname := os.Getenv("POSTGRES_DB")

    template := "host=%v port=%v user=%v password=%v dbname=%v sslmode=disable"
    dsn := fmt.Sprintf(template, host, port, user, password, dbname)
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Could not connect to database.")
    }

    err = db.AutoMigrate(&Post{})
    if err != nil {
        panic("Could not migrate posts.")
    }

    err = db.AutoMigrate(&Comment{})
    if err != nil {
        panic("Could not migrate comments.")
    }

    err = db.AutoMigrate(&User{})
    if err != nil {
        panic("Could not migrate users.")
    }
}

// @title Simple blogging platform API based on PostgreSQL and Redis databases
// @version 1.0
// @description Simple blogging platform API created in Golang with the use of Echo framework, GORM ORM library,
// @description PostgreSQL database for storing objects, and Redis for storing session data with integrated Prometheus
// @description and Swagger dockerized.

// @contact.name Aleksander Kurczyk
// @contact.url http://github.com/akurczyk
// @contact.email akurczyk000@gmail.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
    setupDb()

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
// TODO: SET_NULL to Author
// TODO: Redis
// TODO: PostgreSQL start problem
