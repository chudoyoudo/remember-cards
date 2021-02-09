package postgres

import (
    "github.com/golobby/container"
    "github.com/pkg/errors"
    "gorm.io/gorm"

    "github.com/chudoyoudo/remember-cards/questions"
)

type dao struct {
    c *gorm.DB
}

func (dao *dao) Create(q *questions.Question) error {
    c := dao.getConnection()

    result := c.Create(q)
    if result.Error != nil {
        return errors.Wrap(result.Error, "Can't create question via gorm")
    }

    return nil
}

func (dao *dao) getConnection() *gorm.DB {
    if dao.c == nil {
        container.Make(&dao.c)
    }
    return dao.c
}
