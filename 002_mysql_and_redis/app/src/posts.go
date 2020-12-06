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
    Post struct {
        ID         uint           `json:"id" gorm:"primarykey"`
        CreatedAt  time.Time      `json:"created_at"`
        UpdatedAt  time.Time      `json:"updated_at"`
        DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
        AuthorID   uint           `json:"author_id"`
        Author     User           `json:"author"`
        Title      string         `json:"title" gorm:"size:255"`
        Content    string         `json:"content"`
    }

    PostIn struct {
        Title   string `json:"title" validate:"required,min=3"`
        Content string `json:"content" validate:"required,min=3"`
    }
)

func getPostOrError(context echo.Context) (*Post, int) {
    var post Post

    id, err := strconv.Atoi(context.Param("id"))
    if err != nil {
        return nil, http.StatusBadRequest
    }

    result := sqlClient.First(&post, id)
    if result.Error != nil {
        return nil, http.StatusNotFound
    }

    return &post, 0
}

// listPosts godoc
// @Summary List Posts
// @Tags posts
// @Produce json
// @Param author_id query int false "Author ID"
// @Success 200 {array} Post
// @Router /posts [get]
func listPosts(context echo.Context) error {
    var posts []Post

    query := sqlClient
    if authorID := context.QueryParam("author_id"); authorID != "" {
        query = query.Where("AuthorID = ?", authorID)
    }
    query.Find(&posts)

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
    postIn := new(PostIn)
    if err := context.Bind(postIn); err != nil {
        return err
    }
    if err := context.Validate(postIn); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    post := new(Post)
    post.Author = context.Get("User").(User)
    post.Title = postIn.Title
    post.Content = postIn.Content

    result := sqlClient.Create(&post)
    if result.Error != nil {
        return result.Error
    }

    return context.JSON(http.StatusCreated, post)
}

// retrievePost godoc
// @Summary Retrieve Post
// @Tags posts
// @Produce json
// @Param id path int true "ID"
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
// @Param id path int true "ID"
// @Param post body PostIn true "Post"
// @Security ApiKeyAuth
// @Success 200 {object} Post
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /posts/{id} [put]
func updatePost(context echo.Context) error {
    post, err := getPostOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    if post.AuthorID != context.Get("User").(User).ID {
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
    sqlClient.Save(&post)

    return context.JSON(http.StatusOK, post)
}

// deletePost godoc
// @Summary Delete Post
// @Tags posts
// @Param id path int true "ID"
// @Security ApiKeyAuth
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /posts/{id} [delete]
func deletePost(context echo.Context) error {
    post, err := getPostOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    if post.AuthorID != context.Get("User").(User).ID {
        return context.NoContent(http.StatusForbidden)
    }

    sqlClient.Delete(&post)

    return context.NoContent(http.StatusNoContent)
}
