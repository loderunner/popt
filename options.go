// Original work, Copyright 2017 Pantomath SAS
// Modified work, Copyright (c) 2019 Charles Francoise
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Option
//
// Use Option to define configuration options for your program. This allows quick and consistent bindings for
// configuration file options, flags and environment variables.
//
// Start by configuring an option for your program.
//
//	var nameOption = popts.Option{
//		Name:    "name",
//		Default: "World",
//		Usage:   "the name of the person you wish to greet",
//		Env:     "HELLO_NAME",
//		Flag:    "name",
//		Short:   "n",
//	}
//
// At init time, add the option to your program.
//
//	func init() {
//		// Add name option
//		if err := popt.AddOption(nameOption, pflag.CommandLine); err != nil {
//			panic(err)
//		}
//	}
//
// When running your executable, bind your option and parse the flags.
//
//	func main() {
//		// Bind env var and flag to viper
//		popt.BindOption(nameOption, pflag.CommandLine)
//
//		// Parse command-line flags
//		pflag.Parse()
//
//		// Read configuration file
//		viper.SetConfigName("hello")
//		viper.AddConfigPath(".")
//		viper.SetConfigType("yaml")
//		viper.ReadInConfig()
//
//		fmt.Println("Hello", viper.GetString(nameOption.Name))
//	}
//
// Example usage:
//
//	$ ./hello
//	Hello World
//
//	$ ./hello -h
//	Usage of ./hello:
//	  -n, --name string   the name of the person you wish to greet (default "World")
//	pflag: help requested
//
//	$ ./hello --name="Steve"
//	Hello Steve
//
//	$ HELLO_NAME="Brooklyn" ./hello
//	Hello Brooklyn
//
//	$ HELLO_NAME="Brooklyn" ./hello --name="Steve"
//	Hello Steve
//
//	$ echo 'name: "Sunshine"' > hello.yaml
//	$ ./hello
//	Hello Sunshine
package popt

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Option describes a confifuration option for the program.
type Option struct {
	Name    string      // The name of the option. Supports viper nesting using dot '.' in names.
	Default interface{} // The default value of the option. Mandatory, as the default value is used to infer the option type.
	Usage   string      // A description of the option.

	Flag  string // The name of the command-line flag.
	Short string // A shorthand for the flag (optional).

	Env string // An environment variable to bind this option to (optional).
}

// AddOption adds an option to the program. If opt.Default is set, it sets the default value in viper. If flags is not
// nil and opt.Flag is set, the option is configuration-only. If opt.Name is empty, and opt.Flag is set, the
// option is flag-only. Use this when setting up flags and configuration options, typically at init time.
func AddOption(opt Option, flags *pflag.FlagSet) error {
	// Set default
	if opt.Name != "" && opt.Default != nil {
		viper.SetDefault(opt.Name, opt.Default)
	}

	// Set flag
	if flags != nil && opt.Flag != "" {
		switch def := opt.Default.(type) {
		case bool:
			flags.BoolP(opt.Flag, opt.Short, def, opt.Usage)
		case int:
			flags.IntP(opt.Flag, opt.Short, def, opt.Usage)
		case float64:
			flags.Float64P(opt.Flag, opt.Short, def, opt.Usage)
		case string:
			flags.StringP(opt.Flag, opt.Short, def, opt.Usage)
		case time.Duration:
			flags.DurationP(opt.Flag, opt.Short, def, opt.Usage)
		default:
			return fmt.Errorf("unsupported option type: %T", def)
		}
	}

	return nil
}

// AddOptions calls AddOption on a list of Options returning the first error it encounters, or nil if none occurred.
func AddOptions(opts []Option, flags *pflag.FlagSet) error {
	for _, o := range opts {
		if err := AddOption(o, flags); err != nil {
			return fmt.Errorf("failed to add option: %s", err)
		}
	}
	return nil
}

// BindOption binds the environment variables and flags to viper. Use this when running the executable, typically at
// the start of a cobra command.
func BindOption(opt Option, flags *pflag.FlagSet) error {
	// Bind environment variable
	if opt.Name != "" && opt.Env != "" {
		if err := viper.BindEnv(opt.Name, opt.Env); err != nil {
			return err
		}
	}

	// Bind flag
	if flags != nil && opt.Flag != "" {
		flag := flags.Lookup(opt.Flag)
		if flag == nil {
			return fmt.Errorf("flag %s not found", opt.Flag)
		}
		if opt.Name != "" {
			viper.BindPFlag(opt.Name, flag)
		}
	}

	return nil
}

// BindOptions calls BindOption on a list of Options returning the first error it encounters, or nil if none occurred.
func BindOptions(opts []Option, flags *pflag.FlagSet) error {
	for _, o := range opts {
		if err := BindOption(o, flags); err != nil {
			return fmt.Errorf("failed to add option: %s", err)
		}
	}
	return nil
}

// AddAndBindOption is a helper function that calls AddOption followed by BindOption. Returns an error if either fails.
// Useful in most cases where there is only one FlagSet.
func AddAndBindOption(opt Option, flags *pflag.FlagSet) error {
	if err := AddOption(opt, flags); err != nil {
		return fmt.Errorf("failed to add option: %s", err)
	}
	if err := BindOption(opt, flags); err != nil {
		return fmt.Errorf("failed to bind option: %s", err)
	}
	return nil
}

// AddAndBindOptions calls AddAndBindOption on a list of Options returning the first error it encounters, or nil if none
// occurred.
func AddAndBindOptions(opts []Option, flags *pflag.FlagSet) error {
	for _, o := range opts {
		if err := AddAndBindOption(o, flags); err != nil {
			return err
		}
	}
	return nil
}
