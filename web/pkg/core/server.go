package core

import "os"

type Server interface {
	Start() error
}

func Run(srv Server) {
	if err := srv.Start(); err != nil {
		os.Exit(2)
		return
	}
}
