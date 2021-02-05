package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/golobby/container"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "github.com/chudoyoudo/remember-cards/questions"
    _ "github.com/chudoyoudo/remember-cards/questions/postgres"
)

func init() {
    initPostgres()
}

func main() {

    q := &questions.Question{
        Title:      "title",
        Body:       "body",
        GroupId:    1,
        UserId:     2,
        IsFailed:   true,
        RepeatTime: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
        Step:       3,
    }

    var u questions.Usecase
    container.Make(&u)
    q, err := u.Add(q)

    if err != nil {
        log.Fatalln(err)
    }

    fmt.Println(q.ID)
}

func initPostgres() {
    container.Singleton(func() *gorm.DB {
        dsn, found := os.LookupEnv("POSTGRES_DSN")
        if !found {
            dsn = "host=localhost port=5432 user=postgres password=123 dbname=rc"
        }

        db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            log.Fatalf("Can't connect to postgresql. Error %s", err)
        }

        err = db.AutoMigrate(&questions.Question{})
        if err != nil {
            log.Fatalf("Can't migrate questions table. Error %s", err)
        }

        return db
    })
}
