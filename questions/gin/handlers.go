package gin

import (
    "net/http"
    "strconv"

    rest_api_response_formatter "github.com/chudoyoudo/rest-api-response-formatter"
    "github.com/gin-gonic/gin"
    "github.com/golobby/container"
    "github.com/pkg/errors"

    "github.com/chudoyoudo/remember-cards/questions"

    errors_formatter "github.com/chudoyoudo/errors-formatter"
    log "github.com/sirupsen/logrus"
)

func RegisterHandlers(r *gin.Engine, middleware ...gin.HandlerFunc) {
    v1 := r.Group("/v1").Use(middleware...)
    v1.POST("/question", addHandler)
    v1.GET("/question", listHandler)
    v1.PUT("/question/:id", correctHandler)
    v1.GET("/question/:id", viewHandler)
    v1.DELETE("/question/:id", deleteHandler)
}

type questionData struct {
    Title   string `json:"title" binding:"required"`
    Body    string `json:"body" binding:"required"`
    GroupId uint64 `json:"groupId" binding:"required"`
}

func (d *questionData) Bind(q *questions.Question) {
    if d.Title != "" {
        q.Title = d.Title
    }
    if d.Body != "" {
        q.Body = d.Body
    }
    if d.GroupId != 0 {
        q.GroupId = d.GroupId
    }
}

type filter struct {
    GroupId []uint64 `form:"groupId"`
    UserId  []uint64 `form:"userId"`
    Limit   int      `form:"limit"`
    Offset  int      `form:"offset"`
}

func (f *filter) ToConds() *map[string]interface{} {
    result := map[string]interface{}{}

    if len(f.GroupId) > 0 {
        result[questions.QuestionGroupId] = f.GroupId
    }
    if len(f.UserId) > 0 {
        result[questions.QuestionUserId] = f.UserId
    }

    return &result
}

