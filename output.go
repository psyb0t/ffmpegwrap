package ffmpegwrap

import (
	"strings"
	"fmt"
)

type Output struct {
	IO
}

func NewOutput() *Output {
	output := &Output{}
	output.Flags = make(map[string]string)

	return output
}

func (o *Output) SetPath(path string) *Output {
	o.IO.SetPath(path)
	return o
}

func (o *Output) SetFlags(flags map[string]string) *Output {
	o.IO.SetFlags(flags)
	return o
}

func (o *Output) Compile() string {
	var flagsStrings []string
	for k, v := range o.Flags {
		flagsStrings = append(flagsStrings, k, v)
	}

	return fmt.Sprintf("%s %s", strings.Join(flagsStrings, " "), o.Path)
}