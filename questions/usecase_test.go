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
