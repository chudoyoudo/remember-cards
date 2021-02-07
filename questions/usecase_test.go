package questions

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

type daoMock struct {
    mock.Mock
}

func (m *daoMock) Create(q *Question) (*Question, error) {
    args := m.Called(q)
    return args.Get(0).(*Question), args.Error(1)
}

func TestAddQuestionWithDaoSuccess(t *testing.T) {
    qIn := &Question{ID: 0}
    qOut := &Question{ID: 1}

    dao := &daoMock{}
    dao.On("Create", qIn).Return(qOut, nil)
    u := usecase{dao: dao}

    qResult, errResult := u.Add(qIn)

    require.Equal(t, nil, errResult, "dao вернул пустую ошибку, однако usecase вернул не пустую ошибку")
    require.Equal(t, qOut, qResult, "usecase вернул не тот объект question, который вернул dao")
    dao.AssertCalled(t, "Create", qIn)
}

func TestAddQuestionWithDaoFail(t *testing.T) {
    var qEmpty *Question
    qIn := &Question{ID: 0}
    errOut := errors.New("Test error")

    dao := &daoMock{}
    dao.On("Create", qIn).Return(qEmpty, errOut)
    u := usecase{dao: dao}

    qResult, errResult := u.Add(qIn)

    require.Equal(t, errOut, errResult, "usecase вернул не ту ошибку, которую вернул dao")
    require.Equal(t, qEmpty, qResult, "usecase вернул не пустой объект question")
    dao.AssertCalled(t, "Create", qIn)
}
