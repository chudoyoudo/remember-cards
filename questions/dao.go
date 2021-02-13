package questions

type Dao interface {
    Create(q *Question) error
    Update(q *Question) error
    Delete(conds ...interface{}) error
    Find(limit, offset int, conds ...interface{}) (list *[]Question, more bool, err error)
}
