package gin

import (
    "testing"

    "github.com/pkg/errors"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "github.com/chudoyoudo/remember-cards/questions"
)

type usecaseMock struct {
    mock.Mock
}

func (m *usecaseMock) Add(q *questions.Question) error {
    args := m.Called(q)
    return args.Error(0)
}

func (m *usecaseMock) Correct(q *questions.Question) error {
    args := m.Called(q)
    return args.Error(0)
}

func (m *usecaseMock) Delete(conds []interface{}) error {
    args := m.Called(conds)
    return args.Error(0)
}

func (m *usecaseMock) Find(conds *map[string]interface{}, order *[]interface{}, limit, offset int) (list *[]questions.Question, more bool, err error) {
    args := m.Called(conds, order, limit, offset)
    return args.Get(0).(*[]questions.Question), args.Bool(1), args.Error(2)
}

//-----------
//--- Add ---
//-----------

func Test_handler_add_usecase_calls_is_correct(t *testing.T) {
    qIn := &questions.Question{}

    uc := &usecaseMock{}
    uc.On("Add", qIn).Return(nil)

    _ = addQuestion(uc, qIn)

    addCalls := 1
    if !uc.AssertNumberOfCalls(t, "Add", addCalls) {
        t.Errorf("Метод Add у usecase должен вызваться %d раз", addCalls)
        t.Fail()
    }
}

func Test_handler_add_when_usecase_work_success_result_error_is_empty(t *testing.T) {
    qIn := &questions.Question{}

    uc := &usecaseMock{}
    uc.On("Add", qIn).Return(nil)

    errResult := addQuestion(uc, qIn)

    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
}

func Test_handler_add_usecase_work_wrong_result_error_not_empty_and_have_info_from_usecase(t *testing.T) {
    qIn := &questions.Question{}
    usecaseErr := errors.New("Usecase mock error")

    uc := &usecaseMock{}
    uc.On("Add", qIn).Return(usecaseErr)

    errResult := addQuestion(uc, qIn)

    require.NotNil(t, errResult, "Возвращаемая ошибка не должна быть пустой")
    require.ErrorIs(t, errResult, usecaseErr, "Возвращаемая ошибка должна содержать информацию из usecase")
}

func Test_handler_add_when_usecase_work_success_result_question_contains_data_from_usecase(t *testing.T) {
    qIn := &questions.Question{ID: 0}

    uc := &usecaseMock{}
    uc.On("Add", qIn).Return(nil).Run(func(args mock.Arguments) {
        qOut := args.Get(0).(*questions.Question)
        qOut.ID = 1
    })

    _ = addQuestion(uc, qIn)

    assert.Equal(t, uint64(1), qIn.ID, "Результирующий объект question должен иметь изменения, внесенные в него в usecase")
}

//---------------
//--- Correct ---
//---------------

func Test_handler_correct_usecase_calls_is_correct(t *testing.T) {
    qIn := &questions.Question{}

    uc := &usecaseMock{}
    uc.On("Correct", qIn).Return(nil)

    _ = correctQuestion(uc, qIn)

    correctCalls := 1
    if !uc.AssertNumberOfCalls(t, "Correct", correctCalls) {
        t.Errorf("Метод Correct у usecase должен вызваться %d раз", correctCalls)
        t.Fail()
    }
}

func Test_handler_correct_when_usecase_work_success_result_error_is_empty(t *testing.T) {
    qIn := &questions.Question{}

    uc := &usecaseMock{}
    uc.On("Correct", qIn).Return(nil)

    errResult := correctQuestion(uc, qIn)

    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
}

func Test_handler_correct_usecase_work_wrong_result_error_not_empty_and_have_info_from_usecase(t *testing.T) {
    qIn := &questions.Question{}
    usecaseErr := errors.New("Usecase mock error")

    uc := &usecaseMock{}
    uc.On("Correct", qIn).Return(usecaseErr)

    errResult := correctQuestion(uc, qIn)

    require.NotNil(t, errResult, "Возвращаемая ошибка не должна быть пустой")
    require.ErrorIs(t, errResult, usecaseErr, "Возвращаемая ошибка должна содержать информацию из usecase")
}

