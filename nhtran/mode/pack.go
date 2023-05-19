package mode

import (
	"gohtran/nhtran/params"
	"log"
	"strconv"
)

type pack struct {
	header    []byte
	body      []byte
	temp      []byte //
	plaintext bool

	complete bool

	pack interface{}
}

func NewPack(data []byte) *pack {
	p := &pack{}
	if len(data) < params.HeaderLen {
		p.temp = data
		return p
	}

	if params.HeaderRegexp.Match(data[:params.HeaderLen]) {
		p.header = data[:params.HeaderLen]
		if len(data) == params.HeaderLen {
			return p
		}
		p.body = data[params.HeaderLen:]
	} else {
		// Plaintext
		p.header = params.Header
		p.body = data
	}

	p.parse()
	return p
}

func (p *pack) push(data []byte) {
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
	if deletionLength < len(data) {
		p.pack = NewPack(data[deletionLength:])
	}
}

func (p *pack) parse() {
	singlePackLength := p.GetPackLen()
	switch {
	case singlePackLength == len(p.body):
		p.complete = true
		return
	//case singlePackLength > len(p.body):
	case singlePackLength < len(p.body):
		completeBody := p.body[:singlePackLength]
		p.body = completeBody
		p.complete = true

		data := p.body[singlePackLength:]
		nextPack := NewPack(data)
		p.pack = nextPack
	}
}

func (p pack) GetPackLen() int {
	length := string(p.header[params.PackStartLocal:params.PackEndLocal])
	l, err := strconv.Atoi(length)
	if err != nil {
		log.Println("decode header error")
	}
	return l
}

func (p *pack) getDeep() *pack {
	if p.pack == nil {
		return nil
	}

	return p.pack.(*pack)
}