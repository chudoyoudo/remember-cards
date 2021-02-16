package questions

import (
    "time"

    "github.com/golobby/container"
    "github.com/pkg/errors"
)

type Usecase interface {
    Add(q *Question) error
    Correct(q *Question) error
    Delete(conds []interface{}) error
    Find(conds *map[string]interface{}, order *[]interface{}, limit, offset int) (list *[]Question, more bool, err error)
}

type usecase struct {
    dao Dao
    now time.Time
}

func (u *usecase) Add(q *Question) error {
    originalStep := q.Step
    originalRepeatTime := q.RepeatTime
    originalIsFailed := q.IsFailed

    q.Step = 1
    q.RepeatTime = u.getRepeatTime(q.Step)
    q.IsFailed = false

    dao := u.getDao()
    err := dao.Create(q)
    if err != nil {
        q.Step = originalStep
        q.RepeatTime = originalRepeatTime
        q.IsFailed = originalIsFailed
        return errors.Wrap(err, "Can't create question via dao")
    }

    return nil
}

func (u *usecase) Correct(q *Question) error {
    dao := u.getDao()
    fields := []string{QuestionGroupId, questionTitle, questionBody}
    err := dao.Update(q, fields)
    if err != nil {
        return errors.Wrap(err, "Can't update question via dao")
    }
    return nil
}

func (u *usecase) Delete(conds []interface{}) error {
    dao := u.getDao()
    err := dao.Delete(conds...)
    if err != nil {
        return errors.Wrapf(err, "Can't delete question via dao by conds %v", conds)
    }
    return nil
}

func (u *usecase) Find(conds *map[string]interface{}, order *[]interface{}, limit, offset int) (list *[]Question, more bool, err error) {
    dao := u.getDao()
    list, more, err = dao.Find(conds, order, limit, offset)
    if err != nil {
        return list, more, errors.Wrapf(err, "Can't delete question via dao by conds %v", conds)
    }
    return list, more, err
}

func (u *usecase) getDao() Dao {
    if u.dao == nil {
        container.Make(&u.dao)
    }
    return u.dao
}

func (u *usecase) getNow() time.Time {
    var emptyTime time.Time
    if u.now == emptyTime {
        return time.Now()
    }
    return u.now
}

func (u *usecase) getRepeatTime(step uint8) time.Time {
    now := u.getNow()
    switch step {
    case 1:
        return now.Add(time.Minute * 30)
    case 2:
        return now.Add(time.Hour * 24 * 14)
    case 3:
        return now.Add(time.Hour * 24 * 60)
    default:
        return now.Add(time.Hour * 24 * 90)
    }
}
