package main

import (
    "encoding/json"
    "github.com/labstack/echo/v4"
    "github.com/labstack/gommon/random"
    "golang.org/x/crypto/bcrypt"
    "net/http"
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
    var user_obj User
    var user_json string
    var err error

    user_json, err = redis_db.Get(ctx, token).Result()
    if err != nil {
        return false, nil
    }

    err = json.Unmarshal([]byte(user_json), &user_obj)
    if err != nil {
        panic(err)
    }

    context.Set("token", token)
    context.Set("User", user_obj)

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
    var user_obj User
    var user_json []byte
    var err error

    result := db.First(&user_obj, "Name = ?", context.FormValue("name"))
    if result.Error != nil {
        return context.NoContent(http.StatusUnauthorized)
    }

    if !comparePasswords(user_obj.PasswordHash, context.FormValue("password")) {
        return context.NoContent(http.StatusUnauthorized)
    }

    token := random.String(32, random.Alphanumeric)

    user_json, err = json.Marshal(user_obj)
    if err != nil {
        panic(err)
    }

    err = redis_db.Set(ctx, token, string(user_json), 0).Err()
    if err != nil {
        panic(err)
    }

    return context.String(http.StatusCreated, token)
}

// revokeAuthToken godoc
// @Summary Revoke Auth Token / Logout Current User
// @Tags auth
// @Security ApiKeyAuth
// @Success 204
// @Router /token [delete]
func revokeAuthToken(context echo.Context) error {
    token := context.Get("token").(string)
    redis_db.Del(ctx, token)
    return context.NoContent(http.StatusNoContent)
}
