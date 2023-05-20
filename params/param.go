package params

const (
	Listen = "-listen"
	Tran   = "-tran"
	Slave  = "-slave"

	Left         = "-left"
	Right        = "-right"
	Aes     uint = 1
	Gzip    uint = 2
	AesGzip uint = 3

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
		case Aes, Gzip, AesGzip:
			return true
		default:
			return false
		}
	default:
		return false
	}

}
