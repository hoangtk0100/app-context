package appctx

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type appFlagSet struct {
	prefix  string
	flagSet *pflag.FlagSet
}

func newAppFlagSet(prefix string, name string, fs *pflag.FlagSet) *appFlagSet {
	afs := &appFlagSet{
		prefix:  prefix,
		flagSet: fs,
	}

	afs.flagSet.Usage = flagUsages(name, afs)

	return afs
}

func (afs *appFlagSet) GetSampleEnvs() {
	afs.flagSet.VisitAll(func(f *pflag.Flag) {
		if f.Name == "outenv" {
			return
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("## %s (--%s)\n", f.Usage, f.Name))
		sb.WriteString(fmt.Sprintf("#%s=", getEnvName(afs.prefix, f.Name)))

		if !isZeroValue(f, f.DefValue) {
			t := fmt.Sprintf("%T", f.Value)
			if t == "*pflag.stringValue" {
				// put quotes on the value
				sb.WriteString(fmt.Sprintf("%q", f.DefValue))
			} else {
				sb.WriteString(fmt.Sprintf("%v", f.DefValue))
			}
		}

		fmt.Print(sb.String(), "\n\n")
	})
}

// ParseSet parses the given flagset. The specified prefix will be applied to
// the environment variable names.
func (afs *appFlagSet) ParseSet() error {
	var explicit []*pflag.Flag
	var all []*pflag.Flag

	afs.flagSet.Visit(func(f *pflag.Flag) {
		explicit = append(explicit, f)
	})

	var err error
	afs.flagSet.VisitAll(func(f *pflag.Flag) {
		if err != nil {
			return
		}

		all = append(all, f)
		if !contains(explicit, f) {
			name := getEnvName(afs.prefix, f.Name)
			val := viper.GetString(name)
			if !f.Changed && val != "" {
				if ferr := afs.flagSet.Set(f.Name, val); ferr != nil {
					err = fmt.Errorf("failed to set flag %q with value %q", f.Name, val)
				}
			}
		}
	})

	return err
}

func (afs *appFlagSet) Parse(args []string) {
	// Parse flags in command
	err := afs.flagSet.Parse(args)
	if err != nil {
		log.Fatalln("Cannot parse command flags")
	}

	// Set ENV variables to flags value if not passing flags in command
	if err := afs.ParseSet(); err != nil {
		log.Fatalln("Cannot set flags")
	}
}

func isZeroValue(f *pflag.Flag, value string) bool {
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}

	if value == z.Interface().(pflag.Value).String() {
		return true
	}

	switch value {
	case "false":
		return true
	case "":
		return true
	case "0":
		return true
	}

	return false
}

func getEnvName(prefix, name string) string {
	name = strings.Replace(name, ".", "_", -1)
	name = strings.Replace(name, "-", "_", -1)

	if prefix != "" {
		name = fmt.Sprintf("%s_%s", prefix, name)
	}

	return strings.ToUpper(name)
}

func contains(list []*pflag.Flag, f *pflag.Flag) bool {
	for _, i := range list {
		if i == f {
			return true
		}
	}

	return false
}

func flagUsages(name string, afs *appFlagSet) func() {
	return func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)

		afs.flagSet.VisitAll(func(f *pflag.Flag) {
			var sb strings.Builder

			sb.WriteString(fmt.Sprintf("  --%s", f.Name))
			name, usage := pflag.UnquoteUsage(f)
			if len(name) > 0 {
				sb.WriteString(" " + name)
			}

			if sb.Len() <= 4 {
				sb.WriteString("\t")
			} else {
				sb.WriteString("\n    \t")
			}

			sb.WriteString(usage)

			if !isZeroValue(f, f.DefValue) {
				t := fmt.Sprintf("%T", f.Value)
				if t == "*pflag.stringValue" {
					sb.WriteString(fmt.Sprintf(" (default %q)", f.DefValue))
				} else {
					sb.WriteString(fmt.Sprintf(" (default %v)", f.DefValue))
				}
			}

			sb.WriteString(fmt.Sprintf(" [$%s]", getEnvName(afs.prefix, f.Name)))
			_, _ = fmt.Fprint(os.Stderr, sb.String(), "\n")
		})
	}
}