func listHandler(c *gin.Context) {
    f := &filter{}
    if err := c.ShouldBindQuery(f); err != nil {
        errData := errors_formatter.FormatErrors(err)
        response := rest_api_response_formatter.GetResponseData(&struct{}{}, &errData)
        c.Negotiate(http.StatusBadRequest, *getNegotiate(response))
        return
    }

    conds := f.ToConds()
    order := &[]interface{}{"id desc"}
    uc := getUsecase()
    ql, more, err := getQuestionList(uc, conds, order, f.Limit, f.Offset)
    if err != nil {
        log.Error(errors.Wrap(err, "Can't get question list"))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    response := rest_api_response_formatter.GetResponseData(gin.H{
        "list": *ql,
        "more": more,
    }, &map[string][]string{})
    c.Negotiate(http.StatusOK, *getNegotiate(response))
}

func viewHandler(c *gin.Context) {
    id := getIdFomRequest(c)
    uc := getUsecase()

    q, err := getQuestion(uc, id)
    if err != nil {
        log.Error(errors.Wrap(err, "Can't get question"))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    if nil == q {
        c.AbortWithStatus(http.StatusNotFound)
        return
    }

    response := rest_api_response_formatter.GetResponseData(*q, &map[string][]string{})
    c.Negotiate(http.StatusOK, *getNegotiate(response))
}

func addHandler(c *gin.Context) {
    d := &questionData{}
    if err := c.Bind(d); err != nil {
        errData := errors_formatter.FormatErrors(err)
        response := rest_api_response_formatter.GetResponseData(&struct{}{}, &errData)
        c.Negotiate(http.StatusBadRequest, *getNegotiate(response))
        return
    }

    q := &questions.Question{}
    d.Bind(q)
    uc := getUsecase()
    if err := addQuestion(uc, q); err != nil {
        log.Error(errors.Wrap(err, "Can't add question"))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    response := rest_api_response_formatter.GetResponseData(*q, &map[string][]string{})
    c.Negotiate(http.StatusOK, *getNegotiate(response))
}

func correctHandler(c *gin.Context) {
    id := getIdFomRequest(c)
    uc := getUsecase()

    q, err := getQuestion(uc, id)
    if err != nil {
        log.Error(errors.Wrapf(err, "Can't get question by id %d", id))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    if nil == q {
        c.AbortWithStatus(http.StatusNotFound)
        return
    }

    d := &questionData{}
    if err := c.Bind(d); err != nil {
        errData := errors_formatter.FormatErrors(err)
        response := rest_api_response_formatter.GetResponseData(struct{}{}, &errData)
        c.Negotiate(http.StatusBadRequest, *getNegotiate(response))
        return
    }

    d.Bind(q)
    if err := correctQuestion(uc, q); err != nil {
        log.Error(errors.Wrap(err, "Can't correct question"))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    response := rest_api_response_formatter.GetResponseData(*q, &map[string][]string{})
    c.Negotiate(http.StatusOK, *getNegotiate(response))
}

func deleteHandler(c *gin.Context) {
    id := getIdFomRequest(c)
    uc := getUsecase()

    q, err := getQuestion(uc, id)
    if err != nil {
        log.Error(errors.Wrap(err, "Can't get question"))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    if nil == q {
        c.AbortWithStatus(http.StatusNotFound)
        return
    }

    if err := deleteQuestion(uc, q); err != nil {
        log.Error(errors.Wrap(err, "Can't delete question"))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    response := rest_api_response_formatter.GetResponseData(&struct{}{}, &map[string][]string{})
    c.Negotiate(http.StatusOK, *getNegotiate(response))
}

func getIdFomRequest(c *gin.Context) uint64 {
    idFromUrl := c.Param("id")
    result, err := strconv.ParseUint(idFromUrl, 10, 64)
    if err != nil {
        log.Error(errors.Wrapf(err, "Can't get question id from request [%v]", idFromUrl))
        return 0
    }
    return result
}

func addQuestion(uc questions.Usecase, q *questions.Question) error {
    err := uc.Add(q)
    if err != nil {
        return errors.Wrapf(err, "Can't add question via usecase")
    }
    return nil
}

func correctQuestion(uc questions.Usecase, q *questions.Question) error {
    err := uc.Correct(q)
    if err != nil {
        return errors.Wrapf(err, "Can't correct question via usecase")
    }
    return nil
}

func deleteQuestion(uc questions.Usecase, q *questions.Question) error {
    if q == nil {
        log.Warn(errors.New("Can't delete question. Question is empty"))
        return nil
    }

    err := uc.Delete([]interface{}{"id=?", q.ID})
    if err != nil {
        return errors.Wrapf(err, "Can't delete question by id %d via usecase", q.ID)
    }

    return nil
}

func getQuestion(uc questions.Usecase, id uint64) (*questions.Question, error) {
    ql, _, err := uc.Find(&map[string]interface{}{"id": id}, &[]interface{}{}, 1, 0)
    if err != nil {
        return nil, errors.Wrapf(err, "Can't get question by id %d via usecase", id)
    }

    if len(*ql) == 0 {
        return nil, err
    }

    return &(*ql)[0], nil
}

func getQuestionList(uc questions.Usecase, conds *map[string]interface{}, order *[]interface{}, limit, offset int) (list *[]questions.Question, more bool, err error) {
    ql, more, err := uc.Find(conds, order, limit, offset)
    if err != nil {
        return nil, false, errors.Wrapf(err, "Can't get question list by conds: %v order: %v limit: %d offset: %d via usecase", conds, order, limit, offset)
    }

    return ql, more, err
}

func getNegotiate(data *gin.H) *gin.Negotiate {
    return &gin.Negotiate{
        Offered: []string{gin.MIMEJSON, gin.MIMEXML},
        Data:    data,
    }
}

func getUsecase() questions.Usecase {
    var uc questions.Usecase
    container.Make(&uc)
    return uc
}
