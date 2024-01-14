# appledore

## 機能

- 記事を投稿できる wiki に近い？
- 記事を編集できる
- 記事を削除できる
- 記事を全文検索できる（タイトル、内容）

### 検索機能

タイトルと内容を全文検索できるようにする。

API にリクエストする時に、title を検索するのか、content を検索するのかを指定する

## DB

PostgreSQL の PGroonga を使って全文検索を行うため、PostgreSQL を使う。

### テーブル設計

- posts
  - id
  - title
  - content
  - created_at
  - updated_at

## frontend

htmx を使う。

## URL 設計

一覧ページ
投稿詳細ページ
投稿編集ページ
新規投稿ページ
検索

Go
Echo
Templ
HTMX
Tailwind
