package main

import (
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "net/http"
    "strings"
    "time"
)

type (
    User struct {
        ID           primitive.ObjectID `bson:"_id" json:"_id"`
        CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
        UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
        Name         string             `bson:"name" json:"name"`
        PasswordHash string             `bson:"password_hash" json:"-"`
        Email        string             `bson:"email" json:"email"`
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
    var err error
    var id primitive.ObjectID
    var user User

    id, err = primitive.ObjectIDFromHex(context.Param("id"))
    if err != nil {
        return nil, http.StatusBadRequest
    }

    filter := bson.M{
        "_id": id,
    }
    err = usersCollection.FindOne(mongoCtx, filter).Decode(&user)
    if err != nil {
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
    var filters []bson.M
    var filter bson.M

    if name := context.QueryParam("name"); name != "" {
        filters = append(filters, bson.M{
            "name": name,
        })
    }
    if email := context.QueryParam("email"); email != "" {
        filters = append(filters, bson.M{
            "email": email,
        })
    }
    if len(filters) > 0 {
        filter = bson.M{
            "$and": filters,
        }
    } else {
        filter = bson.M{}
    }

    cursor, err := usersCollection.Find(mongoCtx, filter)
    if err != nil {
        panic(err)
    }
    defer cursor.Close(mongoCtx)

    for cursor.Next(mongoCtx) {
        var user User

        err := cursor.Decode(&user)
        if err != nil {
            panic(err)
        }

        users = append(users, user)
    }

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
// @Failure 401
// @Failure 500
// @Router /users [post]
func createUserAccount(context echo.Context) error {
    var err error

    userNew := new(UserNew)
    if err = context.Bind(userNew); err != nil {
        return err
    }
    if err = context.Validate(userNew); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    hashedPassword, err := hashAndSalt(userNew.Password)
    if err != nil {
        return err
    }

    user := new(User)
    user.ID = primitive.NewObjectID()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    user.Name = userNew.Name
    user.PasswordHash = hashedPassword
    user.Email = userNew.Email

    _, err = usersCollection.InsertOne(mongoCtx, user)
    if err != nil {
        if strings.Contains(err.Error(), "Duplicate key error.") {
            return context.JSON(http.StatusBadRequest, "User with provided name or email already exists.")
        }
        panic(err)
    }

    return context.JSON(http.StatusCreated, user)
}

// retrieveUserAccount godoc
// @Summary Retrieve User Account
// @Tags users
// @Produce json
// @Param id path string true "ID"
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

    user := context.Get("User").(User)
    user.PasswordHash = hashedPassword
    user.Email = userUpdate.Email

    filter := bson.M{
        "_id": user.ID,
    }
    update := bson.M{
        "$set": user,
    }
    _, err = usersCollection.UpdateOne(mongoCtx, filter, update)
    if err != nil {
        panic(err)
    }

    return context.JSON(http.StatusOK, user)
}

// deleteUserAccount godoc
// @Summary Delete Current User Account
// @Tags users
// @Security ApiKeyAuth
// @Success 204
// @Success 401
// @Router /users [delete]
func deleteUserAccount(context echo.Context) error {
    user := context.Get("User").(User)

    filter := bson.M{
        "_id": user.ID,
    }
    _, err := usersCollection.DeleteOne(mongoCtx, filter)
    if err != nil {
        panic(err)
    }

    // TODO: Delete user sessions here

    return context.NoContent(http.StatusNoContent)
}
