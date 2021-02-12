package main

import (
    "log"
    "os"

    gorm_interface "github.com/chudoyoudo/gorm-interface"
    "github.com/golobby/container"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "github.com/gin-gonic/gin"

    "github.com/chudoyoudo/remember-cards/questions"
    question_gin "github.com/chudoyoudo/remember-cards/questions/gin"
    _ "github.com/chudoyoudo/remember-cards/questions/postgres"
)

func init() {
    initPostgres()
}

func main() {
    RunHttpServer()
}

// Метод запускает http сервер
func RunHttpServer() {
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(gin.Logger())
    question_gin.RegisterHandlers(r)
    if err := r.Run(":8080"); err != nil {
        log.Fatalln(err)
    }
}

func initPostgres() {
    container.Singleton(func() gorm_interface.Connection {
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
