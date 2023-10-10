package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct{}

func (fs *FileServer) Start() {
	ln, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal(err)
	}

	for {
		con, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go fs.read(con)
	}
}

func (fs *FileServer) read(conn net.Conn) {
	buffer := new(bytes.Buffer)

	for {
		var size int64

		binary.Read(conn, binary.LittleEndian, &size)

		n, err := io.CopyN(buffer, conn, size)

		if err != nil {
			fmt.Printf("Cannot read due %s", err.Error())
		}

		fmt.Print(buffer.Bytes())
		fmt.Printf("Received %d bytes over the network\n", n)
	}
}

func sendFile(bufferSize int) error {
	file := make([]byte, bufferSize)

	_, err := io.ReadFull(rand.Reader, file)

	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3000")

	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(bufferSize))

	n, err := io.CopyN(conn, bytes.NewReader(file), int64(bufferSize))

	fmt.Printf("Written %d bytes over the network", n)

	return nil
}

func main() {
	go func() {
		time.Sleep(4 * time.Second)
		sendFile(1e+6)
	}()

	server := &FileServer{}
	server.Start()
}
