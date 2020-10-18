package main

import (
    "github.com/labstack/echo"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type (
    Post struct {
        ID         int       `json:"id" gorm:"primarykey"`
        AuthorID   int       `json:"author_id"`
        Author     User      `json:"author"`
        Title      string    `json:"title"`
        Content    string    `json:"content"`
        CreateDate time.Time `json:"create_date" gorm:"autoCreateTime"`
        ModifyDate time.Time `json:"modify_date" form:"autoUpdateTime"`
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

    result := db.First(&post, id)
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

    query := db
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
// @Success 201 {object} Post
// @Failure 400
// @Failure 500
// @Router /posts [post]
func createPost(context echo.Context) error {
    user := context.Get("User").(*User)

    post := new(Post)
    postIn := new(PostIn)

    if err := context.Bind(postIn); err != nil {
        return err
    }

    if err := context.Validate(postIn); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    post.AuthorID = user.ID
    post.Title = postIn.Title
    post.Content = postIn.Content

    result := db.Create(&post)
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
// @Success 200 {object} Post
// @Failure 400
// @Failure 403
// @Failure 404
// @Router /posts/{id} [put]
func updatePost(context echo.Context) error {
    user := context.Get("User").(*User)

    post, err := getPostOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    if post.AuthorID != user.ID {
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
    db.Save(&post)

    return context.JSON(http.StatusOK, post)
}

// deletePost godoc
// @Summary Delete Post
// @Tags posts
// @Param id path int true "ID"
// @Success 204
// @Failure 400
// @Failure 403
// @Failure 404
// @Router /posts/{id} [delete]
func deletePost(context echo.Context) error {
    user := context.Get("User").(*User)

    post, err := getPostOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    if post.AuthorID != user.ID {
        return context.NoContent(http.StatusForbidden)
    }

    db.Delete(&post)

    return context.NoContent(http.StatusNoContent)
}
