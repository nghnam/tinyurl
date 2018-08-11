package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nghnam/tinyurl/config"
	"github.com/nghnam/tinyurl/handler"
	"github.com/nghnam/tinyurl/keyservice"
	"github.com/nghnam/tinyurl/store"
)

func main() {
	e := echo.New()
	cfg, err := config.NewConfig("config.yaml")
	if err != nil {
		e.Logger.Fatal(err)
	}

	keyCh := make(chan chan string)
	keyService := keyservice.NewKeyService(keyCh)
	keyService.CreateKeys(cfg.AmountOfKey, cfg.LengthOfKey)
	go keyService.Run()

	db := store.NewMapDB()
	h := handler.NewHandler(keyCh, db, cfg.DomainURL)

	e.Use(middleware.Logger())
	e.POST("/create", h.Create)
	e.GET("/:key", h.Redirect)
	e.Logger.Fatal(e.Start(":1323"))
}
