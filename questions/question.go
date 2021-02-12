package questions

import "time"

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
