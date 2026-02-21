package audio

import (
	"bufio"
	"context"
	"fmt"
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
	p.Stop()
	p.mu.Lock()
	defer p.mu.Unlock()

	// 1. Start Request
	resp, err := http.Get(streamURL)
	if err != nil {
		return err
	}

	// 2. Use a buffered reader to prevent small read issues with go-mp3
	bufReader := bufio.NewReaderSize(resp.Body, 32*1024)

	// 3. Decode MP3
	decoder, err := mp3.NewDecoder(bufReader)
	if err != nil {
		resp.Body.Close()
		return fmt.Errorf("mp3 decode failed: %w", err)
	}

	// 4. Init Oto if not already
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

	// 5. Create and Start Player
	player := p.otoCtx.NewPlayer(decoder)
	player.Play()

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	// 6. Cleanup Goroutine
	go func() {
		defer resp.Body.Close()
		defer player.Close()
		<-ctx.Done()
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
