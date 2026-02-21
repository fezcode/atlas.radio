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

	// Check if ffplay exists
	_, err := exec.LookPath("ffplay")
	if err != nil {
		return fmt.Errorf("ffplay not found in PATH. Please install FFmpeg.")
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel

	// -nodisp: no video window
	// -autoexit: exit when stream ends
	p.cmd = exec.CommandContext(ctx, "ffplay", "-nodisp", "-autoexit", streamURL)
	
	err = p.cmd.Start()
	if err != nil {
		cancel()
		return err
	}

	return nil
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
