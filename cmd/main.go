package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	flibusta2 "github.com/matperez/flibusta-parser/internal/flibusta"
	"github.com/matperez/flibusta-parser/internal/pool"
	storage2 "github.com/matperez/flibusta-parser/internal/storage"
	"github.com/matperez/flibusta-parser/internal/work"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var flb flibusta2.Client
var db *gorm.DB

func MakeDBConnection() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		CLI.DbUser,
		CLI.DbPassword,
		CLI.DbServer,
		CLI.DbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateFlibustaClient() flibusta2.Client {
	client, err := flibusta2.NewFlibusta()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Auth(CLI.FlibustaUser, CLI.FlibustaPassword)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func Migrate(db *gorm.DB) {
	bookProto := &storage2.Book{}
	authorProto := &storage2.Author{}
	genreProto := &storage2.Genre{}
	err := db.AutoMigrate(bookProto, authorProto, genreProto)
	if err != nil {
		log.Fatal(err)
	}
}

var CLI struct {
	DbServer         string `help:"Database server address and port" default:"localhost:3306"`
	DbName           string `help:"Database name" default:"flibusta"`
	DbUser           string `help:"Database user name" required:""`
	DbPassword       string `help:"Database user password" required:""`
	FlibustaUser     string `help:"Flibusta user name" required:""`
	FlibustaPassword string `help:"Flibusta user password" required:""`
	Parse            struct {
		WorkersCount int `help:"Workers count." short:"w" default:"4"`
		From         int `arg:"" name:"from" help:"Initial book ID." required:""`
		To           int `arg:"" name:"to" help:"Final book ID." required:""`
	} `cmd:"" help:"Run parsing."`
}

func ParseCLIContext() {
	ctx := kong.Parse(
		&CLI,
		kong.UsageOnError(),
		kong.Name("parser"),
		kong.Description("https://flibusta.is parser"),
	)
	switch ctx.Command() {
	case "parse <from> <to>":
	default:
		panic(ctx.Command())
	}
}

func main() {
	ParseCLIContext()

	db = MakeDBConnection()
	Migrate(db)

	flb = CreateFlibustaClient()

	collector := pool.StartDispatcher(CLI.Parse.WorkersCount, db, flb) // start up worker pool

	for i, job := range work.CreateJobs(CLI.Parse.From, CLI.Parse.To) {
		collector.Work <- pool.Work{BookID: job, ID: i}
	}
}
