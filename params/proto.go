package params

import "regexp"

var (
	HeaderLen       = 6
	Header          = []byte("000000")
	HeaderRegexp, _ = regexp.Compile(`^[$%][#&][0-9]{4}`)

	AesLocal       = 0
	GzipLocal      = 1
	PackStartLocal = 2
	PackEndLocal   = 6

	Plaintext         byte = '0'
	AesEncryptSingle  byte = '$' // aes encrypt
	AesDecryptSingle  byte = '%' // aes decrypt
	GzipEncryptSingle byte = '#' // gzip compression
	GzipDecryptSingle byte = '&' // gzip decompression

	AesDefaultKey = []byte("}3#*%-&*{+>?@8;'")
)
