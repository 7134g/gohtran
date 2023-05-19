package params

const (
	Listen = "-listen"
	Tran   = "-tran"
	Slave  = "-slave"

	Left            = "-left"
	Right           = "-right"
	NewAes     uint = 1
	NewGzip    uint = 2
	NewAesGzip uint = 3

	//Aes     = "-aes"
	//Gzip    = "-gzip"
	//AesGzip = "-e"

	Slice = "-s"
	Help  = "-h"
	Log   = "-log"
)

func ExistParams(p interface{}) bool {
	switch p.(type) {
	case string:
		switch p {
		case Listen, Tran, Slave, Slice, Help, Log:
			return true
		default:
			return false
		}
	case uint:
		switch p {
		case NewAes, NewGzip, NewAesGzip:
			return true
		default:
			return false
		}
	default:
		return false
	}

}
