package main

import (
    "github.com/go-playground/validator"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
)

type (
    CustomValidator struct {
        validator *validator.Validate
    }
)

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func main() {
    e := echo.New()

    // Validator
    e.Validator = &CustomValidator{validator: validator.New()}

    e.GET("/users", listUserAccounts)
    e.POST("/users", createUserAccount)
    e.GET("/users/:name", retrieveUserAccount)
    e.PUT("/users", updateUserAccount, middleware.KeyAuth(checkAuthToken))
    e.DELETE("/users", deleteUserAccount, middleware.KeyAuth(checkAuthToken))

    e.POST("/token", issueToken)
    e.DELETE("/token", revokeToken, middleware.KeyAuth(checkAuthToken))

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
