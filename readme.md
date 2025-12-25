# quote-api

API em Go para gerenciar autores, livros, categorias e citações.

## Estrutura do projeto

A estrutura principal do projeto é:

    `quote-api/`
    ├── `main.go`
    ├── `go.mod`
    ├── `database/`
    │   ├── `db.go`
    │   └── `migrations.go`
    ├── `models/`
    │   ├── `author.go`
    │   ├── `book.go`
    │   ├── `category.go`
    │   ├── `quote.go`
    │   └── `associations.go`
    ├── `dto/`
    │   ├── `author_dto.go`
    │   ├── `book_dto.go`
    │   ├── `category_dto.go`
    │   └── `quote_dto.go`
    ├── `repository/`
    │   ├── `author_repo.go`
    │   ├── `book_repo.go`
    │   ├── `category_repo.go`
    │   └── `quote_repo.go`
    ├── `service/`
    │   ├── `author_service.go`
    │   ├── `book_service.go`
    │   ├── `category_service.go`
    │   └── `quote_service.go`
    └── `handlers/`
        ├── `author_handler.go`
        ├── `book_handler.go`
        ├── `category_handler.go`
        └── `quote_handler.go`

## Como executar

- Execução normal (cria tabelas se não existirem)
  go run main.go

- Reset completo (remove tudo e recria)
  go run main.go -reset

- Criar tabelas e inserir dados de exemplo (seed)
  go run main.go -seed

- Reset + seed
  go run main.go -reset -seed

## Recursos das migrations

-  Criação idempotente — usa `IF NOT EXISTS`, pode rodar múltiplas vezes
-  Foreign Keys com `CASCADE` — deleta relacionamentos automaticamente
-  Constraints — validações no banco (`CHECK`, `UNIQUE`, `NOT NULL`)
-  Índices otimizados — melhor performance em buscas
-  Seed data — dados de exemplo para testes e desenvolvimento
-  Reset functionality — útil para desenvolvimento local
-  Configurações SQLite — WAL mode, cache, pragmas de performance
-  Transaction safety — seed e migrações executadas em transação quando aplicável
-  Preparado para versionamento de migrations futuras

## Estrutura das tabelas (resumo)

- `author`
    - `id` (PK, autoincrement)
    - `name` (NOT NULL)
    - `created_at`
    - `updated_at`

- `book`
    - `id` (PK, autoincrement)
    - `title` (NOT NULL)
    - `isbn` (UNIQUE, nullable)
    - `published_year` (NOT NULL, `CHECK`)
    - `publisher` (nullable)
    - `pages` (NOT NULL, `CHECK`, default 0)
    - `created_at`
    - `updated_at`

- `category`
    - `id` (PK, autoincrement)
    - `name` (NOT NULL, UNIQUE)
    - `created_at`
    - `updated_at`

- `quote`
    - `id` (PK, autoincrement)
    - `book_id` (FK → `book.id`, `CASCADE`)
    - `text` (NOT NULL, `CHECK`)
    - `created_at`
    - `updated_at`

- `book_author`
    - `book_id` (PK, FK → `book.id`, `CASCADE`)
    - `author_id` (PK, FK → `author.id`, `CASCADE`)
    - `order` (`CHECK` >= 1)

- `book_category`
    - `book_id` (PK, FK → `book.id`, `CASCADE`)
    - `category_id` (PK, FK → `category.id`, `CASCADE`)

## Notas

- Arquivo de migrations em `database/migrations.go`.
- Banco padrão usado para desenvolvimento: SQLite (configurações aplicadas no `database`).
- Seed usa transação para garantir consistência.
