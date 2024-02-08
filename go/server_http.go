package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
)

type Server struct {
	srv http.Server
}

func (s *Server) Serve(listener *net.UnixListener) {
	s.srv.Handler = s
	s.srv.Serve(listener)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Print("request: %s", req)
	w.Header().Set("Content-Type", "application/json")

	dec := json.NewDecoder(req.Body)
	enc := json.NewEncoder(w)

	var err error

	var m Message

	err = dec.Decode(&m.Content)

	if err != nil {
		panic(err)
	}

	_, err = m.WriteToChrome()

	if err != nil {
		panic(err)
	}

	_, err = m.ReadFromChrome()

	if err != nil && err != io.EOF {
		panic(err)
	}

	err = enc.Encode(m.Content)

	if err != nil {
		panic(err)
	}
}
