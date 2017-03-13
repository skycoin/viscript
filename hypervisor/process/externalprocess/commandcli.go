package externalprocess

import (
	"os"
	"os/exec"

	"io"

	"bufio"

	"github.com/kr/pty"
)

// var GOPATH string = os.Getenv("GOPATH")

func woof() {
	cmd := exec.Command(os.Args[1])

	currentCmd, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	inChannel := make(chan []byte, 1024)
	outChannel := make(chan []byte, 1024)

	// Wait for inChannel and print it
	go func() {
		for {
			if f := <-inChannel; f != nil {
				println(f)
			}
		}
	}()

	// Wait for the inChannel to respond
	go func() {
		for {
			var b []byte
			currentCmd.Read(b)
			inChannel <- b
		}
	}()

	// Wait for the user input
	go func() {
		for {
			// cycle here
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			outChannel <- []byte(text)
		}
	}()

	// Wait for the out Channel to send it to the command that we are running currently
	go func() {
		for {
			if f := <-outChannel; f != nil {
				currentCmd.Write([]byte(f))
			}
		}
	}()

	io.Copy(os.Stdout, currentCmd)
}
