package questions

type Dao interface {
	Create(q *Question) error
	Update(q *Question, fields []string) error
}
