package ffmpegwrap

import (
	"strings"
	"fmt"
)

type Input struct {
	IO
}

func NewInput() *Input {
	input := &Input{}
	input.Flags = make(map[string]string)

	return input
}

func (i *Input) SetPath(path string) *Input {
	i.IO.SetPath(path)
	return i
}

func (i *Input) SetFlags(flags map[string]string) *Input {
	i.IO.SetFlags(flags)
	return i
}

func (i *Input) Compile() string {
	var flagsStrings []string
	for k, v := range i.Flags {
		flagsStrings = append(flagsStrings, k, v)
	}

	return fmt.Sprintf("%s -i %s", strings.Join(flagsStrings, " "), i.Path)
}