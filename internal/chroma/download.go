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
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.Flags().BoolVar(&opts.clean, "clean", false, "Remove local files before downloading")
	cmd.AddCommand(
		func() *cobra.Command {
			cmd := &cobra.Command{
				Use:   "extensions [NAME]",
				Short: "Download the chromium extension from the Google Market",
				Long: `Download the chromium extensions from the Google Market.

Examples:
	# Download all extensions
	chroma down ext
	
	# Download the ublock-origin extension
	chroma down ext ublock-origin

	# Download the https-everywhere and ublock-origin extensions
	chroma down ext https-everywhere ublock-origin
`,
				Aliases: []string{"ex", "ext", "exten", "extension"},
				RunE: func(cmd *cobra.Command, args []string) (err error) {
					if err = chroma.configure(); err != nil {
						return
					}
					if err = chroma.downloadExtensions(args, opts); err != nil {
						return
					}
					return
				},
			}
			return cmd
		}(),
		func() *cobra.Command {
			cmd := &cobra.Command{
				Use:   "patches [DISTROS]",
				Short: "Download patches for the given distributions for chromium",
				Long: `Download patches for the given distributions for chromium.

Examples:
	# Download all patches
	chroma down patches
	
	# Download the debian and ungoogled patches
	chroma down patches debian ungoogled
`,
				Aliases: []string{"pa", "patch"},
				Args:    cobra.MinimumNArgs(1),
				RunE: func(cmd *cobra.Command, args []string) (err error) {
					if err = chroma.configure(); err != nil {
						return
					}
					if err = chroma.downloadPatches(args, opts); err != nil {
						return
					}
					return
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
			} else {
				err = errors.Errorf("Error: unsupported extension %s", name)
				return
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
		crxfile := path.Join(chroma.extensionsDir, fmt.Sprintf("%s.crx", extName))

		// Download the extension if it doesn't yet exist
		if !sys.Exists(crxfile) {
			log.Infof("Downloading extension %s:%s => %s", extName, extID, sys.SlicePath(crxfile, -3, -1))
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
			if _, err = agent.Download(uri.String(), crxfile); err != nil {
				return
			}
		}

		// Generate the JSON preferences file
		prefPath := path.Join(chroma.extensionsDir, fmt.Sprintf("%s.json", extID))
		if !sys.Exists(prefPath) {
			log.Infof("Generating extension preferences file for %s", extName)

			// Unzip the extension
			tmpDir := path.Join(chroma.extensionsDir, "_tmp")
			log.Infof("Unzipping the extension %s => %s", extName, sys.SlicePath(tmpDir, -3, -1))
			if sys.Exists(tmpDir) {
				sys.RemoveAll(tmpDir)
			}
			if err = zip.ExtractAll(crxfile, tmpDir); err != nil {
				return
			}

			// Read in the extension's manifest.json file
			var m *n.StringMap
			jsonfile := path.Join(tmpDir, "manifest.json")
			if m, err = n.LoadJSONE(jsonfile); err != nil {
				return
			}
			extVer := m.Query("version").A()
			if extVer == "" {
				err = errors.Errorf("failed to extract version from ext manifest file")
				return
			}
			log.Infof("Extracted extension version: %s", extVer)

			// Preferences file
			// https://developer.chrome.com/apps/external_extensions
			prefs := n.NewStringMap(map[string]interface{}{
				"external_crx":     path.Join("/usr/share/chromium/extensions", path.Base(crxfile)),
				"external_version": extVer,
				// Rather than the local file external_crx we can use the upate url below to download them
				//"external_update_url": "https://clients2.google.com/service/update2/crx",
			})
			log.Infof("Creating preference file %s", sys.SlicePath(prefPath, -3, -1))
			if err = prefs.WriteJSON(prefPath); err != nil {
				return
			}
			sys.RemoveAll(tmpDir)
		} else {
			log.Infof("Extension preferences file for %s already exists", extName)
		}
	}

	return
}

// Download patches for the given distributions
func (chroma *Chroma) downloadPatches(distros []string, opts *downloadOpts) (err error) {
	if len(distros) == 0 {
		distros = []string{"debian", "ungoogled"}
	}
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
		if _, err = sys.Move(dstUsedPath, dstNotUsedPath); err != nil {
			return
		}

	// Move used file from not used to used directory
	case used && sys.Exists(dstNotUsedPath):
		log.Infof("Enabling patch %s => %s", dstName, sys.SlicePath(dstUsedPath, -3, -1))
		if _, err = sys.Move(dstNotUsedPath, dstUsedPath); err != nil {
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
