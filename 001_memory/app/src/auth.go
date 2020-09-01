package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/random"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var (
	tokens = map[string]*user{}
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
	// Check whether provided token exists
	if _, ok := tokens[token]; !ok {
		return false, nil
	}

	// Check whether user account associated with the token has been deleted
	if _, ok := users[tokens[token].Name]; !ok {
		delete(tokens, token)
		return false, nil
	}

	context.Set("token", token)
	context.Set("user", tokens[token])
	return true, nil
}

func issueToken(context echo.Context) error {
	user, ok := users[context.FormValue("name")]
	if !ok {
		return context.NoContent(http.StatusUnauthorized)
	}

	if !comparePasswords(user.PasswordHash, context.FormValue("password")) {
		return context.NoContent(http.StatusUnauthorized)
	}

	token := random.String(32, random.Alphanumeric)
	tokens[token] = user

	return context.String(http.StatusOK, token)
}

func revokeToken(context echo.Context) error {
	token := context.Get("token").(string)
	delete(tokens, token)
	return context.NoContent(http.StatusNoContent)
}
