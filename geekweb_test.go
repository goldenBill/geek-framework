package geek_framework

import (
	"github.com/goldenBill/geekweb"
	"log"
	"net/http"
	"testing"
	"time"
)

/*
(1) index
curl -i http://localhost:9999/index
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sat, 09 Jul 2022 13:56:19 GMT
Content-Length: 19

<h1>Index Page</h1>
(2) v1
$ curl -i http://localhost:9999/v1/
HTTP/1.1 200 OK
Content-Type: text/html
Date: Sat, 09 Jul 2022 13:58:06 GMT
Content-Length: 18

<h1>Hello Gee</h1>
(3)
$ curl "http://localhost:9999/v1/hello?name=geektutu"
hello geektutu, you're at /v1/hello
(4)
$ curl "http://localhost:9999/v2/hello/geektutu"
hello geektutu, you're at /hello/geektutu
(5)
$ curl "http://localhost:9999/v2/login" -X POST -d "username=geektutu&password=1234"
{"password":"1234","username":"geektutu"}
(6)
$ curl "http://localhost:9999/hello"
404 NOT FOUND: GET for URL "/hello"
*/

func onlyForV2() geekweb.HandlerFunc {
	return func(c *geekweb.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(http.StatusInternalServerError, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func TestEnigne(t *testing.T) {
	r := geekweb.New()
	r.Use(geekweb.Logger()) // global middleware
	r.GET("/index", func(c *geekweb.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *geekweb.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *geekweb.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *geekweb.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *geekweb.Context) {
			c.JSON(http.StatusOK, geekweb.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}
