package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func echoHandler(conn net.Conn) {
	// generate id based on epoch
	id := time.Now().Unix()
	log.Printf("Handling new connection, id: %d\n", id)

	ticker := time.NewTicker(10 * time.Second)

	defer func() {
		log.Printf("Closing connection %d\n", id)
		conn.Close()
		ticker.Stop()
	}()
	timeoutDuration := 5 * time.Second
	buffReader := bufio.NewReader(conn)

	for {
		select {
		case t := <-ticker.C:
			log.Printf("tick at %s\n", t)
			// buffWriter.WriteString("PING")
			conn.Write([]byte(fmt.Sprintln("PING :tmi.twitch.tv")))
			// wait for PONG
			conn.SetReadDeadline(time.Now().Add(timeoutDuration))
			bytes, err := buffReader.ReadBytes('\n')
			if err != nil {
				log.Println("Client closed connection")
				return
			}
			line := strings.TrimSpace(string(bytes))
			if line != "PONG :tmi.twitch.tv" {
				log.Println("Missing PONG response")
				return
			}

		default:
			// set deadline for reading.
			conn.SetReadDeadline(time.Now().Add(timeoutDuration))

			// read data
			bytes, err := buffReader.ReadBytes('\n')
			if err != nil {
				switch {
				case errors.Is(err, io.EOF):
					log.Println("Client closed connection")
					return
				case errors.Is(err, os.ErrDeadlineExceeded):
					log.Print("read timout continue")
					continue
				default:
					log.Printf("ERROR: (%T) %s\n", err, err)
					continue
				}

			}
			log.Printf(">> %s", bytes)
			conn.Write([]byte(fmt.Sprintf("echo: %s", bytes)))
		}
	}
}

func main() {
	var address, port string
	flag.StringVar(&address, "address", "127.0.0.1", "Echo Server address")
	flag.StringVar(&port, "port", "8888", "Echo Server port")
	flag.Parse()

	log.Println("Starting Echo server")
	// init
	serverAddr := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Panicln(err)
	}
	defer listener.Close()

	for {
		// listen for an connection
		conn, err := listener.Accept()
		if err != nil {
			log.Panicln(err)
		}
		// create go routine for connection
		go echoHandler(conn)
	}
}
