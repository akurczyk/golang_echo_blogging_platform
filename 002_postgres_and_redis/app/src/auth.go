package main

import (
    "github.com/labstack/echo"
    "github.com/labstack/gommon/random"
    "golang.org/x/crypto/bcrypt"
    "net/http"
)

var (
    tokens = map[string]User{}
)

func hashAndSalt(pwd string) (string, error) {
    bytePwd := []byte(pwd)
    hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
    byteHashedPwd := []byte(hashedPwd)
    bytePlainPwd := []byte(plainPwd)
    err := bcrypt.CompareHashAndPassword(byteHashedPwd, bytePlainPwd)
    return err == nil
}

func checkAuthToken(token string, context echo.Context) (bool, error) {
    if _, ok := tokens[token]; !ok {
        return false, nil
    }

    context.Set("token", token)
    context.Set("User", tokens[token])
    return true, nil
}

// issueAuthToken godoc
// @Summary Issue Auth Token / Login
// @Tags auth
// @Accept mpfd
// @Produce plain
// @Param name formData string true "Name"
// @Param password formData string true "Password"
// @Success 201
// @Failure 401
// @Failure 500
// @Router /token [post]
func issueAuthToken(context echo.Context) error {
    var user User

    result := db.First(&user, context.FormValue("name"))
    if result.Error != nil {
        return context.NoContent(http.StatusUnauthorized)
    }

    if !comparePasswords(user.PasswordHash, context.FormValue("password")) {
        return context.NoContent(http.StatusUnauthorized)
    }

    token := random.String(32, random.Alphanumeric)
    tokens[token] = user

    return context.String(http.StatusCreated, token)
}

// revokeAuthToken godoc
// @Summary Revoke Auth Token / Logout Current User
// @Tags auth
// @Success 204
// @Router /token [delete]
func revokeAuthToken(context echo.Context) error {
    token := context.Get("token").(string)
    delete(tokens, token)
    return context.NoContent(http.StatusNoContent)
}
