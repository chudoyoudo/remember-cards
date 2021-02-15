package questions

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_question_to_map_return_all_fields_if_fieldst_list_in_params_is_empty(t *testing.T) {
	rt := time.Now()
	q := &Question{
		ID:         1,
		UserId:     2,
		GroupId:    3,
		Title:      "Title",
		Body:       "Body",
		Step:       4,
		RepeatTime: rt,
		IsFailed:   true,
	}

	expectedMap := map[string]interface{}{
		QuestionUserId:     uint64(2),
		QuestionGroupId:    uint64(3),
		questionTitle:      "Title",
		questionBody:       "Body",
		questionStep:       uint8(4),
		questionRepeatTime: rt,
		questionIsFailed:   true,
	}
	resultMap := q.ToMap([]string{})

	assert.Equal(t, expectedMap, *resultMap, "Возвращаемая мапа не содержит все необходимые дданные")
}

func Test_question_to_map_return_only_requested_fields_data(t *testing.T) {
	rt := time.Now()
	q := &Question{
		ID:         1,
		UserId:     2,
		GroupId:    3,
		Title:      "Title",
		Body:       "Body",
		Step:       4,
		RepeatTime: rt,
		IsFailed:   true,
	}

	expectedMap := map[string]interface{}{
		QuestionUserId:     uint64(2),
		QuestionGroupId:    uint64(3),
		questionTitle:      "Title",
		questionBody:       "Body",
		questionStep:       uint8(4),
		questionRepeatTime: rt,
		questionIsFailed:   true,
	}
	resultMap := q.ToMap([]string{QuestionUserId, QuestionGroupId, questionTitle, questionBody, questionStep, questionRepeatTime, questionIsFailed})

	assert.Equal(t, expectedMap, *resultMap, "Возвращаемая мапа не содержит все необходимые дданные")
}
