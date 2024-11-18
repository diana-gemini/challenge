package controllers

import (
	"github.com/diana-gemini/doodocs/internal/controllers/archive"
	"github.com/diana-gemini/doodocs/internal/controllers/mail"
	"github.com/diana-gemini/doodocs/pkg"
	"github.com/gin-gonic/gin"
)

func Setup(app pkg.Application, router *gin.Engine) {
	controllers := SetupControllers(app)

	archiveRouter := router.Group("/api/archive")
	{
		archiveRouter.POST("/information", controllers.ArchiveController.ArchiveInformation)
		archiveRouter.POST("/files", controllers.ArchiveController.CreateArchive)
	}

	mailRouter := router.Group("/api/mail")
	{
		mailRouter.POST("/file", controllers.MailController.SendEmail)
	}
}

type controllers struct {
	ArchiveController *archive.ArchiveController
	MailController    *mail.MailController
}

func SetupControllers(app pkg.Application) controllers {
	controllers := controllers{}
	controllers.ArchiveController = &archive.ArchiveController{
		Env: app.Env,
	}
	controllers.MailController = &mail.MailController{
		Env: app.Env,
	}
	return controllers
}
