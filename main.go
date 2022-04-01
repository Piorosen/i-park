package main

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"net/http"
	"os"

	"C"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	date := C.__DATE__
	time := C.__TIME__
	log.Println("BUILD DATE : [ " + date + " " + time + " ]")

	log.Println("program start")
	deviceId := os.Getenv("DEVICE_ID")
	run_token := os.Getenv("TOKEN")
	check_token := os.Getenv("RUN_TOKEN")

	log.Println("enviroment load")
	log.Println()
	// req_token := os.Getenv("REQ_TOKEN")
	token := sha512.Sum512([]byte(deviceId + "-" + run_token))
	hexToken := hex.EncodeToString(token[:])

	if hexToken != check_token {
		log.Println("token expired")
		return
	}
	token = sha512.Sum512([]byte(check_token + "-" + date + "-" + time))
	check_token = hex.EncodeToString(token[:])
	log.Println("Request Token : ", check_token)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	log.Println("toggle on")

	toggle := true
	// Routes
	e.POST("/on", func(c echo.Context) error {
		log.Println("toggle on")
		room := c.Request().Header.Get("room")

		if c.Request().Header.Get("access-token") == check_token {
			log.Println("token check success")
			err := SetLight(deviceId, true, room)
			if err != nil {
				log.Println("token check fail", err)
				return c.String(http.StatusOK, "{\"result\": false}")
			} else {
				log.Println("token check success")
				toggle = true
				return c.String(http.StatusOK, "{\"result\": true}")
			}
		} else {
			log.Println("token check fail")
			return c.String(http.StatusOK, "{\"result\": false}")
		}
	})
	e.POST("/off", func(c echo.Context) error {
		log.Println("toggle on")
		room := c.Request().Header.Get("room")

		if c.Request().Header.Get("access-token") == check_token {
			log.Println("token check success")
			err := SetLight(deviceId, false, room)
			if err != nil {
				log.Println("token check fail", err)
				return c.String(http.StatusOK, "{\"result\": false}")
			} else {
				log.Println("token check success")
				toggle = false
				return c.String(http.StatusOK, "{\"result\": true}")
			}
		} else {
			log.Println("token check fail")
			return c.String(http.StatusOK, "{\"result\": false}")
		}
	})
	e.POST("/sequence", func(c echo.Context) error {
		if c.Request().Header.Get("access-token") == check_token {
			log.Println("token check success")
			room := c.Request().Header.Get("room")

			err := SetLight(deviceId, toggle, room)
			if err != nil {
				log.Println("token check fail", err)
				return c.String(http.StatusOK, "{\"result\": false}")
			} else {
				log.Println("token check success")
				textt := "true"
				if !toggle {
					textt = "false"
				}
				toggle = !toggle
				return c.String(http.StatusOK, "{\"result\": true, \"light\": "+textt+"}")
			}
		} else {
			log.Println("token check fail")
			return c.String(http.StatusOK, "{\"result\": false}")
		}
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
