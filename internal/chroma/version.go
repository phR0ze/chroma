package chroma

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// These are used as build time options
var (
	VERSION   string
	BUILDDATE string
	GITCOMMIT string
)

func (chroma *Chroma) newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Show version information",
		Aliases: []string{"v", "ver"},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			chroma.listVersions(os.Stdout)
			return
		},
	}
	return cmd
}

// Get the chroma version
func (chroma *Chroma) listVersions(out io.Writer) {
	fmt.Fprintf(out, `Chroma
-------------------------------------------------------------
Version:           %s
Build Date:        %s
GitCommit:         %s
`, VERSION, BUILDDATE, GITCOMMIT)
}
