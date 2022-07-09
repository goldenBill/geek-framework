package geek_framework

import (
	"github.com/goldenBill/geekweb"
	"net/http"
	"testing"
)

func TestEnigne(t *testing.T) {
	r := geekweb.New()

	r.GET("/", func(c *geekweb.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello/:name", func(c *geekweb.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.POST("/login", func(c *geekweb.Context) {
		c.JSON(http.StatusOK, geekweb.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
