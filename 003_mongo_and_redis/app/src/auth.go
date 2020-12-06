package main

import (
    "encoding/json"
    "github.com/labstack/echo/v4"
    "github.com/labstack/gommon/random"
    "go.mongodb.org/mongo-driver/bson"
    "golang.org/x/crypto/bcrypt"
    "net/http"
    "time"
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
    var userObj User
    var userJson string
    var err error

    userJson, err = redisClient.Get(redisCtx, token).Result()
    if err != nil {
        return false, nil
    }

    _ = redisClient.Expire(redisCtx, token, 1 * time.Hour).Err()

    err = json.Unmarshal([]byte(userJson), &userObj)
    if err != nil {
        panic(err)
    }

    context.Set("token", token)
    context.Set("User", userObj)

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
    var userObj User
    var userJson []byte
    var err error

    name := context.FormValue("name")
    filter := bson.D{{"name", name}}
    err = mongoDatabase.Collection("users").FindOne(mongoCtx, filter).Decode(&userObj)
    if err != nil {
        return context.NoContent(http.StatusUnauthorized)
    }

    if !comparePasswords(userObj.PasswordHash, context.FormValue("password")) {
        return context.NoContent(http.StatusUnauthorized)
    }

    token := random.String(32, random.Alphanumeric)

    userJson, err = json.Marshal(userObj)
    if err != nil {
        panic(err)
    }

    err = redisClient.Set(redisCtx, token, string(userJson), 1 * time.Hour).Err()
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
    redisClient.Del(redisCtx, token)
    return context.NoContent(http.StatusNoContent)
}
