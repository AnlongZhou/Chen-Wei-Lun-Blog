package main

import (
	"html/template"
	"io"
	"log"
	"mime"
	"os"

	"myWeb/assets"
	"myWeb/views"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseFS(views.FS, "*.html")),
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
	Videos VideoList
}

func newPage() Page {
	return Page{
		Data:   newData(),
		Form:   newFormData(),
		Images: newImages(),
		Videos: newVideos(),
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

type Video struct {
	URL         string
	ALT         string
	Description string
}

type Videos = []Video

type VideoList struct {
	Videos Videos
}

func newVideos() VideoList {
	return VideoList{
		Videos: []Video{
			{
				URL:         "/assets/video/video1.mp4",
				ALT:         "Video 1",
				Description: "陳緯倫在演唱會上深情演唱幾分之幾，臺下觀眾也熱情回應",
			},
		},
	}
}

func newImages() ImageList {
	return ImageList{
		Images: []Image{
			{
				URL:         "/assets/image/image1.webp",
				ALT:         "Image 1",
				Description: "陳緯倫第一次學會唱歌",
			},
			{
				URL:         "/assets/image/image2.png",
				ALT:         "Image 2",
				Description: "陳緯倫第一次參加歌唱比賽，就得到了冠軍",
			},
			{
				URL:         "/assets/image/image3.png",
				ALT:         "Image 3",
				Description: "陳緯倫參加知名節目拍的宣傳照",
			},
			{
				URL:         "/assets/image/image4.png",
				ALT:         "Image 4",
				Description: "陳緯倫第一次得到金曲獎最佳新人獎",
			},
			{
				URL:         "/assets/image/image5.png",
				ALT:         "Image 5",
				Description: "陳緯倫努力和麥克風培養感情",
			},
			{
				URL:         "/assets/image/image6.png",
				ALT:         "Image 6",
				Description: "陳緯倫參加節目飢餓遊戲時開啓野蠻模式",
			},
			{
				URL:         "/assets/image/image7.png",
				ALT:         "Image 7",
				Description: "陳緯倫在東南亞巡迴演唱會的盛況",
			},
			{
				URL:         "/assets/image/image8.jpg",
				ALT:         "Image 8",
				Description: "陳緯倫在錄音室專心錄製新專輯",
			},
			{
				URL:         "/assets/image/image9.png",
				ALT:         "Image 9",
				Description: "陳緯倫化名藏身於臺灣某間大學",
			},
			{
				URL:         "/assets/image/image10.jpg",
				ALT:         "Image 10",
				Description: "陳緯倫被發現和緋聞戀人在餐廳親熱",
			},
			{
				URL:         "/assets/image/image11.jpg",
				ALT:         "Image 11",
				Description: "陳緯倫在美國參加葛萊美獎頒獎典禮",
			},
			{
				URL:         "/assets/image/image12.png",
				ALT:         "Image 12",
				Description: "陳偉倫在臺中開演唱會跟粉絲們合照",
			},
			{
				URL:         "/assets/image/image13.png",
				ALT:         "Image 13",
				Description: "陳緯倫成為耶路撒冷帶來和平",
			},
		},
	}
}

func main() {

	mime.AddExtensionType(".webp", "image/webp")
	mime.AddExtensionType(".gif", "image/gif")

	e := echo.New()
	e.Use(middleware.Logger())
	e.StaticFS("/assets", assets.FS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
		log.Println("Port is not set, defaulting to 5000")
	}

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

		if name == "" {
			formData := newFormData()
			formData.Values["comment"] = comment
			formData.Errors["name"] = "Name cannot be empty"

			return c.Render(422, "commentForm", formData)
		}

		message := newComment(name, comment)
		page.Data.Comments = append(page.Data.Comments, message)

		c.Render(200, "commentForm", newFormData())
		return c.Render(200, "oob-comment", message)
	})

	serverAddress := ":" + port
	e.Logger.Fatal(e.Start(serverAddress))
}
