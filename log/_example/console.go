package main

import "gitlab.com/cinemae/gopkg/log"

func main() {
	log.Info("test")
	log.Errorf("test:%d", 100)

	log.WithFields("name", "phachon").Info("with")
}
