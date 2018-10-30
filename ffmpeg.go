package ffmpegwrap

import (
	"fmt"
	"strings"
	"os/exec"
	"bytes"
	"syscall"
	"os"
)

type FFMPEG struct {
	BinPath  string
	Flags    map[string]string
	Input    *Input
	Outputs  []*Output
	Command struct {
		Cmd        *exec.Cmd
		Stdout     bytes.Buffer
		Stderr     bytes.Buffer
		ExitStatus int
	}
	Running bool
}

func NewFFMPEG() *FFMPEG {
	ffmpeg := &FFMPEG{}
	ffmpeg.BinPath = "ffmpeg"
	ffmpeg.Flags = map[string]string{
		"-hide_banner": "",
		"-y": "",
		"-loglevel": "error",
	}

	return ffmpeg
}

func (f *FFMPEG) SetBinPath(path string) {
	f.BinPath = path
}

func (f *FFMPEG) SetLogLevel(loglevel string) {
	f.Flags["-loglevel"] = loglevel
}

func (f *FFMPEG) SetFlags(flags map[string]string) {
	for k, v := range flags {
		f.Flags[k] = v
	}
}

func (f *FFMPEG) SetInput(input *Input) {
	f.Input = input
}

func (f *FFMPEG) AddOutput(output *Output) {
	f.Outputs = append(f.Outputs, output)
}

func (f *FFMPEG) GetStdout() string {
	return f.Command.Stdout.String()
}

func (f *FFMPEG) GetStderr() string {
	return f.Command.Stderr.String()
}

func (f *FFMPEG) Compile() string {
	var outputsStrings []string
	for _, output := range f.Outputs {
		outputsStrings = append(outputsStrings, output.Compile())
	}

	var flagsStrings []string
	for k, v := range f.Flags {
		flagsStrings = append(flagsStrings, k, v)
	}

	return fmt.Sprintf("%s %s %s %s", f.BinPath, strings.Join(flagsStrings, " "),
		f.Input.Compile(), strings.Join(outputsStrings, " "))
}

func (f *FFMPEG) Run() error {
	cmdParts := strings.Fields(f.Compile())
	f.Command.Cmd = exec.Command(cmdParts[0], cmdParts[1:]...)

	f.Command.Cmd.Stdout = &f.Command.Stdout
	f.Command.Cmd.Stderr = &f.Command.Stderr

	err := f.Command.Cmd.Start()
	if err != nil {
		return err
	}

	f.Running = true

	go func(f *FFMPEG) {
		err := f.Command.Cmd.Wait()
		if err != nil {
			exiterr, ok := err.(*exec.ExitError)
			if ok {
				status := exiterr.Sys().(syscall.WaitStatus)
				f.Command.ExitStatus = status.ExitStatus()
			} else {
				f.Command.ExitStatus = 1
			}
		} else {
			status := f.Command.Cmd.ProcessState.Sys().(syscall.WaitStatus)
			f.Command.ExitStatus = status.ExitStatus()
		}

		f.Running = false
	}(f)

	return nil
}

func (f *FFMPEG) Stop() error {
	if f.Command.Cmd == nil {
		return nil
	}

	return f.Command.Cmd.Process.Signal(os.Interrupt)
}

func (f *FFMPEG) Kill() error {
	if f.Command.Cmd == nil {
		return nil
	}

	return f.Command.Cmd.Process.Signal(os.Kill)
}