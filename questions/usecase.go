package questions

import (
    "github.com/golobby/container"
)

type Usecase interface {
    Add(q *Question) (*Question, error)
}

type usecase struct {
    dao Dao
}

func (u *usecase) Add(q *Question) (*Question, error) {
    dao := u.getDao()

    qr, err := dao.Create(q)
    if err != nil {
        return nil, err
    }

    return qr, nil
}

func (u *usecase) getDao() Dao {
    if u.dao == nil {
        container.Make(&u.dao)
    }
    return u.dao
}
