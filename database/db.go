package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error

	DB, err = sql.Open("sqlite3", "./quotes.db")

	if err != nil {
		log.Fatal("Erro ao abrir o banco:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}

	configureSQLite()
	RunMigrations()

	log.Println("database initialized")

}

func configureSQLite() {
	_, err := DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Erro ao habilitar chaves estrangeiras:", err)
	}

	//DB.Exec("PRAGMA journal_mode = WAL")        // Write-Ahead Logging para melhor concorrência
	//DB.Exec("PRAGMA synchronous = NORMAL")      // Balance entre segurança e performance
	//DB.Exec("PRAGMA cache_size = -64000")       // 64MB de cache
	//DB.Exec("PRAGMA temp_store = MEMORY")       // Tabelas temporárias em memória
	//DB.Exec("PRAGMA mmap_size = 30000000000")

	log.Println("Connected to database")
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Closed database connection")
	}
}
