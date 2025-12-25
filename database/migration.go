package database

import "log"

func RunMigrations() {
	log.Println("Running migrations...")

	createAuthorTable()
	createBookTable()
	createBookAuthorTable()
	createBookCategoryTable()
	createQuoteTable()
	createBookCategoryTable()

	log.Println("Migrations done.")
}

func createAuthorTable() {
	query := `
    CREATE TABLE IF NOT EXISTS author (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar tabela author:", err)
	}
	log.Println(" Tabela author criada/verificada")
}

func createBookTable() {
	query := `
    CREATE TABLE IF NOT EXISTS book (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        isbn TEXT UNIQUE,
        published_year INTEGER NOT NULL,
        publisher TEXT,
        pages INTEGER NOT NULL DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        
        CHECK (published_year >= 0 AND published_year <= 9999),
        CHECK (pages >= 0)
    );
    `

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar tabela book:", err)
	}
	log.Println("Tabela book criada/verificada")

}

func createCategoryTable() {
	query := `
    CREATE TABLE IF NOT EXISTS category (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar tabela category:", err)
	}
	log.Println("Tabela category criada/verificada")
}

func createQuoteTable() {
	query := `
    CREATE TABLE IF NOT EXISTS quote (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        book_id INTEGER NOT NULL,
        text TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        
        FOREIGN KEY (book_id) REFERENCES book(id) ON DELETE CASCADE,
        
        CHECK (LENGTH(text) >= 1)
    );
    `

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar tabela quote:", err)
	}
	log.Println("Tabela quote criada/verificada")
}

func createBookAuthorTable() {
	query := `
    CREATE TABLE IF NOT EXISTS book_author (
        book_id INTEGER NOT NULL,
        author_id INTEGER NOT NULL,
        "order" INTEGER NOT NULL DEFAULT 1,
        
        PRIMARY KEY (book_id, author_id),
        FOREIGN KEY (book_id) REFERENCES book(id) ON DELETE CASCADE,
        FOREIGN KEY (author_id) REFERENCES author(id) ON DELETE CASCADE,
        
        CHECK ("order" >= 1)
    );
    `

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("❌ Erro ao criar tabela book_author:", err)
	}
	log.Println("  ✓ Tabela book_author criada/verificada")
}

func createBookCategoryTable() {
	query := `
    CREATE TABLE IF NOT EXISTS book_category (
        book_id INTEGER NOT NULL,
        category_id INTEGER NOT NULL,
        
        PRIMARY KEY (book_id, category_id),
        FOREIGN KEY (book_id) REFERENCES book(id) ON DELETE CASCADE,
        FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
    );
    `

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar tabela book_category:", err)
	}
	log.Println(" Tabela book_category criada/verificada")
}

func DropAllTables() {
	log.Println("Removendo todas as tabelas...")

	tables := []string{
		"DROP TABLE IF EXISTS book_category",
		"DROP TABLE IF EXISTS book_author",
		"DROP TABLE IF EXISTS quote",
		"DROP TABLE IF EXISTS category",
		"DROP TABLE IF EXISTS book",
		"DROP TABLE IF EXISTS author",
	}

	for _, dropQuery := range tables {
		_, err := DB.Exec(dropQuery)
		if err != nil {
			log.Printf("Aviso ao remover tabela: %v", err)
		}
	}

	log.Println("Todas as tabelas removidas!")
}
