package control

import (
	"gohtran/mode"
	"gohtran/params"
	"log"
	"time"
)

type Core struct {
	stop chan struct{}

	Net mode.NetMode
}

func NewCore() *Core {
	c := &Core{stop: make(chan struct{})}
	c.Net.Crypt.SetAesKey(params.AesDefaultKey)
	return c
}

func (c *Core) Run() {
	//go c.connectTimeout()

	switch c.Net.GetDesign() {
	case params.Listen:
		if err := c.Net.Listen(); err != nil {
			log.Fatalln(err)
		}
	case params.Tran:
		if err := c.Net.Tran(); err != nil {
			log.Fatalln(err)
		}
	case params.Slave:
		if err := c.Net.Slave(); err != nil {
			log.Fatalln(err)
		}
	}
}

func (c *Core) connectTimeout() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Fatalln("Unable to link")
		case <-c.stop:
			return
		}
	}
}
