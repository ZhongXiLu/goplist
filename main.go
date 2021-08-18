package main

import (
    "flag"
    "os"
    "os/signal"
    "syscall"
    "time"

    log "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "github.com/zhongxilu/goplist/internal"
)

var (
    refreshRate int
    verbose     bool
    outputFile  string
    quiet       bool
)

func init() {
    viper.SetDefault("PreferencesDir", "$HOME/Library/Preferences")
    viper.SetDefault("TmpPreferencesDir", "/tmp/prefs")
}

func main() {
    // CLI arguments
    flag.IntVar(&refreshRate, "r", 1, "Refresh rate (in seconds)")
    flag.BoolVar(&verbose, "v", false, "Verbose mode (print some extra information)")
    flag.BoolVar(&quiet, "q", false, "Quiet mode (dont print anything)")
    flag.StringVar(&outputFile, "o", "macos.sh", "Output file with the commands for setting the preferences")
    flag.Parse()

    // Logging
    if verbose {
        log.SetLevel(log.DebugLevel)
    } else if quiet {
        log.SetLevel(log.ErrorLevel)
    }

    prefs := internal.SavePreferences(viper.GetString("TmpPreferencesDir"))
    newSettings := make(map[string]string)

    // Write settings to file on program exit (ctrl + c)
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        internal.WritePreferencesToFile(outputFile, newSettings)
        internal.DeletePreferences(prefs)
        os.Exit(0)
    }()

    // Diff settings every x seconds and look for changes in the plist files
    // If any are found, convert them to bash commands and save them in `newSettings`
    for range time.Tick(time.Duration(refreshRate) * time.Second) {
        internal.DiffPreferences(prefs, viper.GetString("PreferencesDir") + "/", newSettings)
        prefs = internal.SavePreferences(viper.GetString("TmpPreferencesDir"))
    }
}
