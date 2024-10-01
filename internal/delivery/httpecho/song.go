package httpecho

import (
	_ "github.com/Magic-Kot/effective-mobile/docs"
	"github.com/Magic-Kot/effective-mobile/internal/controllers"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SetSongRoutes(e *echo.Echo, apiController *controllers.ApiController) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	song := e.Group("/song")
	{
		song.POST("/create", apiController.AddSong)
		song.GET("/all", apiController.GetAllSong)
		song.GET("/get/:id", apiController.GetLyricsSong)
		song.PATCH("/update/:id", apiController.UpdateSong)
		song.DELETE("/delete/:id", apiController.DeleteSong)
	}
}
