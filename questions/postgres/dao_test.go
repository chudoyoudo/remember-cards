package postgres

import (
    "testing"

    gorm_interface "github.com/chudoyoudo/gorm-interface"
    "github.com/pkg/errors"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "gorm.io/gorm"

    "github.com/chudoyoudo/remember-cards/questions"
)

func Test_dao_create_question_when_connection_ok(t *testing.T) {
    qIn := &questions.Question{}
    qExpected := &questions.Question{ID: 1}

    c := &gorm_interface.ConnectionMock{}
    c.On("Create", qIn).Return(&gorm.DB{}).Run(func(args mock.Arguments) {
        qIn := args.Get(0).(*questions.Question)
        qIn.ID = 1
    })
    dao := &dao{c: c}

    errResult := dao.Create(qIn)

    c.AssertNumberOfCalls(t, "Create", 1)
    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
    assert.Equal(t, qExpected.ID, qIn.ID, "ID не заполнен")
}

func Test_dao_create_question_when_connection_failed(t *testing.T) {
    qIn := &questions.Question{}
    connectionErr := errors.New("Connection mock error")

    c := &gorm_interface.ConnectionMock{}
    c.On("Create", qIn).Return(&gorm.DB{Error: connectionErr})
    dao := &dao{c: c}

    errResult := dao.Create(qIn)

    c.AssertNumberOfCalls(t, "Create", 1)
    require.ErrorIs(t, errResult, connectionErr, "Результирующая ошибка должна содержать информацию из connection")
}
