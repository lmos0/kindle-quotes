package models

type BookAuthor struct {
	BookId   int64
	AuthorId int
	Order    int
}

type BookCategory struct {
	BookId     int64
	CategoryId int
}
