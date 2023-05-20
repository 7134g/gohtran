package mode

import (
	"fmt"
	"testing"
)

func TestAESEncrypt(t *testing.T) {
	key := []byte("!@#$%^&*()_+reds")
	origData := []byte("123456")
	fmt.Println(len(origData), string(origData))

	//origData := GetCtx(conn)
	en := AESEncrypt(origData, key)
	fmt.Println(len(en), string(en))
	s := string(en)
	d := []byte(s)

	de := AESDecrypt(d, key)
	fmt.Println(len(de), string(de))

}

func TestCompression(t *testing.T) {
	src := []byte("111111111111111111111111111111111111111111111111111111")
	dst := Compression(src)
	t.Log(src)
	t.Log(dst)
}
