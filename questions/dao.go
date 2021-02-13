package questions

type Dao interface {
    Create(q *Question) error
}
