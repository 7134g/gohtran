package mode

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"gohtran/params"
	"io/ioutil"
	"log"
)

const (
	_ uint = iota
	encrypt
	decrypt
	compression
	decompression
)

type CryptMode struct {
	cryptFunc func([]byte) ([]byte, error)

	script string
	AesKey []byte
}

func (c *CryptMode) SetScript(s string) {
	c.script = s
	//switch s {
	//case params.Aes:
	//	c.cryptFunc = c.Aes
	//case params.Gzip:
	//	c.cryptFunc = c.Gzip
	//case params.AesGzip:
	//	c.cryptFunc = c.AesGzip
	//}
}

func (c *CryptMode) GetScript() string {
	return c.script
}

func (c *CryptMode) Aes(p *pack) error {
	switch p.header[params.AesLocal] {
	case params.AesEncryptSingle:
		p.header[params.AesLocal] = params.AesDecryptSingle
		p.body = AESEncrypt(p.body, c.AesKey)
	case params.AesDecryptSingle:
		p.header[params.AesLocal] = params.AesEncryptSingle
		p.body = AESDecrypt(p.body, c.AesKey)
	case params.Plaintext:
		p.header[params.AesLocal] = params.AesDecryptSingle
		p.body = AESEncrypt(p.body, c.AesKey)
	default:
		return errors.New("crypt error")
	}
	return nil
}

func (c *CryptMode) Gzip(p *pack) error {

	switch p.header[params.GzipLocal] {
	case params.GzipEncryptSingle:
		p.header[params.GzipLocal] = params.GzipDecryptSingle
		p.body = Compression(p.body)
	case params.GzipDecryptSingle:
		p.header[params.GzipLocal] = params.GzipEncryptSingle
		p.body = Decompress(p.body)
	case params.Plaintext:
		p.header[params.GzipLocal] = params.GzipDecryptSingle
		p.body = Compression(p.body)
	default:
		return errors.New("gzip error")
	}

	return nil
}

func (c *CryptMode) AesGzip(p *pack) error {
	switch p.header[params.AesLocal] {
	case params.Plaintext, params.AesEncryptSingle:
		err := c.Aes(p)
		if err != nil {
			return err
		}

		return c.Gzip(p)
	case params.AesDecryptSingle:
		err := c.Gzip(p)
		if err != nil {
			return err
		}

		return c.Aes(p)
	}
	return nil

}

func (c *CryptMode) IsOpen() bool {
	return !(len(c.script) == 0)
}

func AESEncrypt(origData, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	origData = PKCS7Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])

	crypt := make([]byte, len(origData))

	blockMode.CryptBlocks(crypt, origData)

	return crypt
}

func AESDecrypt(crypt, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypt))
	blockMode.CryptBlocks(origData, crypt)
	origData = PKCS7UnPadding(origData)
	return origData
}

func PKCS7Padding(origData []byte, blockSize int) []byte {
	padding := blockSize - len(origData)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(origData, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:length-unPadding]
}

func Compression(oriData []byte) []byte {
	var buf bytes.Buffer
	write := gzip.NewWriter(&buf)
	_, err := write.Write(oriData)
	if err != nil {
		log.Fatalln(err)
	}
	_ = write.Flush()
	return buf.Bytes()
}

func Decompress(oriData []byte) []byte {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(r)
		}
	}()

	var buf bytes.Buffer
	buf.Write(oriData)
	read, err := gzip.NewReader(&buf)
	if err != nil {
		log.Fatalln(err)
	}
	pData, _ := ioutil.ReadAll(read)
	return pData
}
