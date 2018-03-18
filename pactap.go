package main

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"

    "github.com/jessevdk/go-flags"
    "github.com/mitchellh/go-homedir"
)

const VERSION string = "0.0.1"

var opts struct {
    Database bool `short:"D" long:"database" description:"Operate on the package database."`
    Query bool `short:"Q" long:"query" description:"Query the package database."`
    Remove bool `short:"R" long:"remove" description:"Remove packages from system."`
    Sync bool `short:"S" long:"sync" description:"Synchronize packages or groups with remote."`
    Deptest bool `short:"T" long:"deptest" description:"Check dependencies required for a package."`
    Upgrade bool `short:"U" long:"upgrade" description:"Upgrade or add packages via a \"remove-then-add\" process."`
    Files bool `short:"F" long:"files" description:"Query the files database to check which package provides a file."`
    Version bool `short:"V" long:"version" description:"Display version and exit."`

    DatabasePath string `short:"b" long:"dbpath" description:"Specify an alternative database path to \"~/.pactap\"."`
    Arch bool `long:"arch" description:"Specify an alternate architecture."`
    Fancy bool `long:"fancy" description:"Force fancy mode on non-tty systems."` // TODO this is for CLI stuff, will be disabled on detection of non-tty
    Config string `long:"config" description:"Specify an alternate config file."`
    Debug bool `long:"debug" description:"Display debug messages. Use when reporting bugs."`

    // TODO Big options; ones specific to D Q R S T U F
}

func version(){
    fmt.Printf("Pactap v%s\n" +
               "Copyright (C) 2018 Luna Mendes\n", VERSION)
}

func main(){
    if len(os.Args) < 2 {
        os.Args = append(os.Args, "-h")
    }

    _, err := flags.Parse(&opts)
    if err != nil {
        return
    }

    if opts.Version {
        version()
        return
    }

    var bigOpts = make([]bool, 7)
    bigOpts[0] = opts.Database
    bigOpts[1] = opts.Query
    bigOpts[2] = opts.Remove
    bigOpts[3] = opts.Sync
    bigOpts[4] = opts.Deptest
    bigOpts[5] = opts.Upgrade
    bigOpts[6] = opts.Files

    var appliedBigOpts = Filter(bigOpts)

    if len(appliedBigOpts) != 1 {
        fmt.Printf("One big option is required. See -? for a list of options.")
        return
    }

    // TODO: Use the opts for something
    fmt.Printf("Opt Database %t\n", opts.Database)
    fmt.Printf("Opt Query %t\n", opts.Query)
    fmt.Printf("Opt Remove %t\n", opts.Remove)

    homedir, err := homedir.Dir()
    if err != nil {
        panic(err)
    }

    configPath := ".config/pactap/config.toml"

    if runtime.GOOS == "darwin" {
        configPath = "Library/Application Support/pactap/config.toml"
    }

    conf := ReadConfig(filepath.Join(homedir, configPath))
    fmt.Println("raw config:", *conf)

    // Start main program state
    state := &State{
        Config: conf,
        RepoConfig: conf.Repos,
    }

    fmt.Println("raw state", state)

    defer state.Close()

    // TODO: We should start reading our db files, IF ANY
    UpdateRepos(conf)

    // TODO: operate upon args
}


func Filter(vs []bool,) []bool {
    vsf := make([]bool, 0)
    for _, v := range vs {
        if v {
            vsf = append(vsf, v)
        }
    }
    return vsf
}
