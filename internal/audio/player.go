package audio

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

type Player struct {
	cancel context.CancelFunc
	cmd    *exec.Cmd
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(streamURL string) error {
	// Sync Stop before playing new
	p.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	p.cmd = exec.CommandContext(ctx, "ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", streamURL)
	
	err := p.cmd.Start()
	if err != nil {
		cancel()
		return err
	}
	
	// Start a goroutine to wait for the process to end naturally
	go func() {
		_ = p.cmd.Wait()
	}()

	return nil
}

func (p *Player) Stop() {
	if p.cancel != nil {
		p.cancel() // Signal context cancellation
	}
	if p.cmd != nil && p.cmd.Process != nil {
		if runtime.GOOS == "windows" {
			// Synchronously kill the process tree
			_ = exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", p.cmd.Process.Pid)).Run()
		} else {
			_ = p.cmd.Process.Kill()
		}
		// Wait for it to be fully gone before returning
		_ = p.cmd.Wait()
		p.cmd = nil
	}
}
