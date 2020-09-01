package main

import (
    "github.com/labstack/echo"
    "net/http"
    "strings"
)

type (
    user struct {
        Name         string `json:"name"`
        PasswordHash string `json:"-"`
        Email        string `json:"email"`
    }

    userNew struct {
        Name     string `json:"name" validate:"required,min=3"`
        Password string `json:"password" validate:"required,min=6"`
        Email    string `json:"email" validate:"required,email"`
    }

    userUpdate struct {
        Password string `json:"password" validate:"required,min=6"`
        Email    string `json:"email" validate:"required,email"`
    }
)

var (
    users = map[string]*user{}
)

func getUserOrError(context echo.Context) (*user, error) {
    id := context.Param("id")

    user, ok := users[id]
    if !ok {
        return nil, context.NoContent(http.StatusNotFound)
    }

    return user, nil
}

func listUserAccounts(context echo.Context) error {
    return context.JSON(http.StatusOK, users)
}

func createUserAccount(context echo.Context) error {
    user := new(user)
    userNew := new(userNew)

    if err := context.Bind(userNew); err != nil {
        return err
    }

    if err := context.Validate(userNew); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    if _, ok := users[userNew.Name]; ok {
        return context.JSON(http.StatusBadRequest, "User with provided name already exists.")
    }

    hashedPassword, err := hashAndSalt(userNew.Password)
    if err != nil {
        return err
    }

    user.Name = userNew.Name
    user.PasswordHash = hashedPassword
    user.Email = userNew.Email
    users[user.Name] = user

    return context.JSON(http.StatusCreated, user)
}

func retrieveUserAccount(context echo.Context) error {
    user, err := getUserOrError(context)
    if err != nil {
        return err
    }

    return context.JSON(http.StatusOK, user)
}

func updateUserAccount(context echo.Context) error {
    user := context.Get("user").(*user)

    userUpdate := new(userUpdate)

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
    users[user.Name] = user

    return context.JSON(http.StatusOK, user)
}

func deleteUserAccount(context echo.Context) error {
    user := context.Get("user").(*user)
    delete(users, user.Name)
    return context.NoContent(http.StatusNoContent)
}
