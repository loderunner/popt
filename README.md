# popt

[![GoDoc](https://godoc.org/github.com/loderunner/popt?status.svg)](https://godoc.org/github.com/loderunner/popt)

Package popt defines a series of helpers to define options for pflag and viper in one place.

## Option usage

Use Option to define configuration options for your program. This allows quick and consistent bindings for
configuration file options, flags and environment variables.

Start by configuring an option for your program.

```go
var nameOption = popts.Option{
	Name:    "name",
	Default: "World",
	Usage:   "the name of the person you wish to greet",
	Env:     "HELLO_NAME",
	Flag:    "name",
	Short:   "n",
}
```

At init time, add the option to your program.

```go
func init() {
	// Add name option
	if err := popt.AddOption(nameOption, pflag.CommandLine); err != nil {
		panic(err)
	}
}
```

When running your executable, bind your option and parse the flags.

```go
func main() {
	// Bind env var and flag to viper
	popt.BindOption(nameOption, pflag.CommandLine)
	// Parse command-line flags
	pflag.Parse()
	// Read configuration file
	viper.SetConfigName("hello")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	fmt.Println("Hello", viper.GetString(nameOption.Name))
}
```

Example usage:

```
$ ./hello
Hello World

$ ./hello -h
Usage of ./hello:
  -n, --name string   the name of the person you wish to greet (default "World")
pflag: help requested

$ ./hello --name="Steve"
Hello Steve

$ HELLO_NAME="Brooklyn" ./hello
Hello Brooklyn

$ HELLO_NAME="Brooklyn" ./hello --name="Steve"
Hello Steve

$ echo 'name: "Sunshine"' > hello.yaml
$ ./hello
Hello Sunshine
```

You can also define your configuration options in a JSON file, to be loaded at runtime.

`options.json`

```json
[
    {
        "name": "address",
        "default": "localhost",
        "usage": "The address of the remote host",
        "flag": "address",
        "short": "a",
        "env": "HELLO_ADDRESS"
    },
    {
        "name": "port",
        "default": 8080,
        "usage": "The port of the remote host",
        "flag": "port",
        "short": "p",
        "env": "HELLO_PORT"
    },
    {
        "default": false,
        "usage": "Make the operation more talkative",
        "flag": "verbose",
        "short": "v"
    }
]
```

Main

```go
func main() {
	// error handling omitted for brevity
	f, err := os.Open("options.json")
	data, err := ioutil.ReadAll(f)
	var opts []popt.Option
	err = json.Unmarshal(data, &opts)
	err = popt.AddAndBindOptions(opts, pflag.CommandLine)
	pflag.Parse()
}
```

Usage

```
$ ./popt_json -h
Usage of ./popt_json:
-a, --address string   The address of the remote host (default "localhost")
-p, --port string      The port of the remote host (default "8080")
-v, --verbose          Make the operation more talkative
pflag: help requested
```

CAVEAT: Beware that your options defined in JSON will follow JSON typing; in particular, all your numbers will be
float64s.
