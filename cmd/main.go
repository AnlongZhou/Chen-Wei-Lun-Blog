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
		Comments: []Comment{},
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
	Data   Data
	Form   FormData
	Images ImageList
}

func newPage() Page {
	return Page{
		Data:   newData(),
		Form:   newFormData(),
		Images: newImages(),
	}
}

type Image struct {
	URL         string
	ALT         string
	Description string
}

type Images = []Image

type ImageList struct {
	Images Images
}

func newImages() ImageList {
	return ImageList{
		Images: []Image{
			{
				URL:         "/image/image1.webp",
				ALT:         "Image 1",
				Description: "陳緯倫第一次學會唱歌",
			},
			{
				URL:         "/image/image2.png",
				ALT:         "Image 2",
				Description: "陳緯倫第一次參加歌唱比賽，就得到了冠軍",
			},
			{
				URL:         "/image/image3.png",
				ALT:         "Image 3",
				Description: "陳緯倫參加知名節目拍的宣傳照",
			},
			{
				URL:         "/image/image4.png",
				ALT:         "Image 4",
				Description: "陳緯倫第一次得到金曲獎最佳新人獎",
			},
			{
				URL:         "/image/image5.png",
				ALT:         "Image 5",
				Description: "陳緯倫努力和麥克風培養感情",
			},
			{
				URL:         "/image/image6.png",
				ALT:         "Image 6",
				Description: "陳緯倫參加節目飢餓遊戲時開啓野蠻模式",
			},
			{
				URL:         "/image/image7.png",
				ALT:         "Image 7",
				Description: "陳緯倫在東南亞巡迴演唱會的盛況",
			},
		},
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

	e.Logger.Fatal(e.Start(":5000"))
}
