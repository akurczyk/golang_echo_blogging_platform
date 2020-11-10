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
    Comment struct {
        ID           uint           `json:"id" gorm:"primarykey"`
        CreatedAt    time.Time      `json:"created_at"`
        UpdatedAt    time.Time      `json:"updated_at"`
        DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
        AuthorID     uint           `json:"author_id"`
        Author       User           `json:"author"`
        PostID       int            `json:"post_id"`
        Post         Post           `json:"post"`
        Content      string         `json:"content"`
    }

    CommentCreate struct {
        PostID  int    `json:"post_id" validate:"required,numeric"`
        Content string `json:"content" validate:"required,min=3"`
    }

    CommentUpdate struct {
        Content string `json:"content" validate:"required,min=3"`
    }
)

func getCommentOrError(context echo.Context) (*Comment, int) {
    var comment Comment

    id, err := strconv.Atoi(context.Param("id"))
    if err != nil {
        return nil, http.StatusBadRequest
    }

    result := sql_db.First(&comment, id)
    if result.Error != nil {
        return nil, http.StatusNotFound
    }

    return &comment, 0
}

// listComments godoc
// @Summary List Comments
// @Tags comments
// @Produce json
// @Param author_id query int false "Author ID"
// @Param post_id query int false "Post ID"
// @Success 200 {array} Comment
// @Router /comments [get]
func listComments(context echo.Context) error {
    var comments []Comment

    query := sql_db
    if authorID := context.QueryParam("author_id"); authorID != "" {
        query = query.Where("AuthorID = ?", authorID)
    }
    if postID := context.QueryParam("post_id"); postID != "" {
        query = query.Where("PostID = ?", postID)
    }
    query.Find(&comments)

    return context.JSON(http.StatusOK, comments)
}

// createComment godoc
// @Summary Create Comment
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body CommentCreate true "Comment"
// @Security ApiKeyAuth
// @Success 201 {object} Comment
// @Failure 400
// @Failure 500
// @Router /comments [post]
func createComment(context echo.Context) error {
    user := context.Get("User").(User)

    commentCreate := new(CommentCreate)
    post := new(Post)
    comment := new(Comment)

    if err := context.Bind(commentCreate); err != nil {
        return err
    }

    if err := context.Validate(commentCreate); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    result := sql_db.First(&post, commentCreate.PostID)
    if result.Error != nil {
        return context.JSON(http.StatusBadRequest, "Provided post does not exists.")
    }

    comment.Author = user
    comment.Post = *post
    comment.Content = commentCreate.Content
    sql_db.Create(&comment)

    return context.JSON(http.StatusCreated, comment)
}

// retrieveComment godoc
// @Summary Retrieve Comment
// @Tags comments
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} Comment
// @Failure 400
// @Failure 404
// @Router /comments/{id} [get]
func retrieveComment(context echo.Context) error {
    comment, err := getCommentOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    return context.JSON(http.StatusOK, comment)
}

// updateComment godoc
// @Summary Update Comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param comment body CommentUpdate true "Comment"
// @Security ApiKeyAuth
// @Success 200 {object} Comment
// @Failure 400
// @Failure 403
// @Failure 404
// @Router /comments/{id} [put]
func updateComment(context echo.Context) error {
    user := context.Get("User").(User)

    comment, err := getCommentOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    if comment.AuthorID != user.ID {
        return context.NoContent(http.StatusForbidden)
    }

    commentUpdate := new(CommentUpdate)

    if err := context.Bind(commentUpdate); err != nil {
        return err
    }

    if err := context.Validate(commentUpdate); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    comment.Content = commentUpdate.Content
    sql_db.Save(&comment)

    return context.JSON(http.StatusOK, comment)
}

// deleteComment godoc
// @Summary Delete Comment
// @Tags comments
// @Param id path int true "ID"
// @Security ApiKeyAuth
// @Success 204
// @Failure 400
// @Failure 403
// @Failure 404
// @Router /comments/{id} [delete]
func deleteComment(context echo.Context) error {
    user := context.Get("User").(User)

    comment, err := getCommentOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    if comment.AuthorID != user.ID {
        return context.NoContent(http.StatusForbidden)
    }

    sql_db.Delete(&comment)

    return context.NoContent(http.StatusNoContent)
}
