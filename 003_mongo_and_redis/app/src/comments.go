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
    Comment struct {
        ID        primitive.ObjectID `bson:"_id" json:"_id"`
        CreatedAt time.Time          `bson:"created_at" json:"created_at"`
        UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
        Author    User               `bson:"author" json:"author"`
        Content   string             `bson:"content" json:"content"`
    }

    CommentIn struct {
        Content string `json:"content" validate:"required,min=3"`
    }
)

// createComment godoc
// @Summary Create Comment
// @Tags comments
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param comment body CommentIn true "Comment"
// @Security ApiKeyAuth
// @Success 201 {object} Comment
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /posts/{post_id}/comments [post]
func createComment(context echo.Context) error {
    // Validate data and prepare final Comment object
    commentIn := new(CommentIn)
    if err := context.Bind(commentIn); err != nil {
        return err
    }
    if err := context.Validate(commentIn); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    comment := new(Comment)
    comment.ID = primitive.NewObjectID()
    comment.CreatedAt = time.Now()
    comment.UpdatedAt = time.Now()
    comment.Author = context.Get("User").(User)
    comment.Content = commentIn.Content

    // Retrieve Post ID
    postID, err := primitive.ObjectIDFromHex(context.Param("post_id"))
    if err != nil {
        return context.NoContent(http.StatusBadRequest)
    }

    // Prepare and execute query
    filter := bson.M{
        "_id": postID,
    }
    update := bson.M{
        "$push": bson.M{
            "comments": comment,
        },
    }
    result, err := postsCollection.UpdateOne(mongoCtx, filter, update)
    if err != nil {
        panic(err)
    } else if result.MatchedCount != 1 {
        return context.JSON(http.StatusNotFound, comment)
    } else {
        return context.JSON(http.StatusCreated, comment)
    }
}

// updateComment godoc
// @Summary Update Comment
// @Tags comments
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param comment_id path string true "Comment ID"
// @Param comment body CommentIn true "Comment"
// @Security ApiKeyAuth
// @Success 200 {object} Comment
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /posts/{post_id}/comments/{comment_id} [put]
func updateComment(context echo.Context) error {
    // Validate data and prepare final Comment object
    commentIn := new(CommentIn)
    comment := new(Comment)

    if err := context.Bind(commentIn); err != nil {
        return err
    }

    if err := context.Validate(commentIn); err != nil {
        return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
    }

    comment.UpdatedAt = time.Now()
    comment.Content = commentIn.Content

    // Retrieve Post ID and Comment ID
    postID, err := primitive.ObjectIDFromHex(context.Param("post_id"))
    if err != nil {
        return context.NoContent(http.StatusBadRequest)
    }

    commentID, err := primitive.ObjectIDFromHex(context.Param("comment_id"))
    if err != nil {
        return context.NoContent(http.StatusBadRequest)
    }

    // Prepare and execute query
    filter := bson.M{
        "_id": postID,
        "comments._id": commentID,
    }
    update := bson.M{
        "$set": bson.M{
            "comments.$.updated_at": comment.UpdatedAt,
            "comments.$.content": comment.Content,
        },
    }
    result, err := postsCollection.UpdateOne(mongoCtx, filter, update)
    if err != nil {
        panic(err)
    } else if result.MatchedCount != 1 {
        return context.NoContent(http.StatusNotFound)
    } else {
        return context.NoContent(http.StatusOK)
    }
}

// deleteComment godoc
// @Summary Delete Comment
// @Tags comments
// @Param post_id path string true "Post ID"
// @Param comment_id path string true "Comment ID"
// @Security ApiKeyAuth
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 404
// @Router /posts/{post_id}/comments/{comment_id} [delete]
func deleteComment(context echo.Context) error {
    // Retrieve Post ID and Comment ID
    postID, err := primitive.ObjectIDFromHex(context.Param("post_id"))
    if err != nil {
        return context.NoContent(http.StatusBadRequest)
    }

    commentID, err := primitive.ObjectIDFromHex(context.Param("comment_id"))
    if err != nil {
        return context.NoContent(http.StatusBadRequest)
    }

    // Prepare and execute query
    filter := bson.M{
        "_id": postID,
    }
    update := bson.M{
        "$pull": bson.M{
            "comments": bson.M{"_id": commentID},
        },
    }
    result, err := postsCollection.UpdateOne(mongoCtx, filter, update)
    if err != nil {
        panic(err)
    } else if result.MatchedCount != 1 || result.ModifiedCount != 1 {
        return context.NoContent(http.StatusNotFound)
    } else {
        return context.NoContent(http.StatusNoContent)
    }
}
