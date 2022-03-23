package article

import (
	"gogin/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var log = logger.New("article")

// ArticleCtx middleware is used to load an Article object from
// the URL parameters passed through as the request. In case
// the Article could not be found, we stop here and return a 404.
func ArticleCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		var article *Article
		var err error

		if articlePath := c.Param("articlePath"); articlePath != "" {
			article, err = DBGetArticle(articlePath)
			if article == nil {
				article, err = DBGetArticleBySlug(articlePath)
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			c.Abort()
			return
		}
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			c.Abort()
			return
		}
		// set article to the gin context
		c.Set("article", article)
		c.Next()
	}
}

// --
// CRUD Operations
// --

// CreateArticle persists the posted Article and returns it
// back to the client as an acknowledgement.
func CreateArticle(c *gin.Context) {
	article := &Article{}
	// using BindJson method to serialize body with struct
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	article.Title = strings.ToLower(article.Title)
	article.Slug = strings.ReplaceAll(article.Title, " ", "-")
	DBNewArticle(article)

	c.JSON(http.StatusCreated, gin.H{"articles": &article})
}

// ReadArticle returns the specific Article. You'll notice it just
// fetches the Article right off the context, as its understood that
// if we made it this far, the Article must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func ReadArticle(c *gin.Context) {
	data, _ := c.Get("article")
	article := data.(*Article) // type assertion
	c.JSON(http.StatusOK, gin.H{"articles": article})
}

// UpdateArticle updates an existing Article in our persistent store.
func UpdateArticle(c *gin.Context) {
	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	data, _ := c.Get("article")
	oldArticle := data.(*Article) // type assertion
	article := &Article{}
	// using BindJson method to serialize body with struct
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	article.Title = strings.ToLower(article.Title)
	article.Slug = strings.ReplaceAll(article.Title, " ", "-")

	DBUpdateArticle(oldArticle.ID, article)
	c.JSON(http.StatusOK, gin.H{"articles": article})
}

// DeleteArticle removes an existing Article from our persistent store.
func DeleteArticle(c *gin.Context) {
	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	data, _ := c.Get("article")
	article := data.(*Article) // type assertion
	DBRemoveArticle(article.ID)
	c.JSON(http.StatusOK, gin.H{"articles": article})
}

// ListArticles returns all articles from our persistent store.
func ListArticles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"articles": articles})
}
