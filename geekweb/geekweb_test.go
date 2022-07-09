package geekweb

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"testing"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func TestEnigne(t *testing.T) {
	r := Default()
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")
	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *Context) {
		c.HTML(http.StatusOK, "arr.tmpl", H{
			"title":  "geektutu",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/date", func(c *Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", H{
			"title": "geektutu",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})
	r.GET("/panic", func(c *Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	v1 := r.Group("/v1")
	v1.Use(func(c *Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(http.StatusInternalServerError, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	})
	{
		v1.GET("/hello/:name", func(c *Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v1.POST("/login", func(c *Context) {
			c.JSON(http.StatusOK, H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}
