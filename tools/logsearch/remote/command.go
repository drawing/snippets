package remote

import (
	"bytes"
	"errors"

	"github.com/kr/pty"
	"log"
	"os"
	"os/exec"
	"unicode"
)

type RemoteCommand struct {
	shell []byte
	io    *os.File
	addr  string
	pass  string
}

func (remote *RemoteCommand) Close() error {
	remote.io.Close()
	return nil
}

func NewRemote(addr string, pass string) (*RemoteCommand, error) {
	log.Println("remote:", addr)
	c := exec.Command("ssh", "-p36000", addr)
	f, err := pty.Start(c)
	if err != nil {
		log.Println("start:", err)
		return nil, err
	}

	var remote RemoteCommand
	remote.io = f
	remote.addr = addr
	remote.pass = pass

	err = remote.connect()
	if err != nil {
		remote.io.Close()
		log.Println("connect:", err)
		return nil, err
	}

	return &remote, err
}

func (remote *RemoteCommand) connect() error {
	var cache []byte

	echo := make([]byte, 1024)

	echo_pass := false
	echo_yes := false
	finish := false

	for {
		n, err := remote.io.Read(echo)
		if err != nil {
			return err
		}

		for i, v := range echo[:n] {
			if v == byte('\r') {
				echo[i] = ' '
			} else if v == byte('$') || v == byte('#') || v == byte('>') {
				finish = true
			}
		}

		cache = append(cache, echo[0:n]...)

		if finish {
			break
		}

		if !echo_yes && bytes.Contains(cache, []byte("(yes/no)?")) {
			_, err := remote.io.Write([]byte("yes\n"))
			if err != nil {
				return err
			}
			echo_yes = true
		} else if !echo_pass && bytes.Contains(cache, []byte("password")) {
			_, err := remote.io.Write([]byte(remote.pass + "\n"))
			if err != nil {
				return err
			}
			echo_pass = true
		}
	}

	begin := 0
	end := 0
	for i := len(cache) - 1; i >= 0; i-- {
		v := cache[i]
		if v == byte('$') || v == byte('#') || v == byte('>') {
			end = i
		} else if !unicode.IsGraphic(rune(v)) {
			begin = i + 1
			break
		}
	}
	if begin >= end {
		return errors.New("find shell failed")
	}

	remote.shell = cache[begin : end+1]
	log.Println("shell:", string(remote.shell))

	return nil
}

func (remote *RemoteCommand) RunCommand(command string) (string, error) {
	log.Println("RunCommand:", command)

	var cache []byte

	_, err := remote.io.Write([]byte(command + "\n"))
	if err != nil {
		return "", err
	}

	echo := make([]byte, 1024)

	for {
		n, err := remote.io.Read(echo)
		if err != nil {
			return "", err
		}

		for i, v := range echo[:n] {
			if v == byte('\r') {
				echo[i] = ' '
			}
		}

		cache = append(cache, echo[0:n]...)

		if bytes.Contains(cache, remote.shell) {
			break
		}
	}

	offset := bytes.Index(cache, []byte{'\n'})
	if offset != -1 {
		cache = cache[offset+1:]
	} else {
		cache = []byte{}
	}

	offset = bytes.LastIndex(cache, []byte{'\n'})
	if offset != -1 {
		cache = cache[:offset]
	} else {
		cache = []byte{}
	}

	return string(cache), nil
}

func (remote *RemoteCommand) ExecTelnet(host string, command string) (string, error) {
	log.Println("Telnet", host, command)

	var cache []byte

	_, err := remote.io.Write([]byte(host + "\n"))
	if err != nil {
		return "", err
	}

	echo := make([]byte, 1024)
	exec := 0

	for {
		n, err := remote.io.Read(echo)
		if err != nil {
			return "", err
		}

		for i, v := range echo[:n] {
			if v == byte('\r') {
				echo[i] = ' '
			}
			if v == byte('\n') && exec > 0 {
				exec++
			}
		}

		cache = append(cache, echo[0:n]...)

		if exec > 2 && exec != 10000 {
			_, err := remote.io.Write([]byte("quit\n"))
			if err != nil {
				return "", err
			}
			exec = 10000
		}

		if exec == 0 && bytes.Contains(cache, []byte("Escape character")) {
			cache = []byte{}
			_, err := remote.io.Write([]byte(command + "\n"))
			if err != nil {
				return "", err
			}
			exec++
		}
		if bytes.Contains(cache, remote.shell) {
			break
		}
	}

	offset := bytes.LastIndex(cache, []byte{'\n'})
	if offset != -1 {
		cache = cache[:offset]
	} else {
		cache = []byte{}
	}

	return string(cache), nil
}
