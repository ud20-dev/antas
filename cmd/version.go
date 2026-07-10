package cmd

import (
	"fmt"
	"github.com/ud20-dev/antas/internal/data"
)

func PrintVersion(){
	fmt.Println(data.AppName, data.GetVersion())
}