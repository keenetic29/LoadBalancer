package domain

import (
	"net/url"
	"time"
)

type Server struct {
	URL     *url.URL
	Healthy bool
}

func NewServer(rawurl string) (*Server, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &Server{URL: u, Healthy: true}, nil
}

type Client struct {
	ID           string
	Capacity     int
	RatePerSec   int
	RefillPeriod time.Duration
}

func NewClient(id string, capacity, ratePerSec int) *Client {
	return &Client{
		ID:           id,
		Capacity:     capacity,
		RatePerSec:   ratePerSec,
		RefillPeriod: time.Second, // default refill period
	}
}

