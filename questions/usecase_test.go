package questions

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type daoMock struct {
	mock.Mock
}

func (m *daoMock) Create(q *Question) error {
	args := m.Called(q)
	return args.Error(0)
}

func (m *daoMock) Update(q *Question, fields []string) error {
	args := m.Called(q, fields)
	return args.Error(0)
}

// -------------
// ---- Add ----
// -------------

func Test_usecase_add_dao_calls_is_correct(t *testing.T) {
	qIn := &Question{}

	dao := &daoMock{}
	dao.On("Create", qIn).Return(nil)
	u := usecase{dao: dao}

	_ = u.Add(qIn)

	createCalls := 1
	if !dao.AssertNumberOfCalls(t, "Create", createCalls) {
		t.Errorf("Метод Create у dao должен вызваться %d раз", createCalls)
		t.Fail()
	}
}

func Test_usecase_add_when_dao_work_success_result_error_is_empty(t *testing.T) {
	qIn := &Question{}

	dao := &daoMock{}
	dao.On("Create", qIn).Return(nil)
	u := usecase{dao: dao}

	errResult := u.Add(qIn)

	assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
}

func Test_usecase_add_dao_work_wrong_result_error_not_empty_and_have_info_from_connection(t *testing.T) {
	qIn := &Question{}
	daoErr := errors.New("Dao mock error")

	dao := &daoMock{}
	dao.On("Create", qIn).Return(daoErr)
	u := usecase{dao: dao}

	errResult := u.Add(qIn)

	require.NotNil(t, errResult, "Возвращаемая ошибка не должна быть пустой")
	require.ErrorIs(t, errResult, daoErr, "Возвращаемая ошибка должна содержать информацию из dao")
}

func Test_usecase_add_dao_set_correct_default_values_for_question(t *testing.T) {
	now := time.Now()
	qIn := &Question{
		Step:       255,
		IsFailed:   true,
		RepeatTime: now,
	}

	dao := &daoMock{}
	dao.On("Create", qIn).Return(nil)
	u := usecase{
		dao: dao,
		now: now,
	}

	_ = u.Add(qIn)

	assert.Equal(t, uint8(1), qIn.Step, "Step должен быть 1")
	assert.Equal(t, false, qIn.IsFailed, "Флаг IsFailed должен быть false")
	assert.Equal(t, now.Add(time.Minute*30), qIn.RepeatTime, "RepeatTime должно быть +30 минут от текущего времени")
}

func Test_usecase_add_when_dao_work_wrong_set_original_values_for_question(t *testing.T) {
	now := time.Now()
	daoErr := errors.New("Dao mock error")
	qIn := &Question{
		Step:       255,
		IsFailed:   true,
		RepeatTime: now,
	}

	dao := &daoMock{}
	dao.On("Create", qIn).Return(daoErr)
	u := usecase{
		dao: dao,
		now: now,
	}

	_ = u.Add(qIn)

	assert.Equal(t, uint8(255), qIn.Step, "Step должен быть 1")
	assert.Equal(t, true, qIn.IsFailed, "Флаг IsFailed должен быть false")
	assert.Equal(t, now, qIn.RepeatTime, "RepeatTime должно быть +30 минут от текущего времени")
}

func Test_uses_add_when_dao_work_success_result_question_contains_data_from_dao(t *testing.T) {
	qIn := &Question{ID: 0}

	dao := &daoMock{}
	dao.On("Create", qIn).Return(nil).Run(func(args mock.Arguments) {
		qOut := args.Get(0).(*Question)
		qOut.ID = 1
	})
	u := usecase{dao: dao}

	_ = u.Add(qIn)

	assert.Equal(t, uint64(1), qIn.ID, "Результирующий объект question должен иметь изменения, внесенные в него в dao")
}

// -----------------
// ---- Correct ----
// -----------------

var correctFields = []string{questionGroupId, questionTitle, questionBody}

func Test_usecase_correct_dao_calls_is_correct(t *testing.T) {
	qIn := &Question{}

	dao := &daoMock{}
	dao.On("Update", qIn, correctFields).Return(nil)
	u := usecase{dao: dao}

	_ = u.Correct(qIn)

	updateCalls := 1
	if !dao.AssertNumberOfCalls(t, "Update", updateCalls) {
		t.Errorf("Метод Update у dao должен вызваться %d раз", updateCalls)
		t.Fail()
	}
}

func Test_usecase_correct_when_dao_work_success_result_error_is_empty(t *testing.T) {
	qIn := &Question{}

	dao := &daoMock{}
	dao.On("Update", qIn, correctFields).Return(nil)
	u := usecase{dao: dao}

	errResult := u.Correct(qIn)

	assert.Nil(t, errResult, "Возвращаемая ошибка должна быть пустой")
}

func Test_usecase_correct_dao_work_wrong_result_error_not_empty_and_have_info_from_connection(t *testing.T) {
	qIn := &Question{}
	daoErr := errors.New("Dao mock error")

	dao := &daoMock{}
	dao.On("Update", qIn, correctFields).Return(daoErr)
	u := usecase{dao: dao}

	errResult := u.Correct(qIn)

	require.NotNil(t, errResult, "Возвращаемая ошибка не должна быть пустой")
	require.ErrorIs(t, errResult, daoErr, "Возвращаемая ошибка должна содержать информацию из dao")
}

func Test_uses_correct_when_dao_work_success_result_question_contains_data_from_dao(t *testing.T) {
	qIn := &Question{Title: "Title 1"}

	dao := &daoMock{}
	dao.On("Update", qIn, correctFields).Return(nil).Run(func(args mock.Arguments) {
		qOut := args.Get(0).(*Question)
		qOut.Title = "Title 2"
	})
	u := usecase{dao: dao}

	_ = u.Correct(qIn)

	assert.Equal(t, "Title 2", qIn.Title, "Результирующий объект question должен иметь изменения, внесенные в него в dao")
}
