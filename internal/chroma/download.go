package chroma

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (chroma *CHROMA) newDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "download",
		Short:   "Download patches or extensions for chromium",
		Aliases: []string{"down"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(
		func() *cobra.Command {
			cmd := &cobra.Command{
				Use:     "extensions NAME",
				Short:   "Download the chromium extension from the Google Market",
				Aliases: []string{"ext", "exten", "extension"},
				Args:    cobra.ExactArgs(1),
				Run: func(cmd *cobra.Command, args []string) {
					chroma.configure()
					chroma.downloadExtension(args[0])
				},
			}
			return cmd
		}(),
		func() *cobra.Command {
			cmd := &cobra.Command{
				Use:     "patches DISTROS",
				Short:   "Download patches for the given distributions for chromium",
				Aliases: []string{"patch"},
				Args:    cobra.MinimumNArgs(1),
				Run: func(cmd *cobra.Command, distros []string) {
					chroma.configure()
					chroma.downloadPatches(distros)
				},
			}
			return cmd
		}(),
	)
	return cmd
}

// Download the given extension from the Google Market
func (chroma *CHROMA) downloadExtension(extName string) {
	log.Fatal(extName)
}

// Download patches for the given distributions
func (chroma *CHROMA) downloadPatches(distros []string) {
	for _, distro := range distros {
		fmt.Println(distro)
	}
}
