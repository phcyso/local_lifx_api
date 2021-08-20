package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phcyso/local_lifx_api/lights"
)

/* Global Config vars */

/** Defaults */
var port string = "7070"
var configPath string = "."
var lightsRefreshTime int = 60

/* loader */
func loadConfig() {
	cp, exists := os.LookupEnv("CONFIG_PATH")
	if exists {
		configPath = cp
	}

	p, exists := os.LookupEnv("PORT")
	if exists {
		port = p
	}

	r, exists := os.LookupEnv("LIGHTS_REFRESH_TIME")
	if exists {
		tI, err := strconv.Atoi(r)
		if err != nil {
			panic("Not able to parse LIGHTS_REFRESH_TIME, must be a positive integer")
		}
		lightsRefreshTime = tI
	}

}

/* Light setup related functions */

func loadLights() {
	err := lights.InitScenes(configPath)
	if err != nil {
		panic(err)
	}

	err = lights.LoadLights() // TODO Async this?
	if err != nil {
		log.Printf("error loading one or more lights:  %v", err)
	}

	/* Refresh the lights on a schedule */
	ticker := time.NewTicker(time.Duration(time.Duration(lightsRefreshTime)) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				lights.RefreshLights()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

/* Main entrypoint */
func main() {

	loadConfig()
	loadLights()

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	// Routes
	e.GET("/lights/list", listLights)
	e.GET("/lights/all/off", allLightsOff)
	e.GET("/lights/all/on", allLightsOn)
	e.GET("/light/off/:mac", lightOff)
	e.GET("/light/on/:mac", lightOn)
	e.GET("/light/refresh/:mac", refreshLight)

	e.GET("/scenes/list", listscenes)
	e.GET("/scene/run/:id", triggerScene)
	e.POST("/scene/save", saveScene)
	e.POST("/scene/delete/:id", deleteScene)
	e.POST("/scene/modify", modifyScene)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}

/* Routes to the various light functions */
func listLights(c echo.Context) error {
	return c.JSON(http.StatusOK, lights.ListAllLights())
}

func refreshLight(c echo.Context) error {
	return c.JSON(http.StatusOK, lights.RefreshLight(c.Param("mac")))
}

func lightOn(c echo.Context) error {
	return c.JSON(http.StatusOK, lights.LightOn(c.Param("mac")))
}

func lightOff(c echo.Context) error {
	return c.JSON(http.StatusOK, lights.LightOff(c.Param("mac")))
}
func allLightsOff(c echo.Context) error {
	return c.JSON(http.StatusOK, lights.AllLightsOff())
}

func allLightsOn(c echo.Context) error {
	return c.JSON(http.StatusOK, lights.AllLightsOn())
}

func listscenes(c echo.Context) error {
	return c.JSON(http.StatusOK, lights.ListAllScenes())
}

func triggerScene(c echo.Context) error {
	return c.JSON(http.StatusOK, lights.TriggerScene(c.Param("id")))
}

// TODO, these two are basically the same. could dry these two up if touching it anyway.
func saveScene(c echo.Context) error {
	s := new(lights.SceneSaveRequest)
	if err := c.Bind(s); err != nil {
		log.Printf("error unmarshaling scene save request: %v", err)
		return c.JSON(http.StatusBadRequest, fmt.Errorf("bad save request"))
	}
	if err := lights.SaveScene(*s); err != nil {
		c.Logger().Error("error saving scene save request: %v", err)

		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, "OK")
}
func modifyScene(c echo.Context) error {

	s := new(lights.SceneSaveRequest)
	if err := c.Bind(s); err != nil {
		log.Printf("error unmarshaling scene save request: %v", err)
		return c.JSON(http.StatusBadRequest, fmt.Errorf("bad save request"))
	}

	if err := lights.ModifyScene(*s); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, "OK")
}

func deleteScene(c echo.Context) error {
	c.Logger().Info("trying to delete '%v'", c.Param("id"))

	if err := lights.DeleteScene(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, "OK")
}
