package questions

import "time"

type Question struct {
    ID         uint64 `gorm:"primaryKey"`
    UserId     uint64
    GroupId    uint64
    Step       uint8
    Title      string
    Body       string
    RepeatTime time.Time
    IsFailed   bool
}
