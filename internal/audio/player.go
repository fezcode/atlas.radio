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

	// Check for vlc (exceptionally stable for world radio)
	binary, err := exec.LookPath("vlc")
	if err != nil {
		// Fallback to cvlc if available (headless vlc)
		binary, err = exec.LookPath("cvlc")
		if err != nil {
			return fmt.Errorf("vlc/cvlc not found in PATH. Please install VLC.")
		}
	}

	// Use vlc with no-interface mode
	// --intf dummy: no gui
	// --play-and-exit: self-explanatory
	p.cmd = exec.Command(binary, "--intf", "dummy", streamURL)
	
	return p.cmd.Start()
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

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
