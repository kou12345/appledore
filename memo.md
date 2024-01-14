# appledore

## 機能

- 記事を投稿できる wiki に近い？
- 記事を編集できる
- 記事を削除できる
- 記事を全文検索できる（タイトル、内容）

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
