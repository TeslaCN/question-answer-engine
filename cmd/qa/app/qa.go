package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"question-answer-engine/cmd/qa/app/config"
	"question-answer-engine/internal/app/qa/service"
	"strconv"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/api/qa", func(context *gin.Context) {
		query := context.Query("q")
		size, _ := strconv.Atoi(context.Query("size"))
		preTag := context.Query("preTag")
		postTag := context.Query("postTag")
		if size < 1 {
			size = 10
		}
		qa := service.MultiMatchQa(query, nil, size, preTag, postTag)
		context.JSON(http.StatusOK, qa)
	})
	r.PUT("/api/qa", func(context *gin.Context) {
		requestBody := &struct {
			Question     []string `json:"question"`
			Answer       []string `json:"answer"`
			Tag          []string `json:"tag"`
			QuestionTime uint64   `json:"questionTime"`
			AnswerTime   uint64   `json:"answerTime"`
		}{}
		if e := context.ShouldBindJSON(requestBody); e != nil {
		}
		service.Put(requestBody)
	})
	return r
}

func Run() {
	r := setupRouter()
	_ = r.Run(config.Config.Server)
}
