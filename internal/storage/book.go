package storage

import "time"

type Book struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Title      string    `gorm:"index:,class:FULLTEXT;type:VARCHAR(255);not null;required"`
	ReadCount  uint      `gorm:"index"`
	Annotation *string   `gorm:"type:TEXT;index:,class:FULLTEXT"`
	Authors    []*Author `gorm:"many2many:book_authors;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Genres     []*Genre  `gorm:"many2many:book_genres;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Author struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string  `gorm:"index:,class:FULLTEXT;required;not null;"`
	Books     []*Book `gorm:"many2many:book_authors;"`
}

type Genre struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string  `gorm:"required;not null;"`
	Books     []*Book `gorm:"many2many:book_genres;"`
}
