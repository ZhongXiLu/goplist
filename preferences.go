package main

import (
    "fmt"
    "log"
    "os/exec"
)

// Save the current preferences (located at $HOME/Library/Preferences) to the /tmp dir.
func savePreferences(dirName string) (tmpDir string) {
    tmpDir = "/tmp/" + dirName
    cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r $HOME/Library/Preferences/. %s", tmpDir))
    if out, err := cmd.CombinedOutput(); err != nil {
        log.Fatalf(string(out))
    }
    return tmpDir
}

// Move the preferences dir
func movePreferences(sourceDir string, targetDir string) {
    cmd := exec.Command("rm", "-rf", targetDir)
    if out, err := cmd.CombinedOutput(); err != nil {
        log.Fatalf(string(out))
    }
    cmd = exec.Command("mv", sourceDir, targetDir)
    if out, err := cmd.CombinedOutput(); err != nil {
        log.Fatalf(string(out))
    }
}
