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

// ----------------
// ---- Create ----
// ----------------

func Test_dao_create_connection_calls_is_correct(t *testing.T) {
    qIn := &questions.Question{}

    c := &gorm.ConnectionMock{}
    c.On("Create", qIn).Return(c)
    dao := &dao{c: c}

    _ = dao.Create(qIn)

    createCalls := 1
    if !c.AssertNumberOfCalls(t, "Create", createCalls) {
        t.Errorf("Метод Create у connection должен вызваться %d раз", createCalls)
        t.Fail()
    }
}

func Test_dao_create_when_connection_work_success_result_error_is_empty(t *testing.T) {
    qIn := &questions.Question{}

    c := &gorm.ConnectionMock{}
    c.On("Create", qIn).Return(&gorm.ConnectionMock{})
    dao := &dao{c: c}

    errResult := dao.Create(qIn)

    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
}

func Test_dao_create_when_connection_work_success_we_have_correct_result_question(t *testing.T) {
    qIn := &questions.Question{}
    qExpected := &questions.Question{ID: 1}

    c := &gorm.ConnectionMock{}
    c.On("Create", qIn).Return(&gorm.ConnectionMock{}).Run(func(args mock.Arguments) {
        qIn := args.Get(0).(*questions.Question)
        qIn.ID = 1
    })
    dao := &dao{c: c}

    _ = dao.Create(qIn)

    assert.Equal(t, qExpected.ID, qIn.ID, "Результируещий объект question должен содержать данные, пришедшие из connection")
}

func Test_dao_create_when_connection_work_wrong_result_error_not_empty_and_have_info_from_connection(t *testing.T) {
    qIn := &questions.Question{}
    connectionErr := errors.New("Connection mock error")

    c := &gorm.ConnectionMock{}
    c.On("Create", qIn).Return(&gorm.ConnectionMock{Err: connectionErr})
    dao := &dao{c: c}

    errResult := dao.Create(qIn)

    require.NotNil(t, errResult, "Возвращаемая ошибка не должна быть пустой")
    assert.ErrorIs(t, errResult, connectionErr, "Возвращаемая ошибка должна содержать информацию из connection")
}

// ----------------
// ---- Update ----
// ----------------

func Test_dao_update_connection_calls_is_correct(t *testing.T) {
    qIn := &questions.Question{}

    c := &gorm.ConnectionMock{}
    c.On("Save", qIn).Return(c)
    dao := &dao{c: c}

    _ = dao.Update(qIn)

    saveCalls := 1
    if !c.AssertNumberOfCalls(t, "Save", saveCalls) {
        t.Errorf("Метод Save у connection должен вызваться %d раз", saveCalls)
        t.Fail()
    }
}

func Test_dao_update_when_connection_work_success_result_error_is_empty(t *testing.T) {
    qIn := &questions.Question{}

    c := &gorm.ConnectionMock{}
    c.On("Save", qIn).Return(&gorm.ConnectionMock{})
    dao := &dao{c: c}

    errResult := dao.Update(qIn)

    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
}

func Test_dao_update_when_connection_work_success_we_have_correct_result_question(t *testing.T) {
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

    _ = dao.Update(qIn)

    assert.Equal(t, qExpected, qIn, "Результируещий объект question не содержит изменения из connection")
}

func Test_dao_update_when_connection_work_wrong_result_error_not_empty_and_have_info_from_connection(t *testing.T) {
    qIn := &questions.Question{}
    connectionErr := errors.New("Connection mock error")

    c := &gorm.ConnectionMock{}
    c.On("Save", qIn).Return(&gorm.ConnectionMock{Err: connectionErr})
    dao := &dao{c: c}

    errResult := dao.Update(qIn)

    require.NotNil(t, errResult, "Возвращаемая ошибка не должна быть пустой")
    assert.ErrorIs(t, errResult, connectionErr, "Возвращаемая ошибка должна содержать информацию из connection")
}

// ----------------
// ---- Delete ----
// ----------------

func Test_dao_delete_connection_calls_is_correct(t *testing.T) {
    q := &questions.Question{}
    conds := []interface{}{uint64(1)}

    c := &gorm.ConnectionMock{}
    c.On("Delete", q, conds).Return(c)
    dao := &dao{c: c}

    _ = dao.Delete(conds...)

    deleteCalls := 1
    if !c.AssertNumberOfCalls(t, "Delete", deleteCalls) {
        t.Errorf("Метод Delete у connection должен вызваться %d раз", deleteCalls)
        t.Fail()
    }
}

func Test_dao_delete_when_connection_work_success_result_error_is_empty(t *testing.T) {
    q := &questions.Question{}
    conds := []interface{}{uint64(1)}

    c := &gorm.ConnectionMock{}
    c.On("Delete", q, conds).Return(&gorm.ConnectionMock{})
    dao := &dao{c: c}

    errResult := dao.Delete(conds...)

    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
}

func Test_dao_delete_when_connection_work_wrong_result_error_not_empty_and_have_info_from_connection(t *testing.T) {
    q := &questions.Question{}
    conds := []interface{}{uint64(1)}
    connectionErr := errors.New("Connection mock error")

    c := &gorm.ConnectionMock{}
    c.On("Delete", q, conds).Return(&gorm.ConnectionMock{Err: connectionErr})
    dao := &dao{c: c}

    errResult := dao.Delete(conds...)

    require.NotNil(t, errResult, "Возвращаемая ошибка не должна быть пустой")
    assert.ErrorIs(t, errResult, connectionErr, "Возвращаемая ошибка должна содержать информацию из connection")
}

