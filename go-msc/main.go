package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	mss "github.com/whyrusleeping/go-multistream"
	"io"
	"net"
	"strings"
)

func EchoHandler(rwc io.ReadWriteCloser) error {
	defer rwc.Close()
	_, err := io.Copy(rwc, rwc)
	return err
}

func Client(target string, proto string) error {
	con, err := net.Dial("tcp", target)
	if err != nil {
		return err
	}
	defer con.Close()

	err = mss.SelectProtoOrFail(proto, con)
	if err != nil {
		return err
	}

	data := make([]byte, 4096)
	rand.Read(data)
	errs := make(chan error)

	go func() {
		_, err := con.Write(data)
		errs <- err
	}()

	go func() {
		resp := make([]byte, 4096)
		_, err := io.ReadFull(con, resp)
		if err != nil {
			errs <- err
			return
		}

		if !bytes.Equal(resp, data) {
			errs <- fmt.Errorf("data wasnt the same!!")
		}
		errs <- nil
	}()

	err = <-errs
	if err != nil {
		return err
	}
	err = <-errs
	if err != nil {
		return err
	}
	return nil
}

func Server(listen string, protos []string) error {
	list, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	mux := mss.NewMultistreamMuxer()
	for _, p := range protos {
		mux.AddHandler(p, EchoHandler)
	}

	for {
		con, err := list.Accept()
		if err != nil {
			return err
		}

		err = mux.Handle(con)
		if err != nil {
			return err
		}
	}
}

func main() {
	client := flag.Bool("client", false, "specify to run as client")
	addr := flag.String("addr", ":5050", "address to dial/listen on")
	proto := flag.String("protos", "/test", "which protocols to test")

	flag.Parse()

	protos := strings.Split(*proto, ",")
	if *client {
		for _, p := range protos {
			fmt.Printf("testing proto: %s . . . ", p)
			err := Client(*addr, p)
			if err != nil {
				fmt.Printf("client failed on protocol: %s\nerror: %s\n", p, err)
				return
			}

			fmt.Println("pass!")
		}
	} else {
		err := Server(*addr, protos)
		if err != nil {
			fmt.Printf("server failed: %s\n", err)
			return
		}
	}
}
