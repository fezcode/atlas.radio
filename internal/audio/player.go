package audio

import (
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

type Player struct {
	mu      sync.Mutex
	closer  io.Closer
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(streamURL string) error {
	p.Stop()
	p.mu.Lock()
	defer p.mu.Unlock()

	resp, err := http.Get(streamURL)
	if err != nil {
		return err
	}
	p.closer = resp.Body

	// Decode the stream
	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		resp.Body.Close()
		return err
	}

	// speaker.Init is safe to call multiple times in beep/v2
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		streamer.Close()
		resp.Body.Close()
		return err
	}

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		p.Stop()
	})))

	return nil
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	speaker.Clear()
	if p.closer != nil {
		p.closer.Close()
		p.closer = nil
	}
}
