package questions

import "time"

const (
    QuestionId         = "id"
    QuestionUserId     = "userId"
    QuestionGroupId    = "groupId"
    QuestionTitle      = "title"
    QuestionBody       = "body"
    QuestionStep       = "step"
    QuestionRepeatTime = "repeatTime"
    QuestionIsFailed   = "isFailed"
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
            case QuestionUserId:
                result[field] = q.UserId
            case QuestionGroupId:
                result[field] = q.GroupId
            case QuestionTitle:
                result[field] = q.Title
            case QuestionBody:
                result[field] = q.Body
            case QuestionStep:
                result[field] = q.Step
            case QuestionRepeatTime:
                result[field] = q.RepeatTime
            case QuestionIsFailed:
                result[field] = q.IsFailed
            }
        }
        return &result
    }

    return &map[string]interface{}{
        QuestionUserId:     q.UserId,
        QuestionGroupId:    q.GroupId,
        QuestionTitle:      q.Title,
        QuestionBody:       q.Body,
        QuestionStep:       q.Step,
        QuestionRepeatTime: q.RepeatTime,
        QuestionIsFailed:   q.IsFailed,
    }
}
