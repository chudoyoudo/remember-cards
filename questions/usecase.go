package questions

import (
    "github.com/golobby/container"
    "github.com/pkg/errors"
)

type Usecase interface {
    Add(q *Question) error
}

type usecase struct {
    dao Dao
}

func (u *usecase) Add(q *Question) error {
    dao := u.getDao()

    err := dao.Create(q)
    if err != nil {
        return errors.Wrap(err, "Can't create question via dao")
    }

    return nil
}

func (u *usecase) getDao() Dao {
    if u.dao == nil {
        container.Make(&u.dao)
    }
    return u.dao
}
