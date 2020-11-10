package main

import (
    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type (
    User struct {
        ID           uint           `json:"id" gorm:"primarykey"`
        CreatedAt    time.Time      `json:"created_at"`
        UpdatedAt    time.Time      `json:"updated_at"`
        DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
        Name         string         `json:"name" gorm:"uniqueIndex;size:255"`
        PasswordHash string         `json:"-" gorm:"size:255"`
        Email        string         `json:"email" gorm:"uniqueIndex;size:255"`
    }

    UserNew struct {
        Name     string `json:"name" validate:"required,min=3"`
        Password string `json:"password" validate:"required,min=6"`
        Email    string `json:"email" validate:"required,email"`
    }

    UserUpdate struct {
        Password string `json:"password" validate:"required,min=6"`
        Email    string `json:"email" validate:"required,email"`
    }
)

func getUserOrError(context echo.Context) (*User, int) {
    var user User

    id, err := strconv.Atoi(context.Param("id"))
    if err != nil {
        return nil, http.StatusBadRequest
    }

    result := sql_db.First(&user, id)
    if result.Error != nil {
        return nil, http.StatusNotFound
    }

    return &user, 0
}

// listUserAccounts godoc
// @Summary List Users Accounts
// @Tags users
// @Produce json
// @Param name query string false "Name"
// @Param email query string false "Email"
// @Success 200 {array} User
// @Router /users [get]
func listUserAccounts(context echo.Context) error {
    var users []User

    query := sql_db
    if name := context.QueryParam("name"); name != "" {
        query = query.Where("Name = ?", name)
    }
    if email := context.QueryParam("email"); email != "" {
        query = query.Where("Email = ?", email)
    }
    query.Find(&users)

    return context.JSON(http.StatusOK, users)
}

// createUserAccount godoc
// @Summary Create User Account
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserNew true "User"
// @Success 201 {object} User
// @Failure 400
// @Failure 500
// @Router /users [post]
func createUserAccount(context echo.Context) error {
    user := new(User)
    userNew := new(UserNew)

    if err := context.Bind(userNew); err != nil {
        return err
    }

    if err := context.Validate(userNew); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    hashedPassword, err := hashAndSalt(userNew.Password)
    if err != nil {
        return err
    }

    user.Name = userNew.Name
    user.PasswordHash = hashedPassword
    user.Email = userNew.Email

    result := sql_db.Create(&user)
    if result.Error != nil {
        return context.JSON(http.StatusBadRequest, "User with provided name already exists.")
    }

    return context.JSON(http.StatusCreated, user)
}

// retrieveUserAccount godoc
// @Summary Retrieve User Account
// @Tags users
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} User
// @Failure 400
// @Failure 404
// @Router /users/{id} [get]
func retrieveUserAccount(context echo.Context) error {
    user, err := getUserOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    return context.JSON(http.StatusOK, user)
}

// updateUserAccount godoc
// @Summary Update Current User Account
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body UserUpdate true "User"
// @Success 200 {object} User
// @Failure 400
// @Failure 403
// @Failure 404
// @Router /users [put]
func updateUserAccount(context echo.Context) error {
    user := context.Get("User").(User)

    userUpdate := new(UserUpdate)

    if err := context.Bind(userUpdate); err != nil {
        return err
    }

    if err := context.Validate(userUpdate); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    hashedPassword, err := hashAndSalt(userUpdate.Password)
    if err != nil {
        return err
    }

    user.PasswordHash = hashedPassword
    user.Email = userUpdate.Email
    sql_db.Save(&user)

    return context.JSON(http.StatusOK, user)
}

// deleteUserAccount godoc
// @Summary Delete Current User Account
// @Tags users
// @Security ApiKeyAuth
// @Success 204
// @Success 400
// @Failure 403
// @Failure 404
// @Router /users [delete]
func deleteUserAccount(context echo.Context) error {
    user := context.Get("User").(User)

    sql_db.Delete(&user)

    return context.NoContent(http.StatusNoContent)
}
