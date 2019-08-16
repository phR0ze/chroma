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

func (chroma *CHROMA) newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Show version information",
		Aliases: []string{"v", "ver"},
		Run: func(cmd *cobra.Command, args []string) {
			chroma.listVersions(os.Stdout)
		},
	}
	return cmd
}

// Get the chroma version
func (chroma *CHROMA) listVersions(out io.Writer) {
	fmt.Fprintf(out, `Chroma
-------------------------------------------------------------
Version:           %s
Build Date:        %s
GitCommit:         %s
`, VERSION, BUILDDATE, GITCOMMIT)
}
