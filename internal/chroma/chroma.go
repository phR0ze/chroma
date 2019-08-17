package chroma

import (
	"fmt"
	"os"
	"path"

	"github.com/phR0ze/n"
	"github.com/phR0ze/n/pkg/sys"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Distros defines the supported distribution patch sets
type Distros struct {
	arch    string
	cyber   string
	debian  string
	inox    string
	iridium string
}

var (
	gDistros = Distros{"arch", "cyber", "debian", "inox", "iridium"}

	// Supported extensions
	gExtensions = map[string]string{
		"https-everywhere":     "gcbommkclmclpchllfjekcdonpmejbdp", // Automatically use HTTPS security where possible
		"scriptsafe":           "oiigbmnaadbkfbmpbfijlflahbdbdgdf", //
		"smartup-gestures":     "bgjfekefhjemchdeigphccilhncnjldn", // Better mouse gestures for Chromium
		"tampermonkey":         "dhdgffkkebhmkfjojejmpbldmpobfkfo", // World's most popular userscript manager
		"ublock-origin":        "cjpalhdlnbpafiamejdnhcphjbkeiagm", // An efficient ad-blocker for Chromium
		"ublock-origin-extra":  "pgdnlhfefecpicbbihgmbmffkjpaplco", // Foil early hostile anti-user mechanisms
		"umatrix":              "ogfcmafjalglgifnmanfmnieipoejdcf", //
		"videodownload-helper": "lmjnegcaeklhafolokijcfjliaokphfk", // Video download helper for Chromium
	}

	// Supported patch sets
	gPatchSets = map[string]string{
		"arch":    "https://git.archlinux.org/svntogit/packages.git/tree/trunk?h=packages/chromium",
		"debian":  "https://salsa.debian.org/chromium-team/chromium/tree/master/debian/patches",
		"inox":    "https://github.com/gcarq/inox-patchset",
		"iridium": "https://git.iridiumbrowser.de/cgit.cgi/iridium-browser/commit/?h=patchview",
	}

	// Call out patches used and not used and notes
	// Order is significant
	// --------------------------------------------------------------------------
	gUsedPatches = map[string][]string{
		gDistros.arch: {
			"breakpad-use-ucontext_t.patch",   // Glibc 2.26 does not expose struct ucontext any longer
			"chromium-gn-bootstrap-r17.patch", //
		},

		// Credit to Michael Gilber
		gDistros.debian: {
			"manpage.patch", // Adds simple doc with link to documentation website

			"gn/parallel.patch",   // Respect specified number of parllel jobs when bootstrapping
			"gn/narrowing.patch",  // Silence narrowing warnings when bootstrapping gn
			"gn/buildflags.patch", // Support build flags passed in the --args to gn

			"disable/promo.patch",                // Disable ad promo system by default
			"disable/fuzzers.patch",              // Disable fuzzers as they aren't built anyway and only used for testing
			"disable/google-api-warning.patch",   // Disables Google's API key warning when they are removed from the PKGBUILD
			"disable/external-components.patch",  // Disable loading: Enhanced bookmarks, HotWord, ZipUnpacker, GoogleNow
			"disable/device-notifications.patch", // Disable device discovery notifications

			"fixes/mojo.patch",                  // Fix mojo layout test build error
			"fixes/chromecast.patch",            // Disable chromecast unless flag GOOGLE_CHROME_BUILD set
			"fixes/ps-print.patch",              // Add postscript(ps) printing capabiliy
			"fixes/gpu-timeout.patch",           // Increase GPU timeout from 10sec to 20sec
			"fixes/widevine-revision.patch",     // Set widevine version as undefined
			"fixes/connection-message.patch",    // Update connection message to suggest updating your proxy if you can't get connected.
			"fixes/chromedriver-revision.patch", // Set as undefined, Chromedriver allows for automated testing of chromium
		},

		gDistros.cyber: {
			"00-master-preferences.patch",         // Configure the master preferences to be in /etc/chromium/master_preferences
			"01-disable-default-extensions.patch", // Apply on top of debian patches, disables cloud print and feedback
		},

		// Credit to Michael Egger -> patches/inox/LICENSE
		// https://github.com/gcarq/inox-patchset
		//
		// Default Settings
		// ------------------------------------------------------------------------
		// DefaultCookiesSettings                            CONTENT_SETTING_DEFAULT
		// EnableHyperLinkAuditing 	                        false
		// CloudPrintSubmitEnabled 	                        false
		// NetworkPredictionEnabled 	                        false
		// BackgroundModeEnabled 	                          false
		// BlockThirdPartyCookies 	                          true
		// AlternateErrorPagesEnabled 	                      false
		// SearchSuggestEnabled 	                            false
		// AutofillEnabled 	                                false
		// Send feedback to Google if preferences are reset 	false
		// BuiltInDnsClientEnabled                         	false
		// SignInPromoUserSkipped 	                          true
		// SignInPromoShowOnFirstRunAllowed 	                false
		// ShowAppsShortcutInBookmarkBar 	                  false
		// ShowBookmarkBar 	                                true
		// PromptForDownload 	                              true
		// SafeBrowsingEnabled 	                            false
		// EnableTranslate 	                                false
		// LocalDiscoveryNotificationsEnabled 	              false
		gDistros.inox: {
			"0001-fix-building-without-safebrowsing.patch", // Required when the PGKBUILD has safebrowing disabled
			"0003-disable-autofill-download-manager.patch", // Disables HTML AutoFill data transmission to Google
			"0006-modify-default-prefs.patch",              // Set default settings as described in header
			"0007-disable-web-resource-service.patch",      //
			"0008-restore-classic-ntp.patch",               // The new NTP (New Tag Page) pulls from Google including tracking identifier
			"0009-disable-google-ipv6-probes.patch",        // Change IPv6 DNS probes to Google over to k.root-servers.net
			"0010-disable-gcm-status-check.patch",          // Disable Google Cloud-Messaging status probes, GCM allows direct msg to device
			"0014-disable-translation-lang-fetch.patch",    // Disable language fetching from Google when settings are opened the first time
			"0015-disable-update-pings.patch",              // Disable update pings to Google
			"0016-chromium-sandbox-pie.patch",              // Hardening sandbox with Position Independent code, originally from openSUSE
			"0017-disable-new-avatar-menu.patch",           // Disable Google Avatar signin menu
			"0018-disable-first-run-behaviour.patch",       // Modifies first run to prevent data leakage
			"0019-disable-battery-status-service.patch",    // Disable battry status service as it can be used for tracking
			"0021-disable-rlz.patch",                       // Disable RLZ
			"9000-disable-metrics.patch",                   // Disable metrics
			"9001-disable-profiler.patch",                  // Disable profiler
		},
	}

	// Not used patches and a description as to why not
	// --------------------------------------------------------------------------
	gNotUsedPatches = map[string][]string{
		gDistros.arch: {
			"chromium-widevine.patch", // Using debian as this one uses a variable
		},
		gDistros.debian: {
			"master-preferences.patch",          // Use custom cyber patch instead
			"disable/third-party-cookies.patch", // Already covered in inox/0006-modify-default-prefs'
			"gn/bootstrap.patch",                // Fix errors in gn's bootstrapping script, using arch bootstrap instead
			"fixes/crc32.patch",                 // Fix inverted check, using arch crc32c-string-view-check.patch instead
			"system/nspr.patch",                 // Build using the system nspr library
			"system/icu.patch",                  // Backwards compatibility for older versions of icu
			"system/vpx.patch",                  // Remove VP9 support because debian libvpx doesn"t support VP9 yet
			"system/gtk2.patch",                 //
			"system/lcms2.patch",                //
			"system/event.patch",                // Build using the system libevent library
		},
		gDistros.inox: {
			// Disables Hotword, Google Now/Feedback/Webstore/Hangout, Cloud Print, Speech synthesis
			// I like keeping the Webstore and Hangout features so will roll my own patch in cyberlinux
			"0004-disable-google-url-tracker.patch",   // No URL tracking (Google saves your location) also breaks omnibar search
			"0005-disable-default-extensions.patch",   // see above
			"0011-add-duckduckgo-search-engine.patch", // Adds DuckDuckGo as default search engine, still changeable in settings
			"0012-branding.patch",                     // Want to keep the original Chromium branding
			"0013-disable-missing-key-warning.patch",  // Disables warning, using debian patch instead
			"0020-launcher-branding.patch",            // Want to keep the original Chromium branding
			"breakpad-use-ucontext_t.patch",           // Already included by Arch Linux
			"chromium-gn-bootstrap-r17.patch",         // Already included by Arch Linux
			"chromium-widevine.patch",                 // Already included by Arch Linux
			"crc32c-string-view-check.patch",          // Already included by Arch Linux
			"chromium-libva-version.patch",            // Arch doesn't use it
			"chromium-vaapi-r14.patch",                // Arch doesn't use it
		},
	}
)

// Opts allows for passing in complicated arguments
type Opts struct {
	Root    string // path where the chromium PKGBUILD will be located
	Testing bool
}

// Chroma instance
type Chroma struct {
	cmd    *cobra.Command // root cobra command
	quiet  bool           // don't emit anything except errors
	debug  bool           // debug this session
	dryrun bool           // make no changes

	// Custom app state
	rootDir       string // path to dir where chromium PKGBUILD is
	pkgbuild      string // path to the PKGBUILD
	patchesDir    string // path to the patches dir in the chromium package
	extensionsDir string // path to the src/exentions dir in the chromium package
	version       string // target version of chrome pulled from the PKGBUILD
}

// New initializes the CLI with the given options
func New(o ...*Opts) (chroma *Chroma) {
	chroma = &Chroma{}

	// Configure startup options
	//----------------------------------------------------------------------------------------------
	opts := &Opts{}
	if len(o) > 0 {
		opts = o[0]
	}

	n.SetOnEmpty(&VERSION, "999.999.999")
	var boilerPlate = fmt.Sprintf("Chroma %s [%s (Git %s)]\n", VERSION, BUILDDATE, GITCOMMIT)

	if opts.Root == "" {
		chroma.rootDir = sys.Pwd()
	} else {
		chroma.rootDir = opts.Root
	}

	// Configure Cobra CLI commands
	// All subcommands should be defined in a seperate file.
	//----------------------------------------------------------------------------------------------
	chroma.cmd = &cobra.Command{
		Use: "chroma",
		Long: fmt.Sprintf(`%sautomation for chromium patches and plugins.

Examples:
  # Check current versions
  chroma version

  # Check current persisted context
  chroma use
`,
			boilerPlate),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	chroma.cmd.AddCommand(
		chroma.newDownloadCmd(),
		chroma.newVersionCmd(),
	)

	// --debug
	chroma.cmd.PersistentFlags().BoolVar(&chroma.debug, "debug", false, "Print out debug info")

	// --dry-run
	chroma.cmd.PersistentFlags().BoolVar(&chroma.dryrun, "dry-run", false, "Make no changes")

	// -q --quiet
	chroma.cmd.PersistentFlags().BoolVarP(&chroma.quiet, "quiet", "q", false, "Use Error level logging and don't emit extraneous output")

	// --pkgbuild
	chroma.cmd.PersistentFlags().StringVar(&chroma.pkgbuild, "pkgbuild", "", "Use this specific PKGBUILD to derive pathes from")

	// Setup logging after we've read in the env variables
	chroma.setupLogging()

	return
}

// Execute the CLI
func (chroma *Chroma) Execute(args ...string) (err error) {
	if len(args) > 0 {
		os.Args = append([]string{"chroma"}, args...)
	}
	return chroma.cmd.Execute()
}

// configure should be run before anything from the Cobra Run functions
// to configure and validate the environment.
func (chroma *Chroma) configure() (err error) {

	// Configure paths based off pkgbuild and root path
	// ---------------------------------------------------------------------------------------------
	if chroma.pkgbuild != "" {
		chroma.rootDir = path.Dir(chroma.pkgbuild)
	} else {
		chroma.pkgbuild = path.Join(chroma.rootDir, "PKGBUILD")
	}
	chroma.patchesDir = path.Join(chroma.rootDir, "patches")
	chroma.extensionsDir = path.Join(chroma.rootDir, "src", "extensions")

	// Validate the chromium PKGBUILD
	// ---------------------------------------------------------------------------------------------
	if !sys.Exists(chroma.pkgbuild) {
		err = errors.Errorf("chromium PKGBUILD coudn't be found")
		chroma.logFatal(err)
	}

	// Parse out the chromium version from the PKGBUILD
	exp := `(?m)^pkgver=(.*)$`
	if chroma.version, err = sys.ExtractString(chroma.pkgbuild, exp); err != nil || chroma.version == "" {
		err = errors.Errorf("failed to extract the chromium version from the PKGBUILD")
		chroma.logFatal(err)
	}

	// Boiler plate for all commands
	// ---------------------------------------------------------------------------------------------
	chroma.printf("Chromium Ver:    %s\n", chroma.version)
	chroma.printf("PKBUILD Path:    %s\n", chroma.pkgbuild)
	chroma.println()
	return
}
