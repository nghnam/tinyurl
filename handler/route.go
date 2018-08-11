package handler

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/nghnam/tinyurl/store"
)

type Handler struct {
	keyCh     chan chan string
	db        *store.MapDB
	domainUrl string
}

func NewHandler(keyCh chan chan string, db *store.MapDB, domainUrl string) *Handler {
	return &Handler{
		keyCh:     keyCh,
		db:        db,
		domainUrl: domainUrl,
	}
}

func (h *Handler) Create(c echo.Context) error {
	// Validate input
	url := c.FormValue("url")
	if !govalidator.IsRequestURL(url) {
		return c.String(http.StatusBadRequest, "URL is not valid")
	}

	// Connect to Key Generation service
	reqCh := make(chan string)
	h.keyCh <- reqCh
	key, ok := <-reqCh
	if !ok {
		return c.String(http.StatusNoContent, "")
	}

	// Update db
	h.db.Update(key, url)

	return c.String(http.StatusCreated, fmt.Sprintf("%v/%v", h.domainUrl, key))
}

func (h *Handler) Redirect(c echo.Context) error {
	key := c.Param("key")
	url, err := h.db.Lookup(key)
	if err != nil {
		return c.String(http.StatusNotFound, "Key not found")
	}
	return c.Redirect(http.StatusMovedPermanently, url)
}
