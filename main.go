package main

import (
	"gohtran/control"
	"gohtran/mode"
	"gohtran/params"
	"log"
	"os"
)

func main() {
	args := os.Args
	argc := len(args) - 1

	var core = control.NewCore()
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case params.Listen, params.Tran, params.Slave:
			core.Net.SetDesign(args[i])
			if argc < i+2 {
				log.Fatalln("params is error")
			}
			core.Net.FirstParam = args[i+1]
			core.Net.SecondParam = args[i+2]
			i += 2
			break
		case params.Aes, params.AesGzip:
			core.Net.Crypt.SetScript(args[i])
			if argc > i+1 && !params.ExistParams(args[i+1]) {
				core.Net.Crypt.AesKey = []byte(args[i+1])
			} else {
				core.Net.Crypt.AesKey = params.AesDefaultKey
			}
			break
		case params.Gzip:
			core.Net.Crypt.SetScript(args[i])
			break
		case params.Slice:
			mode.Slice()
			break
		case params.Log:
			if argc <= i+1 && params.ExistParams(args[i+1]) {
				log.Fatalln("you need to specify the file path")
			}
			mode.Log(args[i+1])
			break
		case params.Help:
			mode.Help()
			os.Exit(0)

		}

	}

	core.Run()
}
