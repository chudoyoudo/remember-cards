package postgres

import (
    gorm "github.com/chudoyoudo/gorm-interface"
    "github.com/pkg/errors"

    "github.com/chudoyoudo/remember-cards/questions"
)

type dao struct {
    c gorm.Connection
}

func (dao *dao) Create(q *questions.Question) error {
    result := dao.getConnection().Create(q)
    err := result.Error()
    if err != nil {
        return errors.Wrapf(err, "Can't create question via gorm %v", *q)
    }
    return nil
}

func (dao *dao) Update(q *questions.Question) error {
    result := dao.getConnection().Save(q)
    err := result.Error()
    if err != nil {
        return errors.Wrapf(err, "Can't update question via gorm %v", *q)
    }
    return nil
}

func (dao *dao) Delete(conds ...interface{}) error {
    result := dao.getConnection().Delete(&questions.Question{}, conds...)
    err := result.Error()
    if err != nil {
        return errors.Wrapf(err, "Can't delete question via gorm by conds %v", conds)
    }
    return nil
}

func (dao *dao) Find(limit, offset int, conds ...interface{}) (list *[]questions.Question, more bool, err error) {
    ql := []questions.Question{}
    c := dao.getConnection()
    var result gorm.Connection

    if limit > 0 {
        limit++
    }

    if limit == 0 && offset == 0 {
        result = c.Find(&ql, conds...)
    } else if limit == 0 && offset != 0 {
        result = c.Offset(offset).Find(&ql, conds...)
    } else if limit != 0 && offset == 0 {
        result = c.Limit(limit).Find(&ql, conds...)
    } else {
        result = c.Limit(limit).Offset(offset).Find(&ql, conds...)
    }

    if limit > 0 {
        limit--
    }

    err = result.Error()
    if err != nil {
        return &ql, false, errors.Wrapf(err, "Can't find question via gorm by conds %v", conds)
    }

    more = false
    if limit > 0 && len(ql) >= limit+1 {
        ql = ql[:limit]
        more = true
    }

    return &ql, more, nil
}

func (dao *dao) getConnection() gorm.Connection {
    if dao.c == nil {
        dao.c = gorm.NewConnection()
    }
    return dao.c
}
