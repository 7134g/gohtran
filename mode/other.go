package mode

import (
	"fmt"
	"log"
	"os"
)

func Slice() {
	f, err := os.OpenFile(os.DevNull, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(f)
}

func Log(logFileName string) {
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(f)
}

func Help() {
	fmt.Println("+-----------------------------help information--------------------------------+")
	fmt.Println(`usage: "-listen port1 port2" #example: "gohtran -listen 8888 3389" `)
	fmt.Println(`       "-tran port1 ip:port2" #example: "gohtran -tran 8888 1.1.1.1:3389" `)
	fmt.Println(`       "-slave ip1:port1 ip2:port2" #example: "gohtran -slave 127.0.0.1:3389 1.1.1.1:8888" `)
	fmt.Println(`       "-e enable gzip and aes functionality`)
	fmt.Println(`       "-aes enable aes functionality, parameters is key, defaults to 16 bits`)
	fmt.Println(`       "-gzip enable gzip functionality`)
	fmt.Println(`       "-h program documentation`)
	fmt.Println(`       "-s silent mode,no information is displayed`)
	fmt.Println(`       "-log output transferred data to file`)
	fmt.Println(`============================================================`)
	fmt.Println("If you see start transmit, that means the data channel is established")
}
