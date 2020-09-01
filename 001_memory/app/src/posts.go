package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	post struct {
		ID				int			`json:"id"`
		AuthorName		string		`json:"author_name"`
		Title			string		`json:"title"`
		Content			string		`json:"content"`
		CreateDate		time.Time	`json:"create_date"`
		ModifyDate		time.Time	`json:"modify_date"`
	}

	postIn struct {
		Title			string		`json:"title" validate:"required,min=3"`
		Content			string		`json:"content" validate:"required,min=3"`
	}
)

var (
	posts		= map[int]*post{}
	postsSeq	= 1
)

func getPostOrError(context echo.Context) (*post, error) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		return nil, context.NoContent(http.StatusBadRequest)
	}

	post, ok := posts[id]
	if !ok {
		return nil, context.NoContent(http.StatusNotFound)
	}

	return post, nil
}

func listPosts(context echo.Context) error {
	if authorName := context.QueryParam("author_name"); authorName != "" {
		filteredPosts := map[int]*post{}
		for id, post := range posts {
			if post.AuthorName == authorName {
				filteredPosts[id] = post
			}
		}

		return context.JSON(http.StatusOK, filteredPosts)
	}

	return context.JSON(http.StatusOK, posts)
}

func createPost(context echo.Context) error {
	user := context.Get("user").(*user)

	post := new(post)
	postIn := new(postIn)

	if err := context.Bind(postIn); err != nil {
		return err
	}

	if err := context.Validate(postIn); err != nil {
		return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
	}

	post.ID = postsSeq
	post.AuthorName = user.Name
	post.Title = postIn.Title
	post.Content = postIn.Content
	post.CreateDate = time.Now()
	post.ModifyDate = time.Now()
	posts[post.ID] = post
	postsSeq++

	return context.JSON(http.StatusCreated, post)
}

func retrievePost(context echo.Context) error {
	post, err := getPostOrError(context)
	if err != nil {
		return err
	}

	return context.JSON(http.StatusOK, post)
}

func updatePost(context echo.Context) error {
	user := context.Get("user").(*user)

	post, err := getPostOrError(context)
	if err != nil {
		return err
	}

	if post.AuthorName != user.Name {
		return context.NoContent(http.StatusForbidden)
	}

	postIn := new(postIn)

	if err := context.Bind(postIn); err != nil {
		return err
	}

	if err := context.Validate(postIn); err != nil {
		return context.JSON(http.StatusBadRequest, strings.Split(err.Error(), "\n"))
	}

	post.Title = postIn.Title
	post.Content = postIn.Content
	post.ModifyDate = time.Now()
	posts[post.ID] = post

	return context.JSON(http.StatusOK, post)
}

func deletePost(context echo.Context) error {
	user := context.Get("user").(*user)

	post, err := getPostOrError(context)
	if err != nil {
		return err
	}

	if post.AuthorName != user.Name {
		return context.NoContent(http.StatusForbidden)
	}

	delete(posts, post.ID)

	return context.NoContent(http.StatusNoContent)
}
