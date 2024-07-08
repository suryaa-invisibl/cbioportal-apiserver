package main

import (
	"net/http"
	"os"

	"github.com/invisibl-cloud/cbioportal-apiserver/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	g := e.Group(os.Getenv("CONTEXT_PATH"))

	g.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// cbioportal apis
	g.POST("/api/upload", func(c echo.Context) error {
		err := routes.Upload(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, "success")
	})
	g.POST("/api/studies/upload", func(c echo.Context) error {
		err := routes.Upload(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, "success")
	})
	g.GET("/api/studies/get-filters", func(c echo.Context) error {
		data, err := routes.GetFilters(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"data": data})
	})
	g.GET("/api/studies/apply-filters", func(c echo.Context) error {
		data, err := routes.GetStudiesWithFilters(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"data": data})
	})
	g.POST("/api/studies/apply-filters", func(c echo.Context) error {
		data, err := routes.GetStudiesWithFiltersV2(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]any{"data": data})
	})

	e.Logger.Fatal(e.Start(":9000"))
}
