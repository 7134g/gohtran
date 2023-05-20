package mode

import (
	"fmt"
	"gohtran/params"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type NetMode struct {
	wg sync.WaitGroup

	Crypt CryptMode

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
		return err
	}
	if err := checkPort(n.SecondParam); err != nil {
		return err
	}

	serve1, err := createServer(fmt.Sprintf("127.0.0.1:%s", n.FirstParam))
	if err != nil {
		return err
	}
	serve2, err := createServer(fmt.Sprintf("127.0.0.1:%s", n.SecondParam))
	if err != nil {
		return err
	}
	log.Printf("Link [127.0.0.1:%s] and [127.0.0.1:%s] are successfully established\n",
		n.FirstParam, n.SecondParam)
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

		log.Println("client ready")
		n.forward(conn1, conn2)
	}

}

func (n *NetMode) Tran() error {
	if err := checkPort(n.FirstParam); err != nil {
		return err
	}
	if err := checkAddress(n.SecondParam); err != nil {
		return err
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
				log.Println(err)
				time.Sleep(time.Second)
				continue
			}
			log.Println("client ready")
			n.forward(conn, target)
			break
		}

	}

}

func (n *NetMode) Slave() error {
	if err := checkAddress(n.FirstParam); err != nil {
		return err
	}
	if err := checkAddress(n.SecondParam); err != nil {
		return err
	}

	for {
		var target1, target2 net.Conn
		var err error
		for {
			target1, err = createDial(n.FirstParam)
			if err == nil {
				break
			} else {
				log.Println(err)
				time.Sleep(time.Second)
			}
		}
		for {
			target2, err = createDial(n.SecondParam)
			if err == nil {
				break
			} else {
				log.Println(err)
				time.Sleep(time.Second)
			}
		}
		log.Println("client ready")
		n.forward(target1, target2)

	}

}

func (n *NetMode) forward(leftConn net.Conn, rightConn net.Conn) {
	n.wg.Add(2)

	if n.Crypt.GetStata() {
		switch n.Crypt.SideParams {
		case params.Left:
			go n.cryptConn(rightConn, leftConn, true)
			go n.cryptConn(leftConn, rightConn, false)
			break
		case params.Right:
			go n.cryptConn(rightConn, leftConn, false)
			go n.cryptConn(leftConn, rightConn, true)
			break
		}
	} else {
		go n.connCopy(leftConn, rightConn)
		go n.connCopy(rightConn, leftConn)
	}

	n.wg.Wait()
}

func (n *NetMode) connCopy(c1 net.Conn, c2 net.Conn) {
	defer func() {
		_ = c1.Close()
		n.wg.Done()
	}()

	_, _ = io.Copy(c1, c2)

}

func (n *NetMode) cryptConn(wConn net.Conn, rConn net.Conn, plaintext bool) {
	defer func() {
		_ = wConn.Close()
		n.wg.Done()
	}()

	var netPack *packet
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
			case err != nil:
				return
			case rCont > 0:
				data := chunk[:rCont]
				switch {
				case netPack == nil:
					netPack = newPacket(data, plaintext)
					netPack = n.sendData(netPack, wConn)
				case netPack != nil:
					netPack.push(data)
					netPack = n.sendData(netPack, wConn)
					if netPack.complete {
						netPack = nil
					}
				}

			}
		}

	}
}

func (n *NetMode) sendData(pack *packet, wConn net.Conn) *packet {
	for {
		if !pack.complete {
			return pack
		}

		if len(pack.body) <= 0 {
			return pack
		}
		_, _ = wConn.Write(n.buildPackage(pack))

		pack = pack.getDeep()
		if pack == nil {
			return nil
		}

	}
}

func (n *NetMode) buildPackage(p *packet) []byte {
	var err error
	switch n.Crypt.OperationParams {
	case params.NewAes:
		err = n.Crypt.Aes(p)
	case params.NewGzip:
		err = n.Crypt.Gzip(p)
	case params.NewAesGzip:
		err = n.Crypt.AesGzip(p)
	default:
		log.Println("cannot find crypt type")
		return nil
	}

	if err != nil {
		log.Println("crypt error")
		return nil
	}

	if p.header == nil {
		// plaintext
		return p.body
	}
	length := fmt.Sprintf("%04d", len(p.body))
	for i, s := range length {
		b := byte(s)
		p.header[i+2] = b
	}

	return append(p.header, p.body...)

}
