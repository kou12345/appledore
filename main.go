package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kou12345/appledore-backend/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	// .envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// DB接続
	db := ConnectDB()
	defer db.Close()

	// handlerにDBインスタンスを渡す
	h := &handler.Handler{DB: db}

	// ルーティング
	e.GET("/", h.GetPosts)
	e.GET("/search", h.Search)
	e.GET("/post/:id", h.GetPost)
	e.POST("/post", h.CreatePost)
	e.PUT("/post/:id", h.UpdatePost)
	e.DELETE("/post/:id", h.DeletePost)

	// サーバー起動
	e.Logger.Fatal(e.Start(":1323"))
}

// DBインスタンスを作成
func ConnectDB() *sql.DB {
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
