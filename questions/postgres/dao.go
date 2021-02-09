package postgres

import (
    "github.com/golobby/container"
    "gorm.io/gorm"

    "github.com/chudoyoudo/remember-cards/questions"
)

type dao struct {
    c *gorm.DB
}

func (dao *dao) Create(q *questions.Question) (*questions.Question, error) {
    c := dao.getConnection()

    result := c.Create(q)
    if result.Error != nil {
        return nil, result.Error
    }

    return q, nil
}

func (dao *dao) getConnection() *gorm.DB {
    if dao.c == nil {
        container.Make(&dao.c)
    }
    return dao.c
}
