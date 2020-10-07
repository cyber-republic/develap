package system

import (
	"fmt"
	"net"
	"os"
)

func GetRandomPort() string {
	listener, err := net.Listen("tcp", ":0")
	defer listener.Close()
	if err != nil {
		panic(err)
	}
	portInt := listener.Addr().(*net.TCPAddr).Port
	portString := fmt.Sprintf("%d", portInt)
	return portString
}

func GetCurrentDir() string {
	var currentDir string
	if pwd, err := os.Getwd(); err == nil {
		currentDir = pwd
	}
	return currentDir
}