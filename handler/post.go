package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/kou12345/appledore-backend/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetPosts(c echo.Context) (err error) {
	db := h.DB

	var posts []model.Post

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
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
			log.Fatal("ScanError: ", err)
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("RowsError: ", err)
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
	}

	return c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetPost(c echo.Context) (err error) {
	id := c.Param("id")

	db := h.DB

	var post model.Post
	err = db.QueryRow(`
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
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
	}

	return c.JSON(http.StatusOK, post)
}

func (h *Handler) Search(c echo.Context) (err error) {
	// query parameterを取得 検索文字列
	searchText := c.QueryParam("search")

	db := h.DB

	var posts []model.Post

	stmt, err := db.Prepare(`
		SELECT
			*
		FROM
			posts
		WHERE 
			content &@ $1 OR title &@ $1;
	`)
	if err != nil {
		log.Fatal("PrepareError: ", err)
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
	}

	rows, err := stmt.Query(searchText)
	if err != nil {
		log.Fatal("QueryError: ", err)
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
	}

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
			log.Fatal("ScanError: ", err)
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
		}

		posts = append(posts, post)
	}

	return c.JSON(http.StatusOK, posts)
}

func (h *Handler) CreatePost(c echo.Context) (err error) {
	title := c.FormValue("title")
	content := c.FormValue("content")

	// validation
	if title == "" || content == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "タイトルと内容を入力してください"}
	}
	// TODO 文字数のvalidation

	db := h.DB
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
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
	}

	var id string
	err = stmt.QueryRow(title, content, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		log.Fatal("ExecError: ", err)
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
	}

	return c.JSON(http.StatusOK, id)
}

func (h *Handler) UpdatePost(c echo.Context) (err error) {
	id := c.Param("id")
	title := c.FormValue("title")
	content := c.FormValue("content")

	// validation

	db := h.DB

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
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
	}

	// prepared statementを実行
	_, err = stmt.Exec(title, content, time.Now(), id)
	if err != nil {
		log.Fatal("ExecError: ", err)
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
	}

	return c.JSON(http.StatusOK, id)
}

func (h *Handler) DeletePost(c echo.Context) (err error) {
	id := c.Param("id")

	db := h.DB

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
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
	}

	// prepared statementを実行
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal("ExecError: ", err)
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "サーバーでエラーが発生しました"}
	}

	return c.JSON(http.StatusOK, id)
}
