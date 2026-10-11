package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"volleyreg/internal/app/handler"
	"volleyreg/internal/app/repository"
	"volleyreg/internal/app/storage"
)

func StartServer(repo *repository.Repository, store *storage.Storage) {
	logrus.Info("Server start up")

	h := handler.NewHandler(repo, store)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/styles", "./resources/styles")

	// HTML-страницы (лабораторные 1-2)
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/participants")
	})
	r.GET("/participants", h.GetParticipants)
	r.GET("/participants/:id", h.GetParticipant)
	r.POST("/participants/:id/add-to-team-application", h.AddParticipantToDraft)
	r.GET("/team-applications/:id", h.GetTeamApplication)
	r.POST("/team-applications/:id/delete", h.DeleteTeamApplication)

	// REST API веб-сервиса (лабораторная 3)
	api := r.Group("/api")
	{
		// домен участников (услуги)
		participants := api.Group("/participants")
		participants.GET("", h.ApiGetParticipants)
		participants.GET("/:id", h.ApiGetParticipant)
		participants.POST("", h.ApiCreateParticipant)

		// домен заявок команд
		apps := api.Group("/team-applications")
		apps.GET("/draft", h.ApiGetDraft)
		apps.GET("", h.ApiGetTeamApplications)
		apps.GET("/:id", h.ApiGetTeamApplication)
		apps.PUT("/:id", h.ApiUpdateTeamApplication)
		apps.PUT("/:id/form", h.ApiFormTeamApplication)
		apps.PUT("/:id/moderate", h.ApiModerateTeamApplication)
		apps.DELETE("/:id", h.ApiDeleteTeamApplication)

		// домен м-м: участники в составе заявки
		apps.POST("/draft/participants/:participant_id", h.ApiAddParticipantToDraft)
		apps.PUT("/:id/participants/:participant_id", h.ApiUpdateJerseyNumber)
		apps.DELETE("/:id/participants/:participant_id", h.ApiRemoveParticipant)

		// домен пользователей
		users := api.Group("/users")
		users.POST("/register", h.ApiRegister)
		users.POST("/login", h.ApiLogin)
		users.POST("/logout", h.ApiLogout)
	}

	if err := r.Run(":8080"); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server down")
}
