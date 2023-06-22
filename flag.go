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

type AppFlagSet struct {
	prefix  string
	flagSet *pflag.FlagSet
}

func newAppFlagSet(prefix string, name string, fs *pflag.FlagSet) *AppFlagSet {
	afs := &AppFlagSet{
		prefix:  prefix,
		flagSet: fs,
	}

	afs.flagSet.Usage = flagUsages(name, afs)

	return afs
}

func (afs *AppFlagSet) GetSampleEnvs() {
	afs.flagSet.VisitAll(func(f *pflag.Flag) {
		if f.Name == "outenv" {
			return
		}

		s := fmt.Sprintf("## %s (--%s)\n", f.Usage, f.Name)
		s += fmt.Sprintf("#%s=", getEnvName(afs.prefix, f.Name))

		if !isZeroValue(f, f.DefValue) {
			t := fmt.Sprintf("%T", f.Value)
			if t == "*pflag.stringValue" {
				// put quotes on the value
				s += fmt.Sprintf("%q", f.DefValue)
			} else {
				s += fmt.Sprintf("%v", f.DefValue)
			}
		}

		fmt.Print(s, "\n\n")
	})
}

// ParseSet parses the given flagset. The specified prefix will be applied to
// the environment variable names.
func (afs *AppFlagSet) ParseSet() error {
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

func (afs *AppFlagSet) Parse(args []string) {
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
		name = prefix + name
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

func flagUsages(name string, afs *AppFlagSet) func() {
	return func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)

		afs.flagSet.VisitAll(func(f *pflag.Flag) {
			s := fmt.Sprintf("  --%s", f.Name)
			name, usage := pflag.UnquoteUsage(f)
			if len(name) > 0 {
				s += " " + name
			}

			if len(s) <= 4 {
				s += "\t"
			} else {
				s += "\n    \t"
			}

			s += usage

			if !isZeroValue(f, f.DefValue) {
				t := fmt.Sprintf("%T", f.Value)
				if t == "*pflag.stringValue" {
					s += fmt.Sprintf(" (default %q)", f.DefValue)
				} else {
					s += fmt.Sprintf(" (default %v)", f.DefValue)
				}
			}

			s += fmt.Sprintf(" [$%s]", getEnvName(afs.prefix, f.Name))
			_, _ = fmt.Fprint(os.Stderr, s, "\n")
		})
	}
}
