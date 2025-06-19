# todo-fullstack
## 概要
このリポジトリは、React + Vite + TypeScript を使ったフロントエンド、Go (Chi フレームワーク) を使ったバックエンド、Supabase（PostgreSQL）をデータベースに利用したフルスタック Todo アプリのソースコードです。

- フロントエンドは Vercel にデプロイ
- バックエンドは DigitalOcean App Platform にデプロイ
- データベースは Supabase を利用

ユーザー登録・ログイン、タスクの作成・編集・削除・ドラッグ＆ドロップによるステータス更新、カテゴリ・タグ機能などを備えています。

## ローカル開発環境構築
### 1. リポジトリをクローン
```bash
git clone https://github.com/yourusername/todo-fullstack.git
cd todo-fullstack
```
### 2. 環境変数を設定
プロジェクトルートに `.env` ファイルを作成して、以下を記述：
```
# backend
DATABASE_URL="postgresql://postgres:example@db.ukqfqjrnjenfvtayvrvp.example.co:1111/postgres"
SUPABASE_URL="https://example.supabase.co"
SUPABASE_KEY="your_supabase_key"
JWT_SECRET="your_jwt_secret"
CLIENT_ORIGIN="http://localhost:5173"

# frontend
VITE_API_URL=http://localhost:8080/api
```
- フロントエンドは `VITE_` プレフィックスを使います。

### 3. Supabase セットアップ
1. Supabase で新規プロジェクトを作成

2.  Supabase テーブル準備
Supabase SQLエディタで以下を実行:
```
create extension if not exists "uuid-ossp";

create table users (
  id uuid primary key default uuid_generate_v4(),
  email text not null unique,
  password text not null,
  name text,
  avatar_url text,
  timezone text not null,
  created_at timestamp with time zone default now(),
  updated_at timestamp with time zone default now()
);
```
```
create table if not exists public.categories (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  color text not null,
  description text,
  user_id uuid not null,
  created_at timestamp with time zone default now(),
  updated_at timestamp with time zone default now()
);

-- optional: index for quick lookup by user
create index if not exists idx_categories_user_id on public.categories (user_id);
```
```
create table if not exists public.tags (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  user_id uuid not null,
  created_at timestamp with time zone default now(),
  updated_at timestamp with time zone default now()
);

create index if not exists idx_tags_user_id on public.tags (user_id);
```
```
create table if not exists public.todos (
  id uuid primary key default gen_random_uuid(),
  title text not null,
  body text,
  status text not null,
  due_date timestamp with time zone,
  completed_at timestamp with time zone,
  user_id uuid not null,
  category_id uuid references public.categories(id) on delete set null,
  tag_ids uuid[] not null default '{}',
  created_at timestamp with time zone default now(),
  updated_at timestamp with time zone default now()
);
```
```
create index if not exists idx_todos_user_id on public.todos (user_id);
create index if not exists idx_todos_category_id on public.todos (category_id);

CREATE OR REPLACE VIEW todos_with_tag_ids AS
SELECT
  todos.*,
  COALESCE(ARRAY_AGG(todo_tags.tag_id), '{}') AS tag_ids
FROM todos
LEFT JOIN todo_tags ON todos.id = todo_tags.todo_id
GROUP BY todos.id;
```
```
create table todo_tags (
  todo_id uuid references todos(id) on delete cascade,
  tag_id uuid references tags(id) on delete cascade,
  primary key (todo_id, tag_id)
);
```

3. API キーと URL を `.env` に設定
4. バックエンドを起動
```
cd backend
go mod tidy
go run ./cmd
```
- デフォルトで `:8080` で起動します。
- `DATABASE_URL` と `JWT_SECRET` が正しく設定されていることを確認してください。
5. フロントエンドを起動
```
cd frontend
npm install
npm run dev
```
- デフォルトで `http://localhost:5173` が開きます。
### プロジェクト構成
```
todo-fullstack/
├── backend/
│   ├── application/
│   │   ├── category/
│   │   ├── tag/
│   │   ├── todo/
│   │   └── user/
│   ├── cmd/
│   │   └── main.go
│   ├── di/
│   │   └── container.go
│   ├── domain/
│   │   ├── entity/
│   │   ├── repository/
│   │   └── valueobject/
│   ├── infrastructure/
│   │   ├── auth/
│   │   └── database/
│   ├── interface-adapter/
│   │   ├── dto/
│   │   ├── handler/
│   │   └── middleware/
│   ├── .env
│   ├── .env.example
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── node_modules/
│   ├── public/
│   │   ├── logo.png
│   ├── src/
│   │   ├── assets/
│   │   ├── components/
│   │   ├── lib/
│   │   ├── pages/
│   │   ├── services/
│   │   ├── utils/
│   │   ├── App.css
│   │   ├── App.tsx
│   │   ├── index.css
│   │   ├── main.tsx
│   │   ├── types.ts
│   │   └── vite-env.d.ts
│   ├── .env
│   ├── .gitignore
│   ├── components.json
│   ├── eslint.config.js
│   ├── index.html
│   ├── package-lock.json
│   ├── package.json
│   ├── tailwind.config.js
│   ├── tsconfig.app.json
│   ├── tsconfig.json
│   ├── tsconfig.node.json
│   └── vite.config.ts
├── .gitignore
└── README.md
```
