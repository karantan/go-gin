package main

import (
	"gogin/article"
	"gogin/database"
	"gogin/logger"
	"gogin/user"

	"github.com/gin-gonic/gin"
)

//
// REST
// ====
// This example demonstrates a HTTP REST web service with some fixture data.
// Follow along the example and patterns.
//
//
// Boot the server:
// ----------------
// $ go run main.go
//
// Client requests:
// ----------------
// $ curl http://localhost:3000/
// {"message":"Hello World"}
//
// $ curl http://localhost:3000/api/v1/articles/
// {"articles":[{"id":"1","user_id":100,"title":"Hi","slug":"hi"},{"id":"2","user_id":200,"title":"sup","slug":"sup"},{"id":"3","user_id":300,"title":"alo","slug":"alo"},{"id":"4","user_id":400,"title":"bonjour","slug":"bonjour"},{"id":"5","user_id":500,"title":"whats up","slug":"whats-up"}]}
//
// $ curl http://localhost:3000/api/v1/articles/1/
// {"articles":{"id":"1","user_id":100,"title":"Hi","slug":"hi"}}
//
// $ curl -X DELETE  http://localhost:3000/api/v1/articles/1/
// {"articles":{"id":"1","user_id":100,"title":"Hi","slug":"hi"}}
//
// $ curl http://localhost:3000/api/v1/articles/1/
// {"error":"Page not found"}
//
// $ curl -X POST -d '{"id":"will-be-omitted","title":"awesomeness"}' http://localhost:3000/api/v1/articles/
// {"articles":{"id":"91","user_id":0,"title":"awesomeness","slug":"awesomeness"}}
//
// $ curl http://localhost:3000/api/v1/articles/91/
// {"articles":{"id":"91","user_id":0,"title":"awesomeness","slug":"awesomeness"}}
//

const DB = "boltdb.db"

var log = logger.New("main")

func hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello World"})
}

func main() {
	s := CreateNewServer(DB)
	s.MountHandlers()
	s.Router.Run(":3000")
}

type Server struct {
	Router *gin.Engine
	DB     *database.Database
}

func CreateNewServer(dbPath string) *Server {
	db, err := database.GetDatabase(dbPath, false)
	if err != nil {
		log.Fatal(err)
	}
	s := &Server{DB: db}
	s.Router = gin.Default()
	s.Router.SetTrustedProxies([]string{"0.0.0.0"})
	return s
}

func (s *Server) MountHandlers() {
	s.Router.GET("/", hello)
	s.Router.POST("/login", user.Login)

	v1 := s.Router.Group("/api/v1", user.AuthenticationCtx())
	{
		articles := v1.Group("/articles")
		{
			articles.GET("/", article.ListArticles)
			articles.POST("/", article.CreateArticle) // POST /articles

			a := articles.Group("/:articlePath") // articlePath can be articleID or articleSlug
			{
				a.Use(article.ArticleCtx())          // Load the *Article on the request context
				a.GET("/", article.ReadArticle)      // GET /articles/123
				a.PUT("/", article.UpdateArticle)    // PUT /articles/123
				a.DELETE("/", article.DeleteArticle) // DELETE /articles/123
			}
		}
	}
}
