package chroma

import (
	"path"

	"github.com/phR0ze/n"
	"github.com/phR0ze/n/pkg/sys"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (chroma *Chroma) newSortCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sort [DISTROS]",
		Short: "Enable/disable patches according to the internal mapping",
		Long: `Enable/disable patches according to the internal mapping

Examples:
	
	# Sort the debian and ungoogled patches
	chroma sort debian ungoogled
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = chroma.configure(); err != nil {
				return
			}
			for _, distro := range args {
				patches := n.Map(gPatches[distro]).Keys().ToStrs()
				if err = chroma.sortPatches(distro, patches); err != nil {
					return
				}
			}
			return
		},
	}
	return cmd
}

// Enable/disable patches according to the internal gPatches mapping
func (chroma *Chroma) sortPatches(distro string, patches []string) (err error) {
	patchSetDir := path.Join(chroma.patchesDir, distro)

	for _, patch := range patches {

		// Set path name to used or not used
		dstUsedPath := path.Join(patchSetDir, patch)
		dstNotUsedPath := path.Join(patchSetDir, "not-used", patch)
		used := gPatches[distro][patch]
		switch {

		// Move not used file from used to not used directory
		case !used && sys.Exists(dstUsedPath):
			log.Infof("Disabling patch %s => %s", patch, sys.SlicePath(dstUsedPath, -3, -1))
			if _, err = sys.Move(dstUsedPath, dstNotUsedPath); err != nil {
				return
			}

		// Move used file from not used to used directory
		case used && sys.Exists(dstNotUsedPath):
			log.Infof("Enabling patch %s => %s", patch, sys.SlicePath(dstUsedPath, -3, -1))
			if _, err = sys.Move(dstNotUsedPath, dstUsedPath); err != nil {
				return
			}
		}
	}
	return
}
