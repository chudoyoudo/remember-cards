package questions

import "time"

const (
	questionUserId     = "userId"
	questionGroupId    = "groupId"
	questionTitle      = "title"
	questionBody       = "body"
	questionStep       = "step"
	questionRepeatTime = "repeatTime"
	questionIsFailed   = "isFailed"
)

type Question struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`
	UserId     uint64    `json:"userId"`
	GroupId    uint64    `json:"groupId"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Step       uint8     `json:"-"`
	RepeatTime time.Time `json:"repeatTime"`
	IsFailed   bool      `json:"isFailed"`
}

func (q *Question) ToMap(fields []string) *map[string]interface{} {
	if len(fields) > 0 {
		result := map[string]interface{}{}
		for _, field := range fields {
			switch field {
			case questionUserId:
				result[field] = q.UserId
			case questionGroupId:
				result[field] = q.GroupId
			case questionTitle:
				result[field] = q.Title
			case questionBody:
				result[field] = q.Body
			case questionStep:
				result[field] = q.Step
			case questionRepeatTime:
				result[field] = q.RepeatTime
			case questionIsFailed:
				result[field] = q.IsFailed
			}
		}
		return &result
	}

	return &map[string]interface{}{
		questionUserId:     q.UserId,
		questionGroupId:    q.GroupId,
		questionTitle:      q.Title,
		questionBody:       q.Body,
		questionStep:       q.Step,
		questionRepeatTime: q.RepeatTime,
		questionIsFailed:   q.IsFailed,
	}
}
