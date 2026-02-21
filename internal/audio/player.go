package audio

import (
	"context"
	"os/exec"
	"runtime"
	"time"
)

type Player struct {
	cancel context.CancelFunc
	cmd    *exec.Cmd
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(streamURL string) error {
	// 1. Force kill everything first
	p.Stop()
	
	// 2. Short pause to ensure OS releases the audio device
	time.Sleep(200 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	// 3. Start new process
	p.cmd = exec.CommandContext(ctx, "ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", streamURL)
	
	return p.cmd.Start()
}

func (p *Player) Stop() {
	// Signal our specific context if it exists
	if p.cancel != nil {
		p.cancel()
	}

	// NUCLEAR OPTION: Kill all ffplay processes by name.
	// This is the only way to be 100% sure on Windows when dealing with media orphans.
	if runtime.GOOS == "windows" {
		_ = exec.Command("taskkill", "/F", "/IM", "ffplay.exe", "/T").Run()
	} else {
		_ = exec.Command("pkill", "-9", "ffplay").Run()
	}

	// Cleanup our command handle
	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.cmd.Wait()
		p.cmd = nil
	}
}
