package main

type DevopsServer struct {
}

func NewDevopsServer(configFile string) *DevopsServer {
	server := &DevopsServer{}
	return server
}

func (s *DevopsServer) Run() error {
	return nil
}
