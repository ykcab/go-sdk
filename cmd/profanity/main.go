package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/profanity"
	"github.com/blend/go-sdk/ref"
)

var (
	flagRulesFile                      *string
	flagFilesInclude, flagFilesExclude *[]string
	flagDirsInclude, flagDirsExclude   *[]string
	flagVerbose                        *bool
	flagDebug                          *bool
	flagExitFirst                      *bool
)

var (
	_ configutil.Resolver = (*config)(nil)
)

type config struct {
	profanity.Config `yaml:",inline"`
}

// Resolve resolves the config.
func (c *config) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		configutil.SetBool(&c.Verbose, configutil.Bool(flagVerbose), configutil.Bool(c.Verbose), configutil.Bool(ref.Bool(false))),
		configutil.SetBool(&c.FailFast, configutil.Bool(flagDebug), configutil.Bool(c.Debug), configutil.Bool(ref.Bool(false))),
		configutil.SetBool(&c.FailFast, configutil.Bool(flagExitFirst), configutil.Bool(c.FailFast), configutil.Bool(ref.Bool(false))),
		configutil.SetString(&c.RulesFile, configutil.String(*flagRulesFile), configutil.String(c.RulesFile), configutil.String(profanity.DefaultRulesFile)),
		configutil.SetStrings(&c.Files.Include, configutil.Strings(*flagFilesInclude), configutil.Strings(c.Files.Include)),
		configutil.SetStrings(&c.Files.Exclude, configutil.Strings(*flagFilesExclude), configutil.Strings(c.Files.Exclude)),
		configutil.SetStrings(&c.Dirs.Include, configutil.Strings(*flagDirsInclude), configutil.Strings(c.Dirs.Include)),
		configutil.SetStrings(&c.Dirs.Exclude, configutil.Strings(*flagDirsExclude), configutil.Strings(c.Dirs.Exclude)),
	)
}

var configExample = `CONTAINS_EXAMPLE: # id is meant to be a de-duplicating identifier
  description: "please use 'foo.Bar', not a concrete type reference" # description should include remediation steps
  contents: { contains: { include: [ "foo.BarImpl" ] } }

EXCLUDES_EXAMPLE:
  description: "please dont use HerpDerp except in tests"
  contents: { regex: { include: [ "HerpDerp$" ] } }
  files: { exclude: [ "*_test.go" ] }

IMPORTS_EXAMPLE: # you can assert a go AST doesnt contains a given import by glob
  description: "dont include command stuff"
  goImports: { include: [ "github.com/blend/go-sdk/cmd/*" ] }
`

func command() *cobra.Command {
	root := &cobra.Command{
		Use:   "profanity",
		Short: "Enforce profanity rules in a directory tree.",
		Long:  "Enforce profanity rules in a directory tree with inherited rules for each child directory.",
		Example: fmt.Sprintf(`
# Run a basic rules set
profanity --rules=.profanity.yml 

# Run a basic rules set with excluded files by glob
profanity --rules=.profanity.yml --files-exclude="*_test.go"

# Run a basic rules set with included and excluded files by glob
profanity --rules=.profanity.yml --files-include="*.go" --files-exclude="*_test.go"

# An example rule file looks like

""" yaml
%s
"""

For an example rule file (with many more rules), see .profanity.yml in the root of the repo.
`, configExample),
	}

	flagRulesFile = root.Flags().StringP("rules", "r", profanity.DefaultRulesFile, "The rules file to search for in each valid directory")
	flagFilesInclude = root.Flags().StringArray("files-include", nil, "Files to include in glob matching format; can be a csv.")
	flagFilesExclude = root.Flags().StringArray("files-exclude", nil, "Files to exclude in glob matching format; can be a csv.")
	flagDirsInclude = root.Flags().StringArray("dirs-include", nil, "Directories to include in glob matching format; can be a csv.")
	flagDirsExclude = root.Flags().StringArray("dirs-exclude", nil, "Directories to exclude in glob matching format; can be a csv.")
	flagVerbose = root.Flags().BoolP("verbose", "v", false, "If we should show verbose output.")
	flagDebug = root.Flags().BoolP("debug", "d", false, "If we should show debug output.")
	flagExitFirst = root.Flags().Bool("exit-first", false, "If we should fail the run after the first error.")
	return root
}

func main() {
	cmd := command()
	cmd.Run = func(parent *cobra.Command, args []string) {
		var cfg config
		var cfgOptions []configutil.Option
		if flagDebug != nil && *flagDebug {
			cfgOptions = append(cfgOptions, configutil.OptLog(logger.All().WithPath("config")))
		}

		if _, err := configutil.Read(&cfg, cfgOptions...); !configutil.IsIgnored(err) {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		engine := profanity.New(profanity.OptConfig(cfg.Config))
		engine.Stdout = os.Stdout
		engine.Stderr = os.Stderr

		if err := engine.Process(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
			return
		}
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
		return
	}
	os.Exit(0)
}
