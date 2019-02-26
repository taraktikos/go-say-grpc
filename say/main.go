package main

import (
	"context"
	"flag"
	"fmt"
	"go-say/api"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	backend := flag.String("b", "localhost:8080", "address of backend")
	output := flag.String("o", "output.wav", "output")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("usage:\n\t%s test to speak", os.Args[0])
		os.Exit(1)
	}

	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can't connect to %s: %v", *backend, err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("can't close connection: %v", err)
		}
	}()

	client := say.NewTextToSpeechClient(conn)
	text := &say.Text{
		Text: flag.Arg(0),
	}
	resp, err := client.Say(context.Background(), text)
	if err != nil {
		log.Fatalf("couldn't say text %s: %v", text.Text, err)
	}

	if err := ioutil.WriteFile(*output, resp.Audio, 0666); err != nil {
		log.Fatalf("could not write to %s: %v", *output, err)
	}

}
