package main

import (
	"gohtran/control"
	"gohtran/mode"
	"gohtran/params"
	"log"
	"os"
	"strconv"
)

func main() {
	args := os.Args
	argc := len(args) - 1

	var core = control.NewCore()
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case params.Listen, params.Tran, params.Slave:
			core.Net.SetDesign(args[i])
			if argc < i+2 || params.ExistParams(args[i+1]) || params.ExistParams(args[i+2]) {
				log.Fatalln("params is error")
			}
			core.Net.FirstParam = args[i+1]
			core.Net.SecondParam = args[i+2]
			i += 2
			break
		case params.Left, params.Right:
			core.Net.Crypt.Open()
			core.Net.Crypt.SideParams = args[i]
			if argc < i+1 || params.ExistParams(args[i+1]) {
				log.Fatalln("params is error")
			}
			op, err := strconv.Atoi(args[i+1])
			if err != nil {
				log.Fatalln(err)
			}
			core.Net.Crypt.OperationParams = uint(op)
			if argc == i+2 && !params.ExistParams(args[i+2]) {
				core.Net.Crypt.SetAesKey(args[i+2])
				i++
			}
			i++
		case params.Slice:
			mode.Slice()
			break
		case params.Log:
			if argc <= i+1 && params.ExistParams(args[i+1]) {
				log.Fatalln("you need to specify the file path")
			}
			mode.Log(args[i+1])
			i++
			break
		case params.Help:
			mode.Help()
			os.Exit(0)
		default:
			log.Fatalf("%s doesn't exist", args[i])
		}

	}

	core.Run()
}
