package cmd

import (
	"fmt"
	pflag "github.com/spf13/pflag"
	"github.com/ud20-dev/antas/internal/data"
)

func PrintHelp() {
	fmt.Printf("%s - Stack it up until it doesn't compile\n",data.AppName)
	fmt.Printf("Usage: %s [OPTIONS] <path/to/file.pdf>\n",data.AppName)
	fmt.Println("Options:")
	pflag.PrintDefaults()
}