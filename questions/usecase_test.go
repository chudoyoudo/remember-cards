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

func Test_usecase_add_question_success(t *testing.T) {
    now := time.Now()
    qExpected := &Question{ID: 1}
    qIn := &Question{
        ID:         0,
        Step:       255,
        IsFailed:   true,
        RepeatTime: now,
    }

    dao := &daoMock{}
    u := usecase{
        dao: dao,
        now: now,
    }
    dao.On("Create", qIn).Return(nil).Run(func(args mock.Arguments) {
        qIn := args.Get(0).(*Question)
        qIn.ID = 1
    })

    errResult := u.Add(qIn)

    dao.AssertNumberOfCalls(t, "Create", 1)
    assert.Nil(t, nil, errResult, "Возвращаемая ошибка должна быть пустой")
    assert.Equal(t, qExpected.ID, qIn.ID, "ID должен быть заполнен")
    assert.Equal(t, uint8(1), qIn.Step, "Step должен быть 1")
    assert.Equal(t, false, qIn.IsFailed, "Флаг IsFailed должен быть false")
    assert.Equal(t, now.Add(time.Minute*30), qIn.RepeatTime, "RepeatTime должно быть +30 минут от текущего времени")
}

func Test_usecase_add_question_when_dao_failed(t *testing.T) {
    qIn := &Question{ID: 0}
    daoErr := errors.New("Dao mock error")

    dao := &daoMock{}
    dao.On("Create", qIn).Return(daoErr)
    u := usecase{dao: dao}

    errResult := u.Add(qIn)

    dao.AssertNumberOfCalls(t, "Create", 1)
    require.ErrorIs(t, errResult, daoErr, "Результирующая ошибка не содержит информации об ошибках из dao")
}
