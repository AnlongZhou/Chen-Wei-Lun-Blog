package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Comment struct {
	Name    string
	Comment string
}

func newComment(name, comment string) Comment {
	return Comment{
		Name:    name,
		Comment: comment,
	}
}

type Comments = []Comment

type Data struct {
	Comments Comments
}

func newData() Data {
	return Data{
		Comments: []Comment{
			newComment("Dannyyu", "I want this bad boy so much"),
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/css", "css")
	e.Static("/image", "image")

	page := newPage()
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.POST("/comment", func(c echo.Context) error {
		name := c.FormValue("name")
		comment := c.FormValue("comment")

		if comment == "" {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Errors["comment"] = "Comment cannot be empty"

			return c.Render(422, "commentForm", formData)
		}

		message := newComment(name, comment)
		page.Data.Comments = append(page.Data.Comments, message)

		c.Render(200, "commentForm", newFormData())
		return c.Render(200, "oob-comment", message)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
