package main

import (
	"github.com/darwinia-network/token/cmd"
	"log"
	"os"
)

func main()  {

	if err := cmd.Run(os.Args); err != nil{
		log.Fatalln(err)
	}
}


