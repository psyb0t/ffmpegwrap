package ffmpegwrap

type IO struct {
	Path string
	Flags map[string]string
}

func (io *IO) SetPath(path string) {
	io.Path = path
}

func (io *IO) SetFlags(flags map[string]string) {
	for k, v := range flags {
		io.Flags[k] = v
	}
}