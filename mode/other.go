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
	fmt.Println(`       "-left The options are 1, 2, and 3`)
	fmt.Println(`       "-right The options are 1, 2, and 3`)
	fmt.Println(`       "-h program documentation`)
	fmt.Println(`       "-s silent mode,no information is displayed`)
	fmt.Println(`       "-log output transferred data to file`)
	fmt.Println(`The value of 1 corresponds to aes encryption and decryption`)
	fmt.Println(`The value of 2 corresponds to gzip compression and decompression`)
	fmt.Println(`The value of 3 corresponds to the simultaneous use of aes and gzip`)
	fmt.Println(`============================================================`)
	fmt.Println("If you see start transmit, that means the data channel is established")
}
