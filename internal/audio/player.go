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
	p.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	// Use mpv: simple, robust, handles almost any stream
	// --no-video: audio only
	// --msg-level=all=no: silent output
	p.cmd = exec.CommandContext(ctx, "mpv", "--no-video", "--msg-level=all=no", streamURL)
	
	return p.cmd.Start()
}

func (p *Player) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
	if p.cmd != nil && p.cmd.Process != nil {
		if runtime.GOOS == "windows" {
			_ = exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", p.cmd.Process.Pid)).Run()
		} else {
			_ = p.cmd.Process.Kill()
		}
		_ = p.cmd.Wait()
		p.cmd = nil
	}
}
