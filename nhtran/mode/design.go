package mode

import (
	"fmt"
	"gohtran/nhtran/params"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type NetMode struct {
	Crypt CryptMode

	netFunc func() error

	design string

	FirstParam  string
	SecondParam string
}

func (n *NetMode) SetDesign(s string) {
	n.design = s
}

func (n *NetMode) GetDesign() string {
	return n.design
}

func (n *NetMode) Listen() error {
	if err := checkPort(n.FirstParam); err != nil {
		return nil
	}
	if err := checkPort(n.SecondParam); err != nil {
		return nil
	}

	serve1, err := createServer(fmt.Sprintf("127.0.0.1:%s", n.FirstParam))
	if err != nil {
		return err
	}
	serve2, err := createServer(fmt.Sprintf("127.0.0.1:%s", n.SecondParam))
	if err != nil {
		return err
	}

	for {
		conn1, err := accept(serve1)
		if err != nil {
			return err
		}
		conn2, err := accept(serve2)
		if err != nil {
			return err
		}
		if conn1 == nil || conn2 == nil {
			continue
		}
		n.forward(conn1, conn2)
	}

}

func (n *NetMode) Tran() error {
	if err := checkPort(n.FirstParam); err != nil {
		return nil
	}
	if err := checkAddress(n.SecondParam); err != nil {
		return nil
	}

	server, err := createServer(fmt.Sprintf("127.0.0.1:%s", n.FirstParam))
	if err != nil {
		return err
	}
	for {
		conn, err := accept(server)
		if err != nil {
			return err
		}

		for {
			target, err := createDial(n.SecondParam)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			n.forward(conn, target)
			break
		}

	}

}

func (n *NetMode) Slave() error {
	if err := checkAddress(n.FirstParam); err != nil {
		return nil
	}
	if err := checkAddress(n.SecondParam); err != nil {
		return nil
	}

	for {
		var target1, target2 net.Conn
		var err error
		for {
			target1, err = createDial(n.FirstParam)
			if err == nil {
				break
			} else {
				time.Sleep(time.Second)
			}
		}
		for {
			target2, err = createDial(n.SecondParam)
			if err == nil {
				break
			} else {
				time.Sleep(time.Second)
			}
		}
		n.forward(target1, target2)

	}

}

func (n *NetMode) forward(conn1 net.Conn, conn2 net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go n.connCopy(conn1, conn2, &wg)
	go n.connCopy(conn2, conn1, &wg)
	wg.Wait()
}

func (n *NetMode) connCopy(wConn net.Conn, rConn net.Conn, wg *sync.WaitGroup) {
	defer func() {
		_ = wConn.Close()
	}()

	if n.Crypt.IsOpen() {
		n.connReadWrite(wConn, rConn)
	} else {
		// plaintext
		_, _ = io.Copy(wConn, rConn)
	}

	wg.Done()
}

func (n *NetMode) connReadWrite(wConn net.Conn, rConn net.Conn) {
	var netPack *pack

	chunk := make([]byte, 1024)

	for {
		rCont, err := rConn.Read(chunk)
		if err != nil {
			break
		}

		switch {
		case err != nil:
			return
		case rCont > 0:
			switch {
			case netPack == nil:
				netPack = NewPack(chunk)
				netPack = n.sendData(netPack, wConn)
			case netPack != nil:
				netPack.push(chunk)
				netPack = n.sendData(netPack, wConn)
			}
			if netPack.complete {
				netPack = nil
			}

		}

	}
}

func (n *NetMode) sendData(pack *pack, wConn net.Conn) *pack {
	for {
		if !pack.complete {
			return pack
		}

		if len(pack.body) <= 0 {
			return pack
		}
		_, _ = wConn.Write(n.buildPackage(pack))

		if pack.getDeep() == nil {
			return pack
		}
		pack = pack.getDeep()

	}
}

func (n *NetMode) buildPackage(p *pack) []byte {
	var err error
	switch n.GetDesign() {
	case params.Aes:
		err = n.Crypt.Aes(p)
	case params.Gzip:
		err = n.Crypt.Gzip(p)
	case params.AesGzip:
		err = n.Crypt.AesGzip(p)
	default:
		log.Println("cannot find crypt type")
		return nil
	}

	if err != nil {
		log.Println("crypt error")
		return nil
	}

	return append(p.header, p.body...)
}
