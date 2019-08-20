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
	debian    string
	ungoogled string
}

var (
	gDistros = Distros{"debian", "ungoogled"}

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
		"debian":    "https://salsa.debian.org/chromium-team/chromium/raw/master/debian/patches/series",
		"ungoogled": "https://raw.githubusercontent.com/Eloston/ungoogled-chromium/master/patches/series",
		"inox":      "https://github.com/gcarq/inox-patchset",
		"iridium":   "https://git.iridiumbrowser.de/cgit.cgi/iridium-browser/commit/?h=patchview",
	}

	// Call out patches used and not used and notes
	// Order is significant
	// --------------------------------------------------------------------------
	gPatches = map[string]map[string]bool{

		// Credit to Michael Gilber
		gDistros.debian: {
			"00-manpage.patch":                  true,  // Adds simple doc with link to documentation website
			"01-sandbox.patch":                  false, // Debian specific error message to install chromium-sandbox
			"02-master-preferences.patch":       true,  // Look for master preferences in /etc/chromium/master_preferences
			"03-libcxx.patch":                   true,  // Avoid chromium's embedded C++ library when bootstrapping
			"04-parallel.patch":                 true,  // Respect specified number of parllel jobs when bootstrapping
			"05-gcc_skcms_ice.patch":            true,  // GCC ICE with optimized version
			"06-pffffft-buildfix.patch":         true,  // ??
			"07-skia-aarch64-buildfix.patch":    true,  // ??
			"08-wrong-namespace.patch":          false, // gcc: not using as getting inspector protocol errors
			"09-virtual-destructor.patch":       true,  // gcc: a virtual destructor is called without this patch
			"10-explicit-specialization.patch":  true,  // gcc: fix for gcc explicit specialiazation namespace issue
			"11-macro.patch":                    false, // gcc6: can be ignored as arch linux is using gcc 9
			"12-sizet.patch":                    false, // gcc6: can be ignored as arch linux is using gcc 9
			"13-atomic.patch":                   false, // gcc6: can be ignored as arch linux is using gcc 9
			"14-constexpr.patch":                false, // gcc6: can be ignored as arch linux is using gcc 9
			"15-wtf-hashmap.patch":              false, // gcc6: can be ignored as arch linux is using gcc 9
			"16-lambda-this.patch":              false, // gcc6: can be ignored as arch linux is using gcc 9
			"17-map-insertion.patch":            false, // gcc6: can be ignored as arch linux is using gcc 9
			"18-not-constexpr.patch":            false, // gcc6: can be ignored as arch linux is using gcc 9
			"19-move-required.patch":            false, // gcc6: can be ignored as arch linux is using gcc 9
			"20-use-after-move.patch":           false, // gcc6: can be ignored as arch linux is using gcc 9
			"21-ambiguous-overloads.patch":      false, // gcc6: can be ignored as arch linux is using gcc 9
			"22-ambiguous-initializer.patch":    false, // gcc6: can be ignored as arch linux is using gcc 9
			"23-nullptr-copy-construct.patch":   false, // gcc6: can be ignored as arch linux is using gcc 9
			"24-noexcept-redeclaration.patch":   false, // gcc6: can be ignored as arch linux is using gcc 9
			"25-trivially-constructible.patch":  false, // gcc6: can be ignored as arch linux is using gcc 9
			"26-designated-initializers.patch":  false, // gcc6: can be ignored as arch linux is using gcc 9
			"27-specialization-namespace.patch": false, // gcc6: can be ignored as arch linux is using gcc 9
			"28-mojo.patch":                     true,  // Fixes: fix mojo layout test build error
			"29-public.patch":                   true,  // Fixes: method needs to be public
			"30-ps-print.patch":                 true,  // Fixes: add postscript(ps) printing capabiliy
			"31-as-needed.patch":                true,  // Fixes: some libraries fail to link when '--as-needed' is set
			"32-inspector.patch":                false, // Fixes: not using as getting inspector protocol errors
			"33-gpu-timeout.patch":              true,  // Fixes: increase GPU timeout from 10sec to 20sec
			"34-empty-array.patch":              true,  // Fixes: arraysize macro fails for zero length array and add one char
			"35-safebrowsing.patch":             false, // Fixes: fix signedness error when built with gcc affects safe browsing
			"36-sequence-point.patch":           true,  // Fixes: fix undefined order in which expressions are evaluated
			"37-jumbo-namespace.patch":          true,  // Fixes: jumbo build has trouble with these namespaces
			"38-template-export.patch":          true,  // Fixes: implementation of template function must be in header to be exported
			"39-widevine-revision.patch":        true,  // Fixes: set widevine version as undefined
			"40-widevine-locations.patch":       false, // Fixes: arch linux works fine don't need to try alternative location for widevine
			"41-widevine-buildflag.patch":       true,  // Fixes: enable widevine support
			"42-connection-message.patch":       false, // Fixes: hardly seems important to 'update suggest updating your proxy when network is unreachable'
			"43-unrar.patch":                    true,  // Disable: disable support for browsing rar files
			"44-signin.patch":                   false, // Disable: already covered in the ungoogled patches
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
			"56-third-party-cookies.patch":      false, // Disable: covered by the inox patch 0006-modify-default-prefs.patch
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
			"79-event.patch":                    false, // System: might be causing libeevnt build failure - build using the system libevent library
			"80-ffmpeg.patch":                   false, // System: arch linux PKGBUILD has a system lib call out for this already
			"81-jsoncpp.patch":                  true,  // System: use system jsoncpp
			"82-openjpeg.patch":                 true,  // System: build system using openjpeg
			"83-convertutf.patch":               true,  // System: use ICU for UTF8 conversions (eleminates ConvertUTF embedded code copy)
			"84-icu63.patch":                    false, // System: arch linux has newer icu don't need to maintain compt with 63
		},

		// Credit to github.com/Eloston/ungoogled-chromium
		gDistros.ungoogled: {
			"00-chromium-exclude_unwind_tables.patch":                       true,  // inox: Exclude unwind dumps as stack dumps can be unwound by Crashpad at a later time
			"01-0001-fix-building-without-safebrowsing.patch":               true,  // inox: Fix building with 'safe_browsing_mode=0' set
			"02-0003-disable-autofill-download-manager.patch":               true,  // inox: Disables HTML AutoFill data transmission to Google
			"03-0004-disable-google-url-tracker.patch":                      false, // inox: Disable Google tracking your entered urls, but breaks omnibar search
			"04-0005-disable-default-extensions.patch":                      true,  // inox: Disable CloudPrint, Feedback, WebStore, InAppPayments
			"05-0007-disable-web-resource-service.patch":                    true,  // inox: Disables downloading dynamic configuration from Google for chromium
			"06-0009-disable-google-ipv6-probes.patch":                      true,  // inox: Change IPv6 DNS probes to Google over to k.root-servers.net
			"07-0010-disable-gcm-status-check.patch":                        true,  // inox: Disable Google Cloud-Messaging status probes, GCM allows direct msg to device
			"08-0014-disable-translation-lang-fetch.patch":                  true,  // inox: Disable language fetching from Google when settings are opened the first time
			"09-0015-disable-update-pings.patch":                            true,  // inox: Disable update pings to Google
			"10-0017-disable-new-avatar-menu.patch":                         true,  // inox: Disable Google Avatar signin menu
			"11-0021-disable-rlz.patch":                                     true,  // inox: Disable RLZ
			"12-unrar.patch":                                                false, // debian: already covered by debian
			"13-perfetto.patch":                                             false, // debian: already covered by debian
			"14-safe_browsing-disable-incident-reporting.patch":             true,  // iridium: disable safe browsing incident reporting
			"15-safe_browsing-disable-reporting-of-safebrowsing-over.patch": true,  // iridium: disable safe browsing incident reporting
			"16-all-add-trk-prefixes-to-possibly-evil-connections.patch":    true,  // iridium: block outgoing calls to any google servers
			"17-disable-crash-reporter.patch":                               true,  // ungoogled: disable crash reporting
			"18-disable-google-host-detection.patch":                        true,  // ungoogled: disable detecting Google hosts
			"19-replace-google-search-engine-with-nosearch.patch":           false, // ungoogled: leaving in the google search engine
			"20-disable-signin.patch":                                       true,  // ungoogled: disable browser signin
			"21-disable-translate.patch":                                    true,  // ungoogled: disable browser translate
			"22-disable-untraceable-urls.patch":                             true,  // ungoogled: disable additional outgoing URLs not caught by "trk" scheme
			"23-disable-profile-avatar-downloading.patch":                   true,  // ungoogled: disable downloading profile avatar
			"24-disable-gcm.patch":                                          true,  // ungoogled: disable Google Cloud Messaging
			"25-disable-domain-reliability.patch":                           true,  // ungoogled: disable domain reliability component
			"26-block-trk-and-subdomains.patch":                             true,  // ungoogled: block other outgoing URLs
			"27-fix-building-without-one-click-signin.patch":                true,  // ungoogled: fix building without one click signin
			"28-disable-gaia.patch":                                         true,  // ungoogled: ensure can't be activated even without signing in
			"29-disable-fonts-googleapis-references.patch":                  false, // ungoogled: google fonts are alright, leaving in
			"30-disable-webstore-urls.patch":                                false, // ungoogled: still want access to the webstore so leaving this in
			"31-fix-learn-doubleclick-hsts.patch":                           true,  // ungoogled:
			"32-disable-webrtc-log-uploader.patch":                          true,  // ungoogled: disable webrtc log uploader
			"33-use-local-devtools-files.patch":                             true,  // ungoogled: bundle in dev files rather than download them
			"34-disable-network-time-tracker.patch":                         true,  // ungoogled: disable network time tracker
			"35-disable-mei-preload.patch":                                  true,  // ungoogled: disable mei preload
			"36-fix-building-without-safebrowsing.patch":                    true,  // ungoogled: fix building without safebrowsing
			"37-disable-fetching-field-trials.patch":                        true,  // bromite: disable fetching field trials

			"38-chromium-widevine.patch":                                    false, // ungoogled: already covered by debian
			"39-0006-modify-default-prefs.patch":                            true,  // inox: set sane defaults for preferences
			"40-0008-restore-classic-ntp.patch":                             true,  // inox: the new NTP (New Tag Page) pulls from Google including tracking identifier
			"41-0011-add-duckduckgo-search-engine.patch":                    true,  // inox: add duckduckgo search option
			"42-0013-disable-missing-key-warning.patch":                     true,  // inox: disable missing google api key warning
			"43-0016-chromium-sandbox-pie.patch":                            true,  // inox: hardening the sandbox with Position Independent Code(PIE) against ROP exploits
			"44-0018-disable-first-run-behaviour.patch":                     true,  // inox: disable first run behavior
			"45-0019-disable-battery-status-service.patch":                  true,  // inox: disable battery status service
			"46-parallel.patch":                                             false, // debian: already covered
			"47-ps-print.patch":                                             false, // debian: already covered
			"48-inspector.patch":                                            false, // debian: already covered
			"49-connection-message.patch":                                   false, // debian: already covered
			"50-android.patch":                                              false, // debian: already covered
			"51-fuzzers.patch":                                              false, // debian: already covered
			"52-welcome-page.patch":                                         false, // debian: already covered
			"53-google-api-warning.patch":                                   false, // debian: already covered
			"54-device-notifications.patch":                                 false, // debian: already covered
			"55-initialization.patch":                                       false, // debian: already covered
			"56-net-cert-increase-default-key-length-for-newly-gener.patch": true,  // iridium: increase default key length from 1024 => 2056
			"57-mime_util-force-text-x-suse-ymp-to-be-downloaded.patch":     false, // iridium: force download of ymp files
			"58-prefs-only-keep-cookies-until-exit.patch":                   true,  // iridium: set cookies to only be kept unit exit
			"59-prefs-always-prompt-for-download-directory-by-defaul.patch": true,  // iridium: always prompt for download directory by default
			"60-updater-disable-auto-update.patch":                          false, // iridium: auto update is already turned off for Linux
			"61-Remove-EV-certificates.patch":                               false, // iridium: just cosmetics - skipping
			"62-browser-disable-profile-auto-import-on-first-run.patch":     true,  // iridium: disable auto importing stuff on first run
			"63-add-third-party-ungoogled.patch":                            false, // ungoogled: skipping
			"64-disable-formatting-in-omnibox.patch":                        false, // ungoogled: skipping
			"65-popups-to-tabs.patch":                                       true,  // ungoogled: force pop up windows to end up as a new tab
			"66-add-ipv6-probing-option.patch":                              true,  // ungoogled: disable IPV6 probing
			"67-remove-disable-setuid-sandbox-as-bad-flag.patch":            false, // ungoogled: skipping
			"68-disable-intranet-redirect-detector.patch":                   true,  // ungoogled: disable internet redirect detector, stop extraneous dns requests
			"69-enable-page-saving-on-more-pages.patch":                     true,  // ungoogled: allow saving of more documents rather than just HTTP/HTTPS
			"70-disable-download-quarantine.patch":                          true,  // ungoogled: disable file download quarantine, always available
			"71-fix-building-without-mdns-and-service-discovery.patch":      true,  // ungoogled: fix building without mdns and service discovery
			"72-add-flag-to-stack-tabs.patch":                               false, // ungoogled: skipping
			"73-add-flag-to-configure-extension-downloading.patch":          false, // ungoogled: skipping
			"74-add-flag-for-search-engine-collection.patch":                false, // ungoogled: skipping
			"75-add-flag-to-disable-beforeunload.patch":                     false, // ungoogled: skipping
			"76-add-flag-to-force-punycode-hostnames.patch":                 false, // ungoogled: skipping
			"77-searx.patch":                                                true,  // ungoogled: adding awesome private search option
			"78-disable-webgl-renderer-info.patch":                          true,  // ungoogled: removing webgl data leakage
			"79-add-flag-to-show-avatar-button.patch":                       false, // ungoogled: skipping
			"80-add-suggestions-url-field.patch":                            false, // ungoogled: skipping
			"81-add-flag-to-hide-crashed-bubble.patch":                      false, // ungoogled: skipping
			"82-default-to-https-scheme.patch":                              true,  // ungoogled: default urls without a schema to https
			"83-add-flag-to-scroll-tabs.patch":                              false, // ungoogled: skipping
			"84-enable-paste-and-go-new-tab-button.patch":                   true,  // ungoogled: enable paste and go new tab
			"85-fingerprinting-flags-client-rects-and-measuretext.patch":    false, // bromite: skipping
			"86-flag-max-connections-per-host.patch":                        false, // bromite: skipping
			"87-flag-fingerprinting-canvas-image-data-noise.patch":          false, // bromite: skipping
		},

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
		// },
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
