package chroma

import (
	"github.com/phR0ze/n/pkg/opt"
)

var (
	// RootOptKey for the root option
	RootOptKey = "root"
)

// Root option
// -------------------------------------------------------------------------------------------------

// RootOpt passes the path to where the chromium PKGBUILD is
func RootOpt(val string) *opt.Opt {
	return &opt.Opt{Key: RootOptKey, Val: val}
}

// GetRootOpt return the root of hte chromium PKGBUILD
func GetRootOpt(opts []*opt.Opt) string {
	if o := opt.Get(opts, RootOptKey); o != nil {
		if val, ok := o.Val.(string); ok {
			return val
		}
	}
	return ""
}

// DefaultRootOpt returns the option value if found else the one given
// Use this when the Get's default is not desirable.
func DefaultRootOpt(opts []*opt.Opt, val string) string {
	if !opt.Exists(opts, RootOptKey) {
		return val
	}
	return GetRootOpt(opts)
}
