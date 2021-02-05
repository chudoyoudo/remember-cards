package questions

import (
    "github.com/golobby/container"
)

func init() {
    container.Transient(func() Usecase {
        return &usecase{}
    })
}
