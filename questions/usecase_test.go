package questions

import (
    "testing"

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

func Test_usecase_add_question_when_dao_ok(t *testing.T) {
    qIn := &Question{ID: 0}
    qExpected := &Question{ID: 1}

    dao := &daoMock{}
    u := usecase{dao: dao}
    dao.On("Create", qIn).Return(nil).Run(func(args mock.Arguments) {
        qIn := args.Get(0).(*Question)
        qIn.ID = 1
    })

    errResult := u.Add(qIn)

    dao.AssertNumberOfCalls(t, "Create", 1)
    assert.Nil(t, nil, errResult, "Ошибка должна быть пустой")
    assert.Equal(t, qExpected, qIn, "Объект вопроса не содержит изменений, которые в него внес dao")
}

func Test_usecase_add_question_when_dao_failed(t *testing.T) {
    qIn := &Question{ID: 0}
    daoErr := errors.New("Dao mock error")

    dao := &daoMock{}
    dao.On("Create", qIn).Return(daoErr)
    u := usecase{dao: dao}

    errResult := u.Add(qIn)

    dao.AssertNumberOfCalls(t, "Create", 1)
    require.ErrorIs(t, errResult, daoErr, "Результирующая ошибка не содержащая информации об ошибках из dao")
}
