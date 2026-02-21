package audio

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"
)

type Player struct {
	cmd *exec.Cmd
	mu  sync.Mutex
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(streamURL string) error {
	p.Stop()
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if mpv exists (it's the most stable for world radio)
	_, err := exec.LookPath("mpv")
	if err != nil {
		return fmt.Errorf("mpv not found in PATH. Please install mpv.")
	}

	// Use mpv with process group management
	// --no-video: obviously
	// --cache=yes: handle jitter
	// --terminal=no: hide its own output
	p.cmd = exec.Command("mpv", "--no-video", "--cache=yes", "--terminal=no", streamURL)
	
	return p.cmd.Start()
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd != nil && p.cmd.Process != nil {
		if runtime.GOOS == "windows" {
			// On Windows, taskkill /F /T is the most reliable way to kill the tree
			_ = exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", p.cmd.Process.Pid)).Run()
		} else {
			_ = p.cmd.Process.Kill()
		}
		_ = p.cmd.Wait()
		p.cmd = nil
	}
}
