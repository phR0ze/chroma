package chroma

import (
	"fmt"
	"net/url"
	"path"

	"github.com/phR0ze/n"
	"github.com/phR0ze/n/pkg/arch/zip"
	"github.com/phR0ze/n/pkg/net"
	"github.com/phR0ze/n/pkg/net/mech"
	"github.com/phR0ze/n/pkg/sys"
	"github.com/pkg/errors"
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
		Aliases: []string{"do", "down"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.Flags().BoolVar(&opts.clean, "clean", false, "Remove local files before downloading")
	cmd.AddCommand(
		func() *cobra.Command {
			cmd := &cobra.Command{
				Use:     "extensions [NAME]",
				Short:   "Download the chromium extension from the Google Market",
				Aliases: []string{"ex", "ext", "exten", "extension"},
				Run: func(cmd *cobra.Command, args []string) {
					chroma.configure()
					if err := chroma.downloadExtensions(args, opts); err != nil {
						chroma.logFatal(err)
					}
				},
			}
			return cmd
		}(),
		func() *cobra.Command {
			cmd := &cobra.Command{
				Use:     "patches DISTROS",
				Short:   "Download patches for the given distributions for chromium",
				Aliases: []string{"pa", "patch"},
				Args:    cobra.MinimumNArgs(1),
				Run: func(cmd *cobra.Command, args []string) {
					chroma.configure()
					if err := chroma.downloadPatches(args, opts); err != nil {
						chroma.logFatal(err)
					}
				},
			}
			return cmd
		}(),
	)
	return cmd
}

// Download the given extension from the Google Market
func (chroma *Chroma) downloadExtensions(extnames []string, opts *downloadOpts) (err error) {
	log.Infof("Downloading extensions => %s", chroma.extensionsDir)
	exts := map[string]string{}

	// Select extensions is given
	if len(extnames) == 0 {
		exts = gExtensions
	} else {
		for _, name := range extnames {
			if val, ok := gExtensions[name]; ok {
				exts[name] = val
			}
		}
	}

	// Ensure destination directory is clean and ready
	// -----------------------------------------------------------------------------------------
	if opts.clean {
		log.Infof("Removing all files in local extensions dir %s", chroma.extensionsDir)
		if sys.Exists(chroma.extensionsDir) {
			sys.RemoveAll(chroma.extensionsDir)
		}
	}
	if _, err = sys.MkdirP(chroma.extensionsDir); err != nil {
		return
	}

	// Download and process links from the patchset page
	// -----------------------------------------------------------------------------------------
	agent := mech.New()
	for extName, extID := range exts {
		zipfile := path.Join(chroma.extensionsDir, fmt.Sprintf("%s.crx", extName))
		if !sys.Exists(zipfile) {

			// Download the extension
			log.Infof("Downloading extension %s:%s => %s", extName, extID, sys.SlicePath(zipfile, -3, -1))
			var uri *url.URL
			if uri, err = url.Parse("https://clients2.google.com/service/update2/crx"); err != nil {
				err = errors.Wrapf(err, "failed to parse webstore url")
				return
			}
			uri.RawQuery = url.Values{
				"response":    {"redirect"},
				"os":          {"linux"},
				"prodversion": {chroma.chromiumVer},
				"x":           {fmt.Sprintf("id=%s&installsource=ondemand&uc", extID)},
			}.Encode()
			if _, err = agent.Download(uri.String(), zipfile); err != nil {
				return
			}

			// Unzip the extension
			log.Infof("Unzipping the extension %s:%s", extName, extID)
			tmpDir := path.Join(chroma.extensionsDir, "_tmp")
			if sys.Exists(tmpDir) {
				sys.RemoveAll(tmpDir)
			}
			zip.ExtractAll(zipfile, tmpDir)

			// Generate the JSON preferences file
			log.Info("Generating extension preferences file")

		}
	}

	return
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
			for i, entry := range order.ToStrs() {
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
			for i, entry := range order.ToStrs() {
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
