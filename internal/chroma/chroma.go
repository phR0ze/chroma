package chroma

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/phR0ze/n"
	"github.com/phR0ze/n/pkg/sys"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Opts allows for passing in complicated arguments
type Opts struct {
	Root    string // path where the chromium PKGBUILD will be located
	Testing bool
}

// CHROMA context
type CHROMA struct {
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
func New(o ...*Opts) (chroma *CHROMA) {
	chroma = &CHROMA{}

	// Configure startup options
	//----------------------------------------------------------------------------------------------
	opts := &Opts{}
	if len(o) > 0 {
		opts = o[0]
	}
	n.SetOnEmpty(&VERSION, "999.999.999")
	var boilerPlate = fmt.Sprintf("Chroma %s [%s (Git %s)]\n", VERSION, BUILDDATE, GITCOMMIT)

	// Configure paths for the application
	//----------------------------------------------------------------------------------------------
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
func (chroma *CHROMA) Execute(args ...string) (err error) {
	if len(args) > 0 {
		os.Args = append([]string{"chroma"}, args...)
	}
	return chroma.cmd.Execute()
}

// configure should be run before anything from the Cobra Run functions
// to configure and validate the environment.
func (chroma *CHROMA) configure() (err error) {

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
	var lines []string
	if lines, err = sys.ReadLines(chroma.pkgbuild); err != nil {
		chroma.logFatal(err)
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "pkgver=") {
			chroma.version = n.A(line).Trim().Split("=").Last().A()
			break
		}
	}

	// Boiler plate for all commands
	// ---------------------------------------------------------------------------------------------
	chroma.printf("Chromium Ver:    %s\n", chroma.version)
	chroma.printf("PKBUILD Path:    %s\n", chroma.pkgbuild)

	return
}
