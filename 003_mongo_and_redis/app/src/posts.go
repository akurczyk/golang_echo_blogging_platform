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
    Post struct {
        ID        primitive.ObjectID `bson:"_id" json:"_id"`
        CreatedAt time.Time          `bson:"created_at" json:"created_at"`
        UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
        Author    User               `bson:"author" json:"author"`
        Title     string             `bson:"title" json:"title"`
        Content   string             `bson:"content" json:"content"`
        Comments  []Comment          `bson:"comments" json:"comments"`
    }

    PostIn struct {
        Title   string `json:"title" validate:"required,min=3"`
        Content string `json:"content" validate:"required,min=3"`
    }
)

func getPostOrError(context echo.Context) (*Post, int) {
    var err error
    var id primitive.ObjectID
    var post Post

    id, err = primitive.ObjectIDFromHex(context.Param("id"))
    if err != nil {
        return nil, http.StatusBadRequest
    }

    filter := bson.M{
        "_id": id,
    }
    err = postsCollection.FindOne(mongoCtx, filter).Decode(&post)
    if err != nil {
        return nil, http.StatusNotFound
    }

    return &post, 0
}

// listPosts godoc
// @Summary List Posts
// @Tags posts
// @Produce json
// @Param author_id query string false "Author ID"
// @Success 200 {array} Post
// @Router /posts [get]
func listPosts(context echo.Context) error {
    var posts []Post
    var filter = bson.M{}

    if authorID := context.QueryParam("author_id"); authorID != "" {
        id, _ := primitive.ObjectIDFromHex(context.QueryParam("author_id"))
        filter = bson.M{
            "author._id": id,
        }
    }

    cursor, err := postsCollection.Find(mongoCtx, filter)
    if err != nil {
        panic(err)
    }
    defer cursor.Close(mongoCtx)

    for cursor.Next(mongoCtx) {
        var post Post

        err := cursor.Decode(&post)
        if err != nil {
            panic(err)
        }

        posts = append(posts, post)
    }

    return context.JSON(http.StatusOK, posts)
}

// createPost godoc
// @Summary Create Post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body PostIn true "Post"
// @Security ApiKeyAuth
// @Success 201 {object} Post
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /posts [post]
func createPost(context echo.Context) error {
    var err error

    postIn := new(PostIn)
    if err := context.Bind(postIn); err != nil {
        return err
    }
    if err := context.Validate(postIn); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    post := new(Post)
    post.ID = primitive.NewObjectID()
    post.CreatedAt = time.Now()
    post.UpdatedAt = time.Now()
    post.Author = context.Get("User").(User)
    post.Title = postIn.Title
    post.Content = postIn.Content
    post.Comments = []Comment{}

    _, err = postsCollection.InsertOne(mongoCtx, post)
    if err != nil {
        panic(err)
    }

    return context.JSON(http.StatusCreated, post)
}

// retrievePost godoc
// @Summary Retrieve Post
// @Tags posts
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} Post
// @Failure 400
// @Failure 404
// @Router /posts/{id} [get]
func retrievePost(context echo.Context) error {
    post, err := getPostOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    return context.JSON(http.StatusOK, post)
}

// updatePost godoc
// @Summary Update Post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param post body PostIn true "Post"
// @Security ApiKeyAuth
// @Success 200 {object} Post
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /posts/{id} [put]
func updatePost(context echo.Context) error {
    post, code := getPostOrError(context)
    if code != 0 {
        return context.NoContent(code)
    }

    if post.Author.ID != context.Get("User").(User).ID {
        return context.NoContent(http.StatusForbidden)
    }

    postIn := new(PostIn)
    if err := context.Bind(postIn); err != nil {
        return err
    }
    if err := context.Validate(postIn); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    post.Title = postIn.Title
    post.Content = postIn.Content

    filter := bson.M{
        "_id": post.ID,
    }
    update := bson.M{
        "$set": post,
    }
    _, err := postsCollection.UpdateOne(mongoCtx, filter, update)
    if err != nil {
        panic(err)
    }

    return context.JSON(http.StatusOK, post)
}

// deletePost godoc
// @Summary Delete Post
// @Tags posts
// @Param id path string true "ID"
// @Security ApiKeyAuth
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /posts/{id} [delete]
func deletePost(context echo.Context) error {
    post, code := getPostOrError(context)
    if code != 0 {
        return context.NoContent(code)
    }

    if post.Author.ID != context.Get("User").(User).ID {
        return context.NoContent(http.StatusForbidden)
    }

    filter := bson.M{
        "_id": post.ID,
    }
    _, err := postsCollection.DeleteOne(mongoCtx, filter)
    if err != nil {
        panic(err)
    }

    return context.NoContent(http.StatusNoContent)
}
