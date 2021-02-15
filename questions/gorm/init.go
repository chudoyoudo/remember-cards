package gorm

import (
    "github.com/golobby/container"

    "github.com/chudoyoudo/remember-cards/questions"
)

func init() {
    container.Transient(func() questions.Dao {
        return &dao{}
    })
}
