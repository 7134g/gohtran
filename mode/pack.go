package mode

import (
	"gohtran/params"
	"log"
	"strconv"
)

type packet struct {
	header []byte
	body   []byte
	temp   []byte

	complete bool

	pack interface{}
}

func newPacket(data []byte, plaintext bool) *packet {
	p := &packet{}
	if plaintext {
		p.header = params.Header
		p.body = data
		p.complete = true
		return p
	}

	if len(data) < params.HeaderLen {
		p.temp = data
		return p
	}

	p.parse(data)
	return p
}

func (p *packet) parse(data []byte) {
	p.header = data[:params.HeaderLen]
	p.body = data[params.HeaderLen:]
	singlePackLength := p.GetPackLen()
	switch {
	case singlePackLength == len(p.body):
		p.complete = true
		return
	case singlePackLength > len(p.body):
		return
	case singlePackLength < len(p.body):
		completeBody := p.body[:singlePackLength]
		p.body = completeBody
		p.complete = true

		data := p.body[singlePackLength:]
		nextPack := newPacket(data, false)
		p.pack = nextPack
	}
}

func (p *packet) push(data []byte) {
	if p.temp != nil {
		switch {
		case len(data)+len(p.temp) < params.HeaderLen:
			p.temp = append(p.temp, data...)
			return
		case len(data)+len(p.temp) == params.HeaderLen:
			p.header = append(p.temp, data...)
			p.temp = nil
			return
		case len(data)+len(p.temp) > params.HeaderLen:
			halfHeader := data[:params.HeaderLen-len(p.temp)]
			p.header = append(p.temp, halfHeader...)
			p.temp = nil
		}
	}

	singlePackLength := p.GetPackLen()
	deletionLength := singlePackLength - len(p.body)
	p.body = append(p.body, data[:deletionLength]...)
	if len(p.body) == singlePackLength {
		p.complete = true
	}
	if deletionLength < len(data) {
		p.pack = newPacket(data[deletionLength:], false)
	}
}

func (p packet) GetPackLen() int {
	length := string(p.header[params.PackStartLocal:params.PackEndLocal])
	l, err := strconv.Atoi(length)
	if err != nil {
		log.Println("decode header error")
	}
	return l
}

func (p *packet) getDeep() *packet {
	if p.pack == nil {
		return nil
	}

	return p.pack.(*packet)
}
