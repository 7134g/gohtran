package mode

import (
	"fmt"
	"testing"
)

func TestPack(t *testing.T) {
	var (
		data1 = []byte("000007XXXXXXX")

		data2_1 = []byte("000020zzzzzzzzzz")
		data2_2 = []byte("zzzzzzzzzz")

		//data3_1 = []byte("000020zzzzzzzzzz")
		//data3_2 = []byte("zzzzzzzzzz000010qqqqqqqqqq")
		//
		//data4 = []byte("000013aaaaaaaaaaaaa000015aaaaaaaaaaaaabb")
		//
		//data5_1 = []byte("000013aaaaaaaaaaaaa000015aaaaaaaaaaaaabb000006ww")
		//data5_2 = []byte("wwww")
		//
		//data6_1 = []byte("000013aaaaaaaaaaaaa000015aaaaaaaaaaaaabb000006")
		//data6_2 = []byte("wwww")
	)

	t.Run("single", func(t *testing.T) {
		t.Log("need 7")
		p := newPacket(data1, false)
		t.Log(len(p.body))
		t.Log(string(p.body))
		t.Log(p)
	})

	t.Run("single_deletion", func(t *testing.T) {
		t.Log("need 20")
		p := newPacket(data2_1, false)
		t.Log(len(p.body))
		t.Log(string(p.body))
		t.Log(p)

		if !p.complete {
			p.push(data2_2)
			np := p.getDeep()
			if np != nil {
				t.Log("parse error")
			}
		}

		t.Log(len(p.body))
		t.Log(string(p.body))
		t.Log(p)
	})

}

func TestName(t *testing.T) {
	t.Log(fmt.Sprintf("%04d", 10))
	t.Log(fmt.Sprintf("%04d", 102))
	t.Log(fmt.Sprintf("%04d", 1024))
	t.Log(fmt.Sprintf("%04d", 10240))
}
