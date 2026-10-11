package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"volleyreg/internal/app/handler"
	"volleyreg/internal/app/repository"
)

func StartServer() {
	logrus.Info("Server start up")

	repo := repository.NewRepository()
	h := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/styles", "./resources/styles")

	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/participants")
	})
	r.GET("/participants", h.GetParticipants)
	r.GET("/participants/:id", h.GetParticipant)
	r.GET("/team-applications/:id", h.GetTeamApplication)

	if err := r.Run(":8080"); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server down")
}
