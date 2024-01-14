package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

// type Template struct {
// 	templates *template.Template
// }

// func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
// 	return t.templates.ExecuteTemplate(w, name, data)
// }

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// t := &Template{
	// 	templates: template.Must(template.ParseGlob("public/views/*.html")),
	// }

	e := echo.New()
	e.Use(middleware.Logger())

	// e.Renderer = t

	e.GET("/hello", Hello)
	e.POST("/post", Post)

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

func Post(c echo.Context) error {
	title := c.FormValue("title")
	content := c.FormValue("content")

	// validation

	// .envからdsnを取得
	dsn := os.Getenv("DSN")
	if len(dsn) == 0 {
		log.Fatal("DSN is empty")
	}

	// DB接続
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("OpenError: ", err)
	}
	defer db.Close()

	// DB接続確認
	if err := db.Ping(); err != nil {
		log.Fatal("PingError: ", err)
	}

	// postsにINSERT
	// prepared statementを作成
	stmt, err := db.Prepare(`
		INSERT INTO posts (
			title,
			content,
			created_at,
			updated_at
		) VALUES (
			$1,
			$2,
			$3,
			$4
		) RETURNING id
	`)
	if err != nil {
		log.Fatal("PrepareError: ", err)
	}

	var id string
	// prepared statementを実行
	err = stmt.QueryRow(title, content, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		log.Fatal("ExecError: ", err)
	}

	return c.JSON(http.StatusOK, id)
}
