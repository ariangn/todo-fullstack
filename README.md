# todo-go

## 環境変数設定
プロジェクトルートに `.env.example`を参考に `.env` を作成する

## インストール
```
git clone https://github.com/ariangn/todo-go.git
cd todo-go
go mod tidy
go install github.com/google/wire/cmd/wire@latest
cd di && wire && cd ..
```

## Supabase テーブル準備
Supabase SQLエディタで以下を実行:
```
create table users (
  id text primary key,
  email text not null unique,
  password text not null,
  name text,
  avatar_url text,
  timezone text not null,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);
create table categories (
  id text primary key,
  name text not null,
  color text not null,
  description text,
  user_id text not null references users(id) on delete cascade,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);
create table tags (
  id text primary key,
  name text not null,
  user_id text not null references users(id) on delete cascade,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);
create table todos (
  id text primary key,
  title text not null,
  body text,
  status text not null,
  due_date timestamptz,
  completed_at timestamptz,
  user_id text not null references users(id) on delete cascade,
  category_id text references categories(id) on delete set null,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);
create table todo_tags (
  todo_id text not null references todos(id) on delete cascade,
  tag_id text not null references tags(id) on delete cascade,
  primary key(todo_id, tag_id)
);
create view todos_with_tag_ids as
select
  t.*,
  coalesce(array_agg(tt.tag_id) filter (where tt.tag_id is not null), '{}') as tag_ids
from todos t
left join todo_tags tt on tt.todo_id = t.id
group by t.id;
```

## サーバー起動
go run ./cmd
（http://localhost:8080 で待機）

## テスト実行
go test ./interface-adapter/handler

## プロジェクト構成
```
todo-app-go/
├── cmd/
│   └── main.go
├── di/
│   ├── container.go
│   └── wire_gen.go
├── infrastructure/
│   ├── auth/
│   │   └── auth_client.go
│   └── database/
│       ├── supabase_client.go
│       ├── user_repository.go
│       ├── todo_repository.go
│       ├── category_repository.go
│       ├── tag_repository.go
│       └── model/
├── application/
│   ├── user/
│   │   ├── register_use_case.go
│   │   └── login_use_case.go
│   ├── todo/
│   │   ├── create_use_case.go
│   │   ├── list_use_case.go
│   │   ├── find_by_id_use_case.go
│   │   ├── update_use_case.go
│   │   ├── toggle_status_use_case.go
│   │   ├── delete_use_case.go
│   │   └── duplicate_use_case.go
│   ├── category/
│   │   ├── create_use_case.go
│   │   ├── list_use_case.go
│   │   ├── update_use_case.go
│   │   └── delete_use_case.go
│   └── tag/
│       ├── create_use_case.go
│       ├── list_use_case.go
│       ├── update_use_case.go
│       └── delete_use_case.go
├── domain/
│   ├── entity/
│   ├── repository/
│   └── valueobject/
├── interface-adapter/
│   ├── dto/
│   ├── handler/
│   └── middleware/
├── .env
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## API 概要
### 認証不要
POST /api/users/register
POST /api/users/login

### 認証要 
/api/todos（作成/一覧/取得/更新/ステータス切替/削除/複製）
/api/categories（作成/一覧/削除）
/api/tags（作成/一覧/削除）
