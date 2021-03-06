package gorm

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
		return errors.Wrapf(err, "Can't create question via connection %v", *q)
	}
	return nil
}

func (dao *dao) Update(q *questions.Question, fields []string) error {
	data := q.ToMap(fields)
	result := dao.getConnection().Model(q).Updates(*data)
	err := result.Error()
	if err != nil {
		return errors.Wrapf(err, "Can't update question with id %d via connection %v", q.ID, data)
	}
	return nil

}

func (dao *dao) Delete(conds ...interface{}) error {
	result := dao.getConnection().Delete(&questions.Question{}, conds...)
	err := result.Error()
	if err != nil {
		return errors.Wrapf(err, "Can't delete question via connection by conds %v", conds)
	}
	return nil
}

func (dao *dao) Find(conds *map[string]interface{}, order *[]interface{}, limit, offset int) (list *[]questions.Question, more bool, err error) {
	ql := []questions.Question{}
	c := dao.getConnection()

	if limit > 0 {
		c = c.Limit(limit + 1)
	}

	if offset > 0 {
		c = c.Offset(offset)
	}

	if len(*order) > 0 {
		for _, o := range *order {
			c = c.Order(o)
		}
	}

	result := c.Find(&ql, *conds)

	err = result.Error()
	if err != nil {
		return &ql, false, errors.Wrapf(err, "Can't find question via connection by conds %v", conds)
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
