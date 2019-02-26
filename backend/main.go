package main

import (
	"context"
	"flag"
	"fmt"
	"go-say/api"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
)

func main() {
	port := flag.Int("p", 8080, "port")
	flag.Parse()

	log.Printf("listening to port %d", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("could not listen port %d: %v", *port, err)
	}

	s := grpc.NewServer()
	say.RegisterTextToSpeechServer(s, server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("could not serve listener: %v", err)
	}
}

type server struct {
}

func (server) Say(ctx context.Context, text *say.Text) (*say.Speech, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("could not create temp file: %v", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close file %s: %v", f.Name(), err)
	}

	cmd := exec.Command("flite", "-t", text.Text, "-o", f.Name())
	if data, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("flite failed: %s", data)
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("could not read temp file %s: %v", f.Name(), err)
	}

	return &say.Speech{Audio: data}, nil
}
