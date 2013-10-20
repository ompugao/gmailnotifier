package main

import (
	"fmt"
	"github.com/ziutek/glib"
	"github.com/ziutek/gst"
	"os"
	//"time"
)

type Player struct {
	pipe      *gst.Element
	bus       *gst.Bus
	file_path string
	loop      *glib.MainLoop
}

func (p *Player) onMessage(bus *gst.Bus, msg *gst.Message) {
	fmt.Println(msg.GetType())
	switch msg.GetType() {
	case gst.MESSAGE_EOS:
		p.pipe.SetState(gst.STATE_NULL)
	case gst.MESSAGE_ERROR:
		p.pipe.SetState(gst.STATE_NULL)
		err, debug := msg.ParseError()
		fmt.Printf("Error: %s (debug: %s)\n", err, debug)
	}
}

func (p *Player) onEndofStream(bus *gst.Bus, msg *gst.Message) {
	p.pipe.SetState(gst.STATE_NULL)
	p.loop.Quit()
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 2 {
		fmt.Println("usage: ", os.Args[0], " audiofile")
		os.Exit(1)
	}

	p := new(Player)
	p.file_path = os.Args[1]
	if _, err := os.Stat(p.file_path); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", p.file_path)
		return
	}

	p.pipe = gst.ElementFactoryMake("playbin2", "autoplay")
	p.bus = p.pipe.GetBus()
	p.bus.AddSignalWatch()
	//p.bus.Connect("message", (*Player).onMessage, p)
	p.bus.Connect("message::eos", (*Player).onEndofStream, p)
	p.bus.EnableSyncMessageEmission()
	//p.bus.Connect("sync-message::element", (*Player).onSyncMessage, p)

	p.pipe.SetProperty("uri", "file://"+p.file_path)
	p.pipe.SetState(gst.STATE_PLAYING)
	//for {
	//state, _, _ := p.pipe.GetState(gst.CLOCK_TIME_NONE)
	//fmt.Println(state)
	//if state == gst.STATE_NULL {
	//break
	//}
	//time.Sleep(time.Second * 1)
	//}
	p.loop = glib.NewMainLoop(nil)
	p.loop.Run()
}
