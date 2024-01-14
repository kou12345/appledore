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

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.Use(middleware.Logger())

	// TODO UTCをJSTに変更

	e.GET("/", GetPosts)
	e.GET("/post/:id", GetPost)
	e.POST("/post", NewPost)
	e.PUT("/post/:id", UpdatePost)
	e.DELETE("/post/:id", DeletePost)

	e.Logger.Fatal(e.Start(":1323"))
}

type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetPosts(c echo.Context) error {
	db := ConnectionDB()
	defer db.Close()

	var posts []Post

	// postsをSELECT
	rows, err := db.Query(`
		SELECT 
			id,
			title,
			content,
			created_at,
			updated_at
		FROM
			posts;
	`)
	if err != nil {
		log.Fatal("QueryError: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
			log.Fatal("ScanError: ", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("RowsError: ", err)
	}

	return c.JSON(http.StatusOK, posts)
}

func GetPost(c echo.Context) error {

	// TODO Titleで取得できた方が良いかも
	id := c.Param("id")

	// idを元にpostsをSELECT

	db := ConnectionDB()
	defer db.Close()

	var post Post
	err := db.QueryRow(`
		SELECT 
			id,
			title,
			content,
			created_at,
			updated_at
		FROM
			posts
		WHERE
			id = $1;
	`, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		log.Fatal("QueryRowError: ", err)
	}

	return c.JSON(http.StatusOK, post)
}

func NewPost(c echo.Context) error {
	title := c.FormValue("title")
	content := c.FormValue("content")

	// validation

	db := ConnectionDB()
	defer db.Close()

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

func UpdatePost(c echo.Context) error {
	id := c.Param("id")
	title := c.FormValue("title")
	content := c.FormValue("content")

	// validation

	db := ConnectionDB()
	defer db.Close()

	// postsをUPDATE

	// prepared statementを作成
	stmt, err := db.Prepare(`
		UPDATE 
			posts
		SET
			title = $1,
			content = $2,
			updated_at = $3
		WHERE
			id = $4;
	`)
	if err != nil {
		log.Fatal("PrepareError: ", err)
	}

	// prepared statementを実行
	_, err = stmt.Exec(title, content, time.Now(), id)
	if err != nil {
		log.Fatal("ExecError: ", err)
	}

	return c.JSON(http.StatusOK, id)
}

func DeletePost(c echo.Context) error {
	id := c.Param("id")

	db := ConnectionDB()
	defer db.Close()

	// postsをDELETE
	// prepared statementを作成
	stmt, err := db.Prepare(`
		DELETE FROM
			posts
		WHERE
			id = $1;
	`)
	if err != nil {
		log.Fatal("PrepareError: ", err)
	}

	// prepared statementを実行
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal("ExecError: ", err)
	}

	return c.JSON(http.StatusOK, id)
}

// DBインスタンスを作成
func ConnectionDB() *sql.DB {
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

	// DB接続確認
	if err := db.Ping(); err != nil {
		log.Fatal("PingError: ", err)
	}

	return db
}
