package voicebr

import (
	"net/http"
)

type Client struct {
	c *http.Client
}