// --------------
// ---- Find ----
// --------------

func Test_dao_find_when_dont_set_limit_offset_connection_calls_is_correct(t *testing.T) {
    var conds []interface{}

    c := &gorm.ConnectionMock{}
    c.On("Find", &[]questions.Question{}, conds).Return(c)
    dao := &dao{c: c}

    _, _, _ = dao.Find(0, 0, conds...)

    c.AssertNumberOfCalls(t, "Find", 1)
    c.AssertNumberOfCalls(t, "Limit", 0)
    c.AssertNumberOfCalls(t, "Offset", 0)
}

func Test_dao_find_when_dont_set_limit_and_set_offset_connection_calls_is_correct(t *testing.T) {
    var conds []interface{}
    offset := 1

    c := &gorm.ConnectionMock{}
    c.On("Find", &[]questions.Question{}, conds).Return(c)
    c.On("Offset", offset).Return(c)
    dao := &dao{c: c}

    _, _, _ = dao.Find(0, offset, conds...)

    c.AssertNumberOfCalls(t, "Find", 1)
    c.AssertNumberOfCalls(t, "Limit", 0)
    c.AssertNumberOfCalls(t, "Offset", 1)
}

func Test_dao_find_when_dont_set_offset_and_set_limit_connection_calls_is_correct(t *testing.T) {
    var conds []interface{}
    limit := 1

    c := &gorm.ConnectionMock{}
    c.On("Limit", limit+1).Return(c)
    c.On("Find", &[]questions.Question{}, conds).Return(c)
    dao := &dao{c: c}

    _, _, _ = dao.Find(limit, 0, conds...)

    c.AssertNumberOfCalls(t, "Find", 1)
    c.AssertNumberOfCalls(t, "Limit", 1)
    c.AssertNumberOfCalls(t, "Offset", 0)
}

func Test_dao_find_when_set_limit_and_offset_connection_calls_is_correct(t *testing.T) {
    var conds []interface{}
    limit := 1
    offset := 1

    c := &gorm.ConnectionMock{}
    c.On("Limit", limit+1).Return(c)
    c.On("Offset", offset).Return(c)
    c.On("Find", &[]questions.Question{}, conds).Return(c)
    dao := &dao{c: c}

    _, _, _ = dao.Find(limit, offset, conds...)

    findCalls := 1
    limitCalls := 1
    offsetCalls := 1
    if !c.AssertNumberOfCalls(t, "Find", findCalls) {
        t.Errorf("Метод Find у connection должен вызваться %d раз", findCalls)
    }
    if !c.AssertNumberOfCalls(t, "Limit", limitCalls) {
        t.Errorf("Метод Limit у connection должен вызваться %d раз", limitCalls)
    }
    if !c.AssertNumberOfCalls(t, "Offset", offsetCalls) {
        t.Errorf("Метод Offset у connection должен вызваться %d раз", offsetCalls)
    }
}

func Test_dao_find_when_we_have_not_more_then_limit_records_in_connection_result_more_is_false(t *testing.T) {
    var conds []interface{}
    limit := 2

    cFind := &gorm.ConnectionMock{}
    cLimit := &gorm.ConnectionMock{}
    cFind.On("Find", &[]questions.Question{}, conds).Return(&gorm.ConnectionMock{}).Run(func(args mock.Arguments) {
        qlOut := args.Get(0).(*[]questions.Question)
        *qlOut = append(*qlOut, questions.Question{ID: 1}, questions.Question{ID: 2})
    })
    cLimit.On("Limit", limit+1).Return(cFind)

    dao := &dao{c: cLimit}

    _, resultMore, _ := dao.Find(limit, 0, conds...)

    assert.Equal(t, false, resultMore, "Возвращаемый more флаг должно быть false")
}

func Test_dao_find_when_we_have_more_then_limit_records_in_connection_result_more_is_true(t *testing.T) {
    var conds []interface{}
    limit := 2

    cFind := &gorm.ConnectionMock{}
    cFind.On("Find", &[]questions.Question{}, conds).Return(&gorm.ConnectionMock{}).Run(func(args mock.Arguments) {
        qlOut := args.Get(0).(*[]questions.Question)
        *qlOut = append(*qlOut, questions.Question{ID: 1}, questions.Question{ID: 2}, questions.Question{ID: 3})
    })
    c := &gorm.ConnectionMock{}
    c.On("Limit", limit+1).Return(cFind)

    dao := &dao{c: c}

    qlResult, resultMore, _ := dao.Find(limit, 0, conds...)

    assert.Equal(t, true, resultMore, "Возвращаемый more флаг должно быть true")
    assert.Equal(t, limit, len(*qlResult), "Лишние объекты question, использовавшиеся для вычисления флага more, должны быть убраны из возвращаемого списка объектов")
}

func Test_dao_find_when_connection_work_wrong_result_error_not_empty_and_have_info_from_connection(t *testing.T) {
    var conds []interface{}
    qlOut := &[]questions.Question{}
    connectionErr := errors.New("Connection mock error")

    c := &gorm.ConnectionMock{}
    c.On("Find", qlOut, conds).Return(&gorm.ConnectionMock{Err: connectionErr})
    dao := &dao{c: c}

    _, _, errResult := dao.Find(0, 0, conds...)

    require.NotNil(t, errResult, "Возвращаемая ошибка не должна быть пустой")
    assert.ErrorIs(t, errResult, connectionErr, "Возвращаемая ошибка должна содержать информацию из connection")
}
