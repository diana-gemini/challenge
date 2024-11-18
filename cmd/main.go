package main

import (
	"fmt"

	"github.com/diana-gemini/doodocs/internal/controllers"
	"github.com/diana-gemini/doodocs/pkg"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Info("application is running")

	app, err := pkg.App()

	if err != nil {
		log.Fatal(err)
	}

	ginRouter := gin.Default()

	controllers.Setup(app, ginRouter)

	ginRouter.Run(fmt.Sprintf(":%s", app.Env.PORT))
}
