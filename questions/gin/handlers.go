package gin

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/golobby/container"
    "github.com/pkg/errors"

    "github.com/chudoyoudo/remember-cards/questions"

    errors_formatter "github.com/chudoyoudo/errors-formatter"
    log "github.com/sirupsen/logrus"
)

func RegisterHandlers(r *gin.Engine, middleware ...gin.HandlerFunc) {
    v1 := r.Group("/v1").Use(middleware...)
    v1.POST("/question", addQuestion)
}

type addRequest struct {
    Title   interface{} `json:"title" binding:"required"`
    Body    interface{} `json:"body" binding:"required"`
    GroupId interface{} `json:"groupId" binding:"required,numeric"`
}

func addQuestion(c *gin.Context) {
    response := gin.H{
        "data":   []interface{}{},
        "errors": gin.H{},
    }

    r := &addRequest{}
    if err := c.Bind(r); err != nil {
        response["errors"] = errors_formatter.FormatErrors(err)
        c.Negotiate(http.StatusBadRequest, *getNegotiate(&response))
        return
    }

    uc := getUsecase()
    q := &questions.Question{}
    bindQuestion(q, r)
    err := uc.Add(q)
    if err != nil {
        log.Error(errors.Wrapf(err, "Can't add question via usecase"))
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    response["data"] = *q
    c.Negotiate(http.StatusOK, *getNegotiate(&response))
}

func bindQuestion(q *questions.Question, r *addRequest) {
    title, ok := r.Title.(string)
    if ok {
        q.Title = title
    }

    body, ok := r.Body.(string)
    if ok {
        q.Body = body
    }

    fGroupId, ok := r.GroupId.(float64)
    if ok {
        q.GroupId = uint64(fGroupId)
    } else {
        sGroupId, ok := r.GroupId.(string)
        if ok {
            groupId, err := strconv.Atoi(sGroupId)
            if err != nil || groupId < 0 {
                log.Error(errors.Wrapf(err, "Can't convert question's groupId to int [%s]", sGroupId))
            } else {
                q.GroupId = uint64(groupId)
            }
        }
    }
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
