package work

import (
	flibusta2 "github.com/matperez/flibusta-parser/internal/flibusta"
	storage2 "github.com/matperez/flibusta-parser/internal/storage"
	"gorm.io/gorm"
	"log"
)

func CreateJobs(from, to int) []int {
	var jobs []int

	for i := from; i < to; i++ {
		jobs = append(jobs, i)
	}
	return jobs
}

func MapBookToStore(b *flibusta2.Book) *storage2.Book {
	model := &storage2.Book{
		ID:         uint(b.ID),
		Title:      b.Title,
		Annotation: nil,
		Authors:    []*storage2.Author{},
		Genres:     []*storage2.Genre{},
	}
	if b.Annotation != "" && b.Annotation != "отсутствует" {
		model.Annotation = &b.Annotation
	}
	model.ReadCount = uint(b.ReadCount)
	for _, a := range b.Authors {
		model.Authors = append(model.Authors, &storage2.Author{
			ID:   uint(a.ID),
			Name: a.Name,
		})
	}
	for _, g := range b.Genres {
		model.Genres = append(model.Genres, &storage2.Genre{
			ID:    uint(g.ID),
			Title: g.Name,
		})
	}
	return model
}

func DoWork(db *gorm.DB, flb flibusta2.Client, bookId int, workerId int) {
	log.Printf("worker [%d] - created processing book [%d]\n", workerId, bookId)
	book, err := flb.GetBook(bookId)
	if err != nil {
		log.Printf("worker [%d] failed to fetch the book [%d]: %s", workerId, bookId, err.Error())
		return
	}
	model := MapBookToStore(book)
	db.Create(&model)
	db.Save(&model)
	log.Printf("worker [%d] stored the book [%d]", workerId, bookId)
}
