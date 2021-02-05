package questions

import "github.com/golobby/container"

type Usecase interface {
    Add(q *Question) (*Question, error)
}

type usecase struct {
    dao Dao
}

func (u *usecase) Add(q *Question) (*Question, error) {
    dao, err := u.getDao()
    if err != nil {
        return nil, err
    }

    r, err := dao.Create(q)
    if err != nil {
        return nil, err
    }

    return r, nil
}

func (u *usecase) getDao() (Dao, error) {
    if u.dao == nil {
        container.Make(&u.dao)
    }
    return u.dao, nil
}
