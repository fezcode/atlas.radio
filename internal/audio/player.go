package audio

import (
	"context"
	"os/exec"
)

type Player struct {
	cancel context.CancelFunc
	cmd    *exec.Cmd
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(streamURL string) error {
	p.Stop()

	// Use a context so the process is bound to it
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	// -nodisp: no video
	// -autoexit: exit when done
	// -loglevel quiet: don't spam stderr
	p.cmd = exec.CommandContext(ctx, "ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", streamURL)
	
	return p.cmd.Start()
}

func (p *Player) Stop() {
	if p.cancel != nil {
		p.cancel() // Kill the context and the process with it
	}
	if p.cmd != nil && p.cmd.Process != nil {
		// Force kill just in case
		_ = p.cmd.Process.Kill()
		_ = p.cmd.Wait()
	}
}
