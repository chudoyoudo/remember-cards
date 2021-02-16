package questions

type Dao interface {
    Create(q *Question) error
    Update(q *Question, fields []string) error
    Delete(conds ...interface{}) error
    Find(conds *map[string]interface{}, order *[]interface{}, limit, offset int) (list *[]Question, more bool, err error)
}
