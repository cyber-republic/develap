package main

import (
	"fmt"
	"github.com/cyber-republic/develap/cmd"
	"github.com/spf13/viper"
)

var (
	// VERSION is set during build
	VERSION string
)

func main() {
	cmd.Execute(VERSION)
	fmt.Println(viper.GetString("key"))
}
