package params

const (
	Listen = "-listen"
	Tran   = "-tran"
	Slave  = "-slave"

	Aes     = "-aes"
	Gzip    = "-gzip"
	AesGzip = "-e"

	Slice = "-s"
	Help  = "-h"
	Log   = "-log"
)

func ExistParams(p string) bool {
	switch p {
	case Listen, Tran, Slave, Aes, Gzip, AesGzip, Slice, Help, Log:
		return true
	default:
		return false
	}
}
