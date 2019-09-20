package chroma

import "github.com/phR0ze/n/pkg/opt"

var (
	chromiumPath = "~/Projects/cyberlinux/aur/chromium"
)

func newChroma(opts ...*opt.Opt) *Chroma {
	opt.Add(&opts, RootOpt(chromiumPath))

	c := New(opts...)

	return c
}
