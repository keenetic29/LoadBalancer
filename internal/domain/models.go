package domain

import "net/url"

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