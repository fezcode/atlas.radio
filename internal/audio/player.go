package audio

import (
	"context"
	"net/http"
	"sync"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type Player struct {
	otoCtx *oto.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(streamURL string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Stop previous
	if p.cancel != nil {
		p.cancel()
	}

	// 1. Fetch Stream
	resp, err := http.Get(streamURL)
	if err != nil {
		return err
	}

	// 2. Decode MP3
	decoder, err := mp3.NewDecoder(resp.Body)
	if err != nil {
		resp.Body.Close()
		return err
	}

	// 3. Init Oto if not initialized
	if p.otoCtx == nil {
		op := &oto.NewContextOptions{
			SampleRate:   decoder.SampleRate(),
			ChannelCount: 2,
			Format:       oto.FormatSignedInt16LE,
		}
		ctx, ready, err := oto.NewContext(op)
		if err != nil {
			resp.Body.Close()
			return err
		}
		<-ready
		p.otoCtx = ctx
	}

	// 4. Create Player
	player := p.otoCtx.NewPlayer(decoder)
	player.Play()

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	// 5. Monitor and Cleanup
	go func() {
		defer resp.Body.Close()
		defer player.Close()

		<-ctx.Done() // Block until stopped
	}()

	return nil
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}
}
