quote-api/
├── main.go
├── go.mod
├── database/
│   ├── db.go
│   └── migrations.go       
├── models/
│   ├── author.go
│   ├── book.go
│   ├── category.go
│   ├── quote.go
│   └── associations.go    
├── dto/
│   ├── author_dto.go
│   ├── book_dto.go
│   ├── category_dto.go
│   └── quote_dto.go
├── repository/
│   ├── author_repo.go
│   ├── book_repo.go
│   ├── category_repo.go
│   └── quote_repo.go
├── service/
│   ├── author_service.go
│   ├── book_service.go
│   ├── category_service.go
│   └── quote_service.go
└── handlers/
├── author_handler.go
├── book_handler.go
├── category_handler.go
└── quote_handler.go



# Execução normal (cria tabelas se não existirem)
go run main.go

# Reset completo (remove tudo e recria)
go run main.go -reset

# Criar tabelas e inserir dados de exemplo
go run main.go -seed

# Reset + seed
go run main.go -reset -seed
```

---

## Recursos das Migrations:

✅ **Criação idempotente** - Usa `IF NOT EXISTS`, pode rodar múltiplas vezes  
✅ **Foreign Keys com CASCADE** - Deleta relacionamentos automaticamente  
✅ **Constraints** - Validações no banco (CHECK, UNIQUE)  
✅ **Índices otimizados** - Performance em buscas  
✅ **Seed data** - Dados de exemplo para testes  
✅ **Reset functionality** - Útil em desenvolvimento  
✅ **Configurações SQLite** - WAL mode, cache, etc  
✅ **Transaction safety** - Seed usa transação  
✅ **Versionamento** - Prepara para migrations futuras  

---

## Estrutura das tabelas criadas:
```
author
├── id (PK, autoincrement)
├── name (NOT NULL)
├── created_at
└── updated_at

book
├── id (PK, autoincrement)
├── title (NOT NULL)
├── isbn (UNIQUE, nullable)
├── published_year (NOT NULL, CHECK)
├── publisher (nullable)
├── pages (NOT NULL, CHECK, default 0)
├── created_at
└── updated_at

category
├── id (PK, autoincrement)
├── name (NOT NULL, UNIQUE)
├── created_at
└── updated_at

quote
├── id (PK, autoincrement)
├── book_id (FK → book.id, CASCADE)
├── text (NOT NULL, CHECK)
├── created_at
└── updated_at

book_author
├── book_id (PK, FK → book.id, CASCADE)
├── author_id (PK, FK → author.id, CASCADE)
└── order (CHECK >= 1)

book_category
├── book_id (PK, FK → book.id, CASCADE)
└── category_id (PK, FK → category.id, CASCADE)