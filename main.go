package main

import (
    "time"
)

// RefreshRate How often we want to check the new preferences.
const RefreshRate = 3

func main() {
    newSettings := make(map[string]string)

    oldPrefs := savePreferences("oldPrefs")
    for range time.Tick(time.Second * RefreshRate) {
        newPrefs := savePreferences("newPrefs")
        diffPreferences(oldPrefs, newPrefs, newSettings)
        movePreferences(newPrefs, oldPrefs)
    }

    // TODO: catch ctrl+c and rm -rf prefs dirs
}
