package chroma

import (
	"fmt"
	"path"

	"github.com/phR0ze/n"
	"github.com/phR0ze/n/pkg/net"
	"github.com/phR0ze/n/pkg/net/mech"
	"github.com/phR0ze/n/pkg/sys"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type downloadOpts struct {
	clean bool // remove previous files before downloading
}

func (chroma *Chroma) newDownloadCmd() *cobra.Command {
	opts := &downloadOpts{}
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
					if err := chroma.downloadPatches(distros, opts); err != nil {
						chroma.logFatal(err)
					}
				},
			}
			cmd.Flags().BoolVar(&opts.clean, "clean", false, "Remove local files before downloading")
			return cmd
		}(),
	)
	return cmd
}

// Download the given extension from the Google Market
func (chroma *Chroma) downloadExtension(extName string) {
	log.Fatal(extName)
}

// Download patches for the given distributions
func (chroma *Chroma) downloadPatches(distros []string, opts *downloadOpts) (err error) {
	for _, distro := range distros {
		patchSetDir := path.Join(chroma.patchesDir, distro)
		notUsedDir := path.Join(patchSetDir, "not-used")

		// Ensure destination directory is clean and ready
		// -----------------------------------------------------------------------------------------
		if opts.clean {
			log.Infof("Removing all files in local patchset dir %s", patchSetDir)
			if sys.Exists(patchSetDir) {
				sys.RemoveAll(patchSetDir)
			}
		}
		if _, err = sys.MkdirP(notUsedDir); err != nil {
			return
		}

		// Download and process links from the patchset page
		// -----------------------------------------------------------------------------------------
		log.Infof("Downloading patchset %s => %s", distro, patchSetDir)
		agent := mech.New()

		// Handle each distro differently
		// -----------------------------------------------------------------------------------------
		switch distro {
		case "debian":
			var order *n.StringSlice
			if order, err = readOrderFile(distro, patchSetDir); err != nil {
				return
			}

			// Download each of the patches numbering and naming them according to the order file
			for i, entry := range order.SG() {
				uri := net.JoinURL(net.DirURL(gPatchSets[distro]), entry)
				dstName := fmt.Sprintf("%02d-%s", i, path.Base(entry))
				if err = downloadPatch(agent, uri, distro, patchSetDir, dstName); err != nil {
					return
				}
			}
		case "ungoogled":
			var order *n.StringSlice
			if order, err = readOrderFile(distro, patchSetDir); err != nil {
				return
			}

			// Download each of the patches numbering and naming them according to the order file
			for i, entry := range order.SG() {
				uri := net.JoinURL(net.DirURL(gPatchSets[distro]), entry)
				dstName := fmt.Sprintf("%02d-%s", i, path.Base(entry))
				if err = downloadPatch(agent, uri, distro, patchSetDir, dstName); err != nil {
					return
				}
			}
		}
	}
	return
}

// Download the given patch set or relocate it if needed
func downloadPatch(agent *mech.Mech, uri, distro, patchSetDir, dstName string) (err error) {

	// Set path name to used or not used
	dstUsedPath := path.Join(patchSetDir, dstName)
	dstNotUsedPath := path.Join(patchSetDir, "not-used", dstName)
	used := gPatches[distro][dstName]
	switch {

	// Move not used file from used to not used directory
	case !used && sys.Exists(dstUsedPath):
		log.Infof("Disabling patch %s => %s", dstName, sys.SlicePath(dstUsedPath, -3, -1))
		if err = sys.Move(dstUsedPath, dstNotUsedPath); err != nil {
			return
		}

	// Move used file from not used to used directory
	case used && sys.Exists(dstNotUsedPath):
		log.Infof("Enabling patch %s => %s", dstName, sys.SlicePath(dstUsedPath, -3, -1))
		if err = sys.Move(dstNotUsedPath, dstUsedPath); err != nil {
			return
		}

	case !sys.Exists(dstUsedPath) && !sys.Exists(dstNotUsedPath):
		dstPath := dstUsedPath
		if !used {
			dstPath = dstNotUsedPath
		}
		log.Infof("Downloading patch %s => %s", sys.SlicePath(uri, -3, -1), sys.SlicePath(dstPath, -2, -1))
		if _, err = agent.Download(uri, dstPath); err != nil {
			return
		}
	}
	return
}

// Read the order files from disk, downloading if it doesn't exist
func readOrderFile(distro, patchSetDir string) (order *n.StringSlice, err error) {

	// Read in the patch order file, downloading if needed
	orderFile := path.Join(patchSetDir, path.Base(gPatchSets[distro]))
	if !sys.Exists(orderFile) {
		log.Infof("Downloading patch order file %s", gPatchSets[distro])
		if _, err = mech.Download(gPatchSets[distro], orderFile); err != nil {
			return
		}
	}

	// Read in the order file
	var data []string
	if data, err = sys.ReadLines(orderFile); err != nil {
		return
	}
	order = n.S(data)

	// Trim out any empty lines
	order.DropW(func(x n.O) bool {
		return n.ExB(x.(string) == "")
	})

	return
}
