package postgres

import (
    "testing"

    gorm "github.com/chudoyoudo/gorm-interface"
    "github.com/pkg/errors"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "github.com/chudoyoudo/remember-cards/questions"
)

func Test_dao_create_question_when_connection_ok(t *testing.T) {
    qIn := &questions.Question{}
    qExpected := &questions.Question{ID: 1}

    c := &gorm.ConnectionMock{}
    c.On("Create", qIn).Return(&gorm.ConnectionMock{}).Run(func(args mock.Arguments) {
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

    c := &gorm.ConnectionMock{}
    c.On("Create", qIn).Return(&gorm.ConnectionMock{Err: connectionErr})
    dao := &dao{c: c}

    errResult := dao.Create(qIn)

    c.AssertNumberOfCalls(t, "Create", 1)
    require.ErrorIs(t, errResult, connectionErr, "Результирующая ошибка должна содержать информацию из connection")
}

func Test_dao_update_question_when_connection_ok(t *testing.T) {
    qIn := &questions.Question{
        ID:    1,
        Title: "Test 1",
    }
    qExpected := &questions.Question{
        ID:    1,
        Title: "Test 2",
    }

    c := &gorm.ConnectionMock{}
    c.On("Save", qIn).Return(&gorm.ConnectionMock{}).Run(func(args mock.Arguments) {
        qIn := args.Get(0).(*questions.Question)
        qIn.Title = "Test 2"
    })
    dao := &dao{c: c}

    errResult := dao.Update(qIn)

    c.AssertNumberOfCalls(t, "Save", 1)
    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
    assert.Equal(t, qExpected, qIn, "Возвращаемый question не содержит изменения из connection")
}

func Test_dao_update_question_when_connection_failed(t *testing.T) {
    qIn := &questions.Question{}
    connectionErr := errors.New("Connection mock error")

    c := &gorm.ConnectionMock{}
    c.On("Save", qIn).Return(&gorm.ConnectionMock{Err: connectionErr})
    dao := &dao{c: c}

    errResult := dao.Update(qIn)

    c.AssertNumberOfCalls(t, "Save", 1)
    require.ErrorIs(t, errResult, connectionErr, "Результирующая ошибка должна содержать информацию из connection")
}

func Test_dao_delete_question_when_connection_ok(t *testing.T) {
    q := &questions.Question{}
    conds := []interface{}{uint64(1)}
    c := &gorm.ConnectionMock{}
    c.On("Delete", q, conds).Return(&gorm.ConnectionMock{})
    dao := &dao{c: c}

    errResult := dao.Delete(conds...)

    c.AssertNumberOfCalls(t, "Delete", 1)
    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
}

func Test_dao_delete_question_when_connection_failed(t *testing.T) {
    q := &questions.Question{}
    conds := []interface{}{uint64(1)}
    connectionErr := errors.New("Connection mock error")

    c := &gorm.ConnectionMock{}
    c.On("Delete", q, conds).Return(&gorm.ConnectionMock{Err: connectionErr})
    dao := &dao{c: c}

    errResult := dao.Delete(conds...)

    c.AssertNumberOfCalls(t, "Delete", 1)
    require.ErrorIs(t, errResult, connectionErr, "Результирующая ошибка должна содержать информацию из connection")
}

func Test_dao_find_question_without_limit_offset_when_connection_ok(t *testing.T) {
    conds := []interface{}{}
    qlExpected := &[]questions.Question{{ID: 1}}

    c := &gorm.ConnectionMock{}
    c.On("Find", &[]questions.Question{}, conds).Return(&gorm.ConnectionMock{}).Run(func(args mock.Arguments) {
        qlOut := args.Get(0).(*[]questions.Question)
        *qlOut = append(*qlOut, questions.Question{ID: 1})
    })
    dao := &dao{c: c}

    qlResult, more, errResult := dao.Find(0, 0, conds...)

    c.AssertNumberOfCalls(t, "Find", 1)
    c.AssertNumberOfCalls(t, "Limit", 0)
    c.AssertNumberOfCalls(t, "Offset", 0)
    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
    assert.Equal(t, qlExpected, qlResult, "Возвращаемый список question не содержит данные из connection")
    assert.Equal(t, false, more, "Должны возвращаться все записи т.к. лимит не установлен")
}

func Test_dao_find_question_with_limit_without_offset_when_connection_ok(t *testing.T) {
    conds := []interface{}{}
    qlExpected := &[]questions.Question{{ID: 1}}

    cFind := &gorm.ConnectionMock{}
    cLimit := &gorm.ConnectionMock{}
    cFind.On("Find", &[]questions.Question{}, conds).Return(&gorm.ConnectionMock{}).Run(func(args mock.Arguments) {
        qlOut := args.Get(0).(*[]questions.Question)
        *qlOut = append(*qlOut, questions.Question{ID: 1})
        *qlOut = append(*qlOut, questions.Question{ID: 2})
    })
    cLimit.On("Limit", 2).Return(cFind)

    dao := &dao{c: cLimit}

    qlResult, more, errResult := dao.Find(1, 0, conds...)

    cLimit.AssertNumberOfCalls(t, "Limit", 1)
    cFind.AssertNumberOfCalls(t, "Offset", 0)
    cFind.AssertNumberOfCalls(t, "Find", 1)
    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
    assert.Equal(t, qlExpected, qlResult, "Возвращаемый список question не содержит данные из connection")
    assert.Equal(t, true, more, "Должны возвращаться все записи т.к. лимит не установлен")
}

func Test_dao_find_question_when_connection_failed(t *testing.T) {
    conds := []interface{}{}
    qlOut := &[]questions.Question{}
    connectionErr := errors.New("Connection mock error")

    c := &gorm.ConnectionMock{}
    c.On("Find", qlOut, conds).Return(&gorm.ConnectionMock{Err: connectionErr})
    dao := &dao{c: c}

    _, _, errResult := dao.Find(0, 0, conds...)

    c.AssertNumberOfCalls(t, "Find", 1)
    require.ErrorIs(t, errResult, connectionErr, "Результирующая ошибка должна содержать информацию из connection")
}
