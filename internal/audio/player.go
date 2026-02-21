package audio

import (
	"io"
	"net/http"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type Player struct {
	closer io.Closer
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(streamURL string) error {
	p.Stop()

	resp, err := http.Get(streamURL)
	if err != nil {
		return err
	}
	p.closer = resp.Body

	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		p.Stop()
		return err
	}

	// Initialize speaker only once or if format changes
	// For simplicity in this novice version, we init with every new stream
	// speaker.Init is safe to call multiple times if we close previous
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		p.Stop()
	})))

	return nil
}

func (p *Player) Stop() {
	speaker.Clear()
	if p.closer != nil {
		p.closer.Close()
		p.closer = nil
	}
}
