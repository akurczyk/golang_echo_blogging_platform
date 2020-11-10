package main

import (
    "github.com/labstack/echo"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type (
    comment struct {
        ID         int       `json:"id"`
        AuthorName string    `json:"author_name"`
        PostID     int       `json:"post_id"`
        Content    string    `json:"content"`
        CreateDate time.Time `json:"create_date"`
        ModifyDate time.Time `json:"modify_date"`
    }

    commentCreate struct {
        PostID  int    `json:"post_id" validate:"required,numeric"`
        Content string `json:"content" validate:"required,min=3"`
    }

    commentUpdate struct {
        Content string `json:"content" validate:"required,min=3"`
    }
)

var (
    comments    = map[int]*comment{}
    commentsSeq = 1
)

func getCommentOrError(context echo.Context) (*comment, int) {
    id, err := strconv.Atoi(context.Param("id"))
    if err != nil {
        return nil, http.StatusBadRequest
    }

    comment, ok := comments[id]
    if !ok {
        return nil, http.StatusNotFound
    }

    return comment, 0
}

func listComments(context echo.Context) error {
    authorName := context.QueryParam("author_name")
    postID := context.QueryParam("post_id")

    if authorName != "" || postID != "" {
        numericPostID, err := strconv.Atoi(postID)
        if err != nil {
            numericPostID = -1
        }

        filteredComments := map[int]*comment{}
        for id, comment := range comments {
            if comment.AuthorName == authorName || comment.PostID == numericPostID {
                filteredComments[id] = comment
            }
        }

        return context.JSON(http.StatusOK, filteredComments)
    }

    return context.JSON(http.StatusOK, comments)
}

func createComment(context echo.Context) error {
    user := context.Get("user").(*user)

    comment := new(comment)
    commentCreate := new(commentCreate)

    if err := context.Bind(commentCreate); err != nil {
        return err
    }

    if err := context.Validate(commentCreate); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    comment.ID = commentsSeq
    comment.AuthorName = user.Name
    comment.PostID = commentCreate.PostID
    comment.Content = commentCreate.Content
    comment.CreateDate = time.Now()
    comment.ModifyDate = time.Now()
    comments[comment.ID] = comment
    commentsSeq++

    return context.JSON(http.StatusCreated, comment)
}

func retrieveComment(context echo.Context) error {
    comment, err := getCommentOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    return context.JSON(http.StatusOK, comment)
}

func updateComment(context echo.Context) error {
    user := context.Get("user").(*user)

    comment, err := getCommentOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    if comment.AuthorName != user.Name {
        return context.NoContent(http.StatusForbidden)
    }

    commentUpdate := new(commentUpdate)

    if err := context.Bind(commentUpdate); err != nil {
        return err
    }

    if err := context.Validate(commentUpdate); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    comment.Content = commentUpdate.Content
    comment.ModifyDate = time.Now()
    comments[comment.ID] = comment

    return context.JSON(http.StatusOK, comment)
}

func deleteComment(context echo.Context) error {
    user := context.Get("user").(*user)

    comment, err := getCommentOrError(context)
    if err != 0 {
        return context.NoContent(err)
    }

    if comment.AuthorName != user.Name {
        return context.NoContent(http.StatusForbidden)
    }

    delete(comments, comment.ID)

    return context.NoContent(http.StatusNoContent)
}