func Test_handler_correct_when_usecase_work_success_result_question_contains_data_from_usecase(t *testing.T) {
    qIn := &questions.Question{Title: "Title 1"}

    uc := &usecaseMock{}
    uc.On("Correct", qIn).Return(nil).Run(func(args mock.Arguments) {
        qOut := args.Get(0).(*questions.Question)
        qOut.Title = "Title 2"
    })

    _ = correctQuestion(uc, qIn)

    assert.Equal(t, "Title 2", qIn.Title, "Результирующий объект question должен иметь изменения, внесенные в него в usecase")
}

//--------------
//--- Delete ---
//--------------

func Test_handler_delete_usecase_calls_is_correct(t *testing.T) {
    qIn := &questions.Question{ID: 1}
    conds := []interface{}{"id=?", qIn.ID}

    uc := &usecaseMock{}
    uc.On("Delete", conds).Return(nil)

    _ = deleteQuestion(uc, qIn)

    deleteCalls := 1
    if !uc.AssertNumberOfCalls(t, "Delete", deleteCalls) {
        t.Errorf("Метод Delete у usecase должен вызваться %d раз", deleteCalls)
        t.Fail()
    }
}

func Test_handler_delete_when_usecase_work_success_result_error_is_empty(t *testing.T) {
    qIn := &questions.Question{ID: 1}
    conds := []interface{}{"id=?", qIn.ID}

    uc := &usecaseMock{}
    uc.On("Delete", conds).Return(nil)

    errResult := deleteQuestion(uc, qIn)

    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
}

func Test_handler_delete_usecase_work_wrong_result_error_not_empty_and_have_info_from_usecase(t *testing.T) {
    qIn := &questions.Question{ID: 1}
    conds := []interface{}{"id=?", qIn.ID}
    usecaseErr := errors.New("Usecase mock error")

    uc := &usecaseMock{}
    uc.On("Delete", conds).Return(usecaseErr)

    errResult := deleteQuestion(uc, qIn)

    require.NotNil(t, errResult, "Возвращаемая ошибка не должна быть пустой")
    require.ErrorIs(t, errResult, usecaseErr, "Возвращаемая ошибка должна содержать информацию из usecase")
}

//------------
//--- View ---
//------------

func Test_handler_view_usecase_calls_is_correct(t *testing.T) {
    id := uint64(1)
    conds := &map[string]interface{}{"id": id}

    uc := &usecaseMock{}
    uc.On("Find", conds, &[]interface{}{}, 1, 0).Return(&[]questions.Question{}, false, nil)

    _, _ = getQuestion(uc, id)

    findCalls := 1
    if !uc.AssertNumberOfCalls(t, "Find", findCalls) {
        t.Errorf("Метод Find у usecase должен вызваться %d раз", findCalls)
        t.Fail()
    }
}

func Test_handler_view_when_usecase_work_success_result_error_is_empty(t *testing.T) {
    id := uint64(1)
    conds := &map[string]interface{}{"id": id}

    uc := &usecaseMock{}
    uc.On("Find", conds, &[]interface{}{}, 1, 0).Return(&[]questions.Question{}, false, nil)

    _, errResult := getQuestion(uc, id)

    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
}

func Test_handler_view_usecase_work_wrong_result_error_not_empty_and_have_info_from_usecase(t *testing.T) {
    usecaseErr := errors.New("Usecase mock error")
    id := uint64(1)
    conds := &map[string]interface{}{"id": id}

    uc := &usecaseMock{}
    uc.On("Find", conds, &[]interface{}{}, 1, 0).Return(&[]questions.Question{}, false, usecaseErr)

    _, errResult := getQuestion(uc, id)

    require.NotNil(t, errResult, "Возвращаемая ошибка не должна быть пустой")
    require.ErrorIs(t, errResult, usecaseErr, "Возвращаемая ошибка должна содержать информацию из usecase")
}

func Test_handler_view_when_usecase_work_success_result_question_contains_data_from_usecase(t *testing.T) {
    qExpected := &questions.Question{Title: "Title 1"}
    id := uint64(1)
    conds := &map[string]interface{}{"id": id}

    uc := &usecaseMock{}
    uc.On("Find", conds, &[]interface{}{}, 1, 0).Return(&[]questions.Question{*qExpected}, false, nil)

    qResult, _ := getQuestion(uc, id)

    assert.Equal(t, *qExpected, *qResult, "Результирующий объект question должен быть идентичен тому, что вернул usecase")
}

