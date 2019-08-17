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

	// Order files for supported patch sets
	gPatchSets = map[string]string{
		"arch":    "https://git.archlinux.org/svntogit/packages.git/tree/trunk?h=packages/chromium",
		"debian":  "https://salsa.debian.org/chromium-team/chromium/raw/master/debian/patches/series",
		"inox":    "https://github.com/gcarq/inox-patchset",
		"iridium": "https://git.iridiumbrowser.de/cgit.cgi/iridium-browser/commit/?h=patchview",
	}

	// Call out patches used and not used and notes
	// Order is significant
	// --------------------------------------------------------------------------
	gPatches = map[string]map[string]bool{
		// gDistros.arch: {
		// 	"breakpad-use-ucontext_t.patch",   // Glibc 2.26 does not expose struct ucontext any longer
		// 	"chromium-gn-bootstrap-r17.patch", //
		// },

		// Credit to Michael Gilber
		// 	"disable/third-party-cookies.patch", // Already covered in inox/0006-modify-default-prefs'
		gDistros.debian: {
			"00-manpage.patch":                  true,  // Adds simple doc with link to documentation website
			"01-sandbox.patch":                  false, // Debian specific error message to install chromium-sandbox
			"02-master-preferences.patch":       true,  // Look for master preferences in /etc/chromium/master_preferences
			"03-libcxx.patch":                   true,  // Avoid chromium's embedded C++ library when bootstrapping
			"04-parallel.patch":                 true,  // Respect specified number of parllel jobs when bootstrapping
			"05-gcc_skcms_ice.patch":            true,  // GCC ICE with optimized version
			"06-pffffft-buildfix.patch":         true,  // ??
			"07-skia-aarch64-buildfix.patch":    true,  // ??
			"08-wrong-namespace.patch":          true,  // gcc: various methods and classes are using the wrong namespace
			"09-virtual-destructor.patch":       true,  // gcc: a virtual destructor is called without this patch
			"10-explicit-specialization.patch":  true,  // gcc: fix for gcc explicit specialiazation namespace issue
			"11-macro.patch":                    true,  // gcc6: can be ignored as arch linux is using gcc 9
			"12-sizet.patch":                    true,  // gcc6: can be ignored as arch linux is using gcc 9
			"13-atomic.patch":                   true,  // gcc6: can be ignored as arch linux is using gcc 9
			"14-constexpr.patch":                true,  // gcc6: can be ignored as arch linux is using gcc 9
			"15-wtf-hashmap.patch":              true,  // gcc6: can be ignored as arch linux is using gcc 9
			"16-lambda-this.patch":              true,  // gcc6: can be ignored as arch linux is using gcc 9
			"17-map-insertion.patch":            true,  // gcc6: can be ignored as arch linux is using gcc 9
			"18-not-constexpr.patch":            true,  // gcc6: can be ignored as arch linux is using gcc 9
			"19-move-required.patch":            true,  // gcc6: can be ignored as arch linux is using gcc 9
			"20-use-after-move.patch":           true,  // gcc6: can be ignored as arch linux is using gcc 9
			"21-ambiguous-overloads.patch":      true,  // gcc6: can be ignored as arch linux is using gcc 9
			"22-ambiguous-initializer.patch":    true,  // gcc6: can be ignored as arch linux is using gcc 9
			"23-nullptr-copy-construct.patch":   true,  // gcc6: can be ignored as arch linux is using gcc 9
			"24-noexcept-redeclaration.patch":   true,  // gcc6: can be ignored as arch linux is using gcc 9
			"25-trivially-constructible.patch":  true,  // gcc6: can be ignored as arch linux is using gcc 9
			"26-designated-initializers.patch":  true,  // gcc6: can be ignored as arch linux is using gcc 9
			"27-specialization-namespace.patch": true,  // gcc6: can be ignored as arch linux is using gcc 9
			"28-mojo.patch":                     true,  // Fixes: fix mojo layout test build error
			"29-public.patch":                   true,  // Fixes: method needs to be public
			"30-ps-print.patch":                 true,  // Fixes: add postscript(ps) printing capabiliy
			"31-as-needed.patch":                true,  // Fixes: some libraries fail to link when '--as-needed' is set
			"32-inspector.patch":                true,  // Fixes: use inspector_protocol from top level third_party dir
			"33-gpu-timeout.patch":              true,  // Fixes: increase GPU timeout from 10sec to 20sec
			"34-empty-array.patch":              true,  // Fixes: arraysize macro fails for zero length array and add one char
			"35-safebrowsing.patch":             true,  // Fixes: fix signedness error when built with gcc affects safe browsing
			"36-sequence-point.patch":           true,  // Fixes: fix undefined order in which expressions are evaluated
			"37-jumbo-namespace.patch":          true,  // Fixes: jumbo build has troubel with these namespaces
			"38-template-export.patch":          true,  // Fixes: implementation of template function must be in header to be exported
			"39-widevine-revision.patch":        true,  // Fixes: set widevine version as undefined
			"40-widevine-locations.patch":       false, // Fixes: arch linux works fine don't need to try alternative location for widevine
			"41-widevine-buildflag.patch":       true,  // Fixes: enable widevine support
			"42-connection-message.patch":       false, // Fixes: hardly seems importan to 'update suggest updating your proxy when network is unreachable'
			"43-unrar.patch":                    true,  // Disable: disable support for browsing rar files
			"44-signin.patch":                   true,  // Disable: disable browser sign-in
			"45-android.patch":                  true,  // Disable: disable dependency on chrome/android
			"46-fuzzers.patch":                  true,  // Disable: fuzzers as they aren't built anyway and only used for testing
			"47-tracing.patch":                  true,  // Disable: disable tracing which depends on too many sourceless javascript files
			"48-openh264.patch":                 false, // Disable: disable support for openh264
			"49-chromeos.patch":                 true,  // Disable: ??
			"50-perfetto.patch":                 true,  // Disable: disable dependencies on third_party perfetto
			"51-installer.patch":                true,  // Disable: avoid building the chromium installer
			"52-font-tests.patch":               true,  // Disable: disable building font tests
			"53-swiftshader.patch":              true,  // Disable: avoid building the swiftshader library
			"54-welcome-page.patch":             true,  // Disable: do not override the welcome page setting in preferences
			"55-google-api-warning.patch":       true,  // Disable: disable Google's API key warning when they are removed from the PKGBUILD
			"56-third-party-cookies.patch":      true,  // Disable: disable third-party cookies in preferences
			"57-device-notifications.patch":     true,  // Disable: disable device discovery notifications in preferences
			"58-int32.patch":                    true,  // Warning: fit int32_t enum values into 32 bits
			"59-friend.patch":                   true,  // Warning: unfriend classses that friend themselves
			"60-printf.patch":                   true,  // Warning: cast enums to int for use as printf arguments
			"61-attribute.patch":                true,  // Warning: fix gcc optimization but attribute doesn't match warnings
			"62-multichar.patch":                true,  // Warning: crashpad relies on multicharacter integer assignments
			"63-deprecated.patch":               true,  // Warning: ignore deprecated bison directive warnings
			"64-bool-compare.patch":             true,  // Warning: fix gcc bool-compare warnings
			"65-enum-compare.patch":             true,  // Warning: fix gcc warnings about enum comparisions
			"66-sign-compare.patch":             true,  // Warning: fix gcc sign-compare warnings
			"67-initialization.patch":           true,  // Warning: source could be uninitialized
			"68-unused-typedefs.patch":          true,  // Warning: fix type in unused local typedefs
			"69-unused-functions.patch":         true,  // Warning: remove functions that are unused
			"70-null-destination.patch":         true,  // Warning: use stack_buf before possible branching
			"71-int-in-bool-context.patch":      true,  // Warning: fix int in bool context gcc warnings
			"72-vpx.patch":                      false, // System: arch linux supports VP9 so we don't need to disable it in libvpx
			"73-icu.patch":                      false, // System: arch linux PKGBUILD has a system lib call out for this already
			"74-gtk2.patch":                     false, // System: arch linux packages work fine when building against GTK3
			"75-jpeg.patch":                     false, // System: arch linux PKGBUILD has a system lib call out for this already
			"76-lcms.patch":                     true,  // System: use system lcms for pdfium
			"77-nspr.patch":                     true,  // System: build using the system nspr library
			"78-zlib.patch":                     false, // System: arch PKGBUILD has a system lib call out for this already
			"79-event.patch":                    true,  // System: build using the system libevent library
			"80-ffmpeg.patch":                   false, // System: arch linux PKGBUILD has a system lib call out for this already
			"81-jsoncpp.patch":                  true,  // System: use system jsoncpp
			"82-openjpeg.patch":                 true,  // System: build system using openjpeg
			"83-convertutf.patch":               true,  // System: use ICU for UTF8 conversions (eleminates ConvertUTF embedded code copy)
			"84-icu63.patch":                    false, // System: arch linux has newer icu don't need to maintain compt with 63
		},

		// gDistros.cyber: {
		// 	"01-disable-default-extensions.patch", // Apply on top of debian patches, disables cloud print and feedback
		// },

		// // Credit to Michael Egger -> patches/inox/LICENSE
		// // https://github.com/gcarq/inox-patchset
		// //
		// // Default Settings
		// // ------------------------------------------------------------------------
		// // DefaultCookiesSettings                            CONTENT_SETTING_DEFAULT
		// // EnableHyperLinkAuditing 	                        false
		// // CloudPrintSubmitEnabled 	                        false
		// // NetworkPredictionEnabled 	                        false
		// // BackgroundModeEnabled 	                          false
		// // BlockThirdPartyCookies 	                          true
		// // AlternateErrorPagesEnabled 	                      false
		// // SearchSuggestEnabled 	                            false
		// // AutofillEnabled 	                                false
		// // Send feedback to Google if preferences are reset 	false
		// // BuiltInDnsClientEnabled                         	false
		// // SignInPromoUserSkipped 	                          true
		// // SignInPromoShowOnFirstRunAllowed 	                false
		// // ShowAppsShortcutInBookmarkBar 	                  false
		// // ShowBookmarkBar 	                                true
		// // PromptForDownload 	                              true
		// // SafeBrowsingEnabled 	                            false
		// // EnableTranslate 	                                false
		// // LocalDiscoveryNotificationsEnabled 	              false
		// gDistros.inox: {
		// 	"0001-fix-building-without-safebrowsing.patch", // Required when the PGKBUILD has safebrowing disabled
		// 	"0003-disable-autofill-download-manager.patch", // Disables HTML AutoFill data transmission to Google
		// 	"0006-modify-default-prefs.patch",              // Set default settings as described in header
		// 	"0007-disable-web-resource-service.patch",      //
		// 	"0008-restore-classic-ntp.patch",               // The new NTP (New Tag Page) pulls from Google including tracking identifier
		// 	"0009-disable-google-ipv6-probes.patch",        // Change IPv6 DNS probes to Google over to k.root-servers.net
		// 	"0010-disable-gcm-status-check.patch",          // Disable Google Cloud-Messaging status probes, GCM allows direct msg to device
		// 	"0014-disable-translation-lang-fetch.patch",    // Disable language fetching from Google when settings are opened the first time
		// 	"0015-disable-update-pings.patch",              // Disable update pings to Google
		// 	"0016-chromium-sandbox-pie.patch",              // Hardening sandbox with Position Independent code, originally from openSUSE
		// 	"0017-disable-new-avatar-menu.patch",           // Disable Google Avatar signin menu
		// 	"0018-disable-first-run-behaviour.patch",       // Modifies first run to prevent data leakage
		// 	"0019-disable-battery-status-service.patch",    // Disable battry status service as it can be used for tracking
		// 	"0021-disable-rlz.patch",                       // Disable RLZ
		// 	"9000-disable-metrics.patch",                   // Disable metrics
		// 	"9001-disable-profiler.patch",                  // Disable profiler
		// },
	}

	// Not used patches and a description as to why not
	// --------------------------------------------------------------------------
	gNotUsedPatches = map[string][]string{
		gDistros.arch: {
			"chromium-widevine.patch", // Using debian as this one uses a variable
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
	chromiumVer   string // target version of chrome pulled from the PKGBUILD
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
	if chroma.chromiumVer, err = sys.ExtractString(chroma.pkgbuild, exp); err != nil || chroma.chromiumVer == "" {
		err = errors.Errorf("failed to extract the chromium version from the PKGBUILD")
		chroma.logFatal(err)
	}
	if chroma.chromiumVer, err = sys.ExtractString(chroma.pkgbuild, exp); err != nil || chroma.chromiumVer == "" {
		err = errors.Errorf("failed to extract the chromium version from the VERSION file")
		chroma.logFatal(err)
	}

	// Boiler plate for all commands
	// ---------------------------------------------------------------------------------------------
	chroma.printf("Chromium Ver:    %s\n", chroma.chromiumVer)
	chroma.printf("PKBUILD Path:    %s\n", chroma.pkgbuild)
	chroma.println()
	return
}

// validate the chromium version before we continue
func (chroma *Chroma) validateChromiumVersion() (err error) {
	chromium := ""
	exp := `(?m)^pkgver=(.*)$`
	if chroma.chromiumVer, err = sys.ExtractString(chroma.pkgbuild, exp); err != nil || chroma.chromiumVer == "" {
		err = errors.Errorf("failed to extract the chromium version from the PKGBUILD")
		chroma.logFatal(err)
	}
	if chromium, err = sys.ExtractString(chroma.pkgbuild, exp); err != nil || chroma.chromiumVer == "" {
		err = errors.Errorf("failed to extract the chromium version from the VERSION file")
		chroma.logFatal(err)
	}

	// Validate the versions are the same
	if chroma.chromiumVer != chromium {
		err = errors.Errorf("target chromium version in VERSION file is not the same as the PKGBUILD version")
		chroma.logFatal(err)
	}
	return
}
