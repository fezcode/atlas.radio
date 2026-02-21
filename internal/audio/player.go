package audio

import (
	"os/exec"
	"runtime"
)

type Player struct {
	cmd *exec.Cmd
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(streamURL string) error {
	p.Stop()

	// Use ffplay (part of ffmpeg) which is common, or mpv
	// -nodisp: don't show video window
	// -autoexit: exit when stream ends
	binary := "ffplay"
	args := []string{"-nodisp", "-autoexit", streamURL}

	if runtime.GOOS == "windows" {
		// Just in case we need to handle windows paths
	}

	p.cmd = exec.Command(binary, args...)
	return p.cmd.Start()
}

func (p *Player) Stop() {
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd.Wait()
	}
}