//------------
//--- Find ---
//------------

func Test_handler_list_usecase_calls_is_correct(t *testing.T) {
    conds := &map[string]interface{}{"id": 1}
    order := &[]interface{}{"id asc"}
    limit := 1
    offset := 1

    uc := &usecaseMock{}
    uc.On("Find", conds, order, limit, offset).Return(&[]questions.Question{}, false, nil)

    _, _, _ = getQuestionList(uc, conds, order, limit, offset)

    findCalls := 1
    if !uc.AssertNumberOfCalls(t, "Find", findCalls) {
        t.Errorf("Метод Find у usecase должен вызваться %d раз", findCalls)
        t.Fail()
    }
}

func Test_handler_list_when_usecase_work_success_result_error_is_empty(t *testing.T) {
    conds := &map[string]interface{}{"id": 1}
    order := &[]interface{}{"id asc"}
    limit := 1
    offset := 1

    uc := &usecaseMock{}
    uc.On("Find", conds, order, limit, offset).Return(&[]questions.Question{}, false, nil)

    _, _, errResult := getQuestionList(uc, conds, order, limit, offset)

    assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
}

func Test_handler_list_usecase_work_wrong_result_error_not_empty_and_have_info_from_usecase(t *testing.T) {
    usecaseErr := errors.New("Usecase mock error")
    conds := &map[string]interface{}{"id": 1}
    order := &[]interface{}{"id asc"}
    limit := 1
    offset := 1

    uc := &usecaseMock{}
    uc.On("Find", conds, order, limit, offset).Return(&[]questions.Question{}, false, usecaseErr)

    _, _, errResult := getQuestionList(uc, conds, order, limit, offset)

    require.NotNil(t, errResult, "Возвращаемая ошибка не должна быть пустой")
    require.ErrorIs(t, errResult, usecaseErr, "Возвращаемая ошибка должна содержать информацию из usecase")
}

func Test_handler_list_when_usecase_work_success_result_question_contains_data_from_usecase(t *testing.T) {
    qlExpected := &[]questions.Question{{ID: 1}, {ID: 2}, {ID: 3}}
    moreExpected := true
    conds := &map[string]interface{}{}
    order := &[]interface{}{}
    limit := 3
    offset := 1

    uc := &usecaseMock{}
    uc.On("Find", conds, order, limit, offset).Return(qlExpected, moreExpected, nil)

    qlResult, moreResult, _ := getQuestionList(uc, conds, order, limit, offset)

    assert.Equal(t, *qlExpected, *qlResult, "Результирующий список объект question должен быть идентичен тому, что вернул usecase")
    assert.Equal(t, moreExpected, moreResult, "Результирующий флаг more должен быть идентичен тому, что вернул usecase")
}

//--------------
//--- Filter ---
//--------------

func Test_filter_to_conds_retern_correct_conditions(t *testing.T) {
    groupIdList := []uint64{1, 2, 3}
    userIdList := []uint64{3, 4, 5}
    condsExpected := &map[string]interface{}{
        questions.QuestionGroupId: groupIdList,
        questions.QuestionUserId:  userIdList,
    }
    f := &filter{
        GroupId: groupIdList,
        UserId:  userIdList,
    }

    condsResult := f.ToConds()

    assert.Equal(t, *condsExpected, *condsResult, "Результирующий список кондишенов неверный")
}

//---------------------
//--- Question data ---
//---------------------

func Test_question_data_bind_retern_correct_question_object(t *testing.T) {
    groupId := uint64(1)
    title := "Title"
    body := "body"

    qIn := &questions.Question{}
    qExpected := &questions.Question{
        Title:   title,
        Body:    body,
        GroupId: groupId,
    }

    qd := &questionData{
        Title:   title,
        Body:    body,
        GroupId: groupId,
    }
    qd.Bind(qIn)

    assert.Equal(t, *qExpected, *qIn, "Результирующий объект question неверный")
}
