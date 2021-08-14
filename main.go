package main

import (
    "flag"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/zhongxilu/plist/internal"
)

var (
    refreshRate int
    verbose     bool
    outputFile  string
)

func main() {
    // CLI arguments
    flag.IntVar(&refreshRate, "r", 3, "Refresh rate (in seconds)")
    flag.BoolVar(&verbose, "v", false, "Verbose mode")
    flag.StringVar(&outputFile, "o", "macos.sh", "Output file with the commands for settings the preferences")
    flag.Parse()

    newSettings := make(map[string]string)

    // Write settings to file on program exit (ctrl + c)
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        internal.WritePreferencesToFile(outputFile, newSettings)
        // TODO: rm -rf prefs dirs
        os.Exit(0)
    }()

    // Diff settings every x seconds and look for changes in the plist files
    // If any are found, convert them to bash commands and save them in `newSettings`
    oldPrefs := internal.SavePreferences("oldPrefs")
    for range time.Tick(time.Duration(refreshRate) * time.Second) {
        newPrefs := internal.SavePreferences("newPrefs")
        internal.DiffPreferences(oldPrefs, newPrefs, newSettings, verbose)
        internal.MovePreferences(newPrefs, oldPrefs)
    }
}
