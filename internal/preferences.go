package internal

import (
    "fmt"
    "os"
    "os/exec"

    log "github.com/sirupsen/logrus"
)

// SavePreferences Save the current preferences (located at $HOME/Library/Preferences) to the /tmp dir.
func SavePreferences(dirName string) (tmpDir string) {
    tmpDir = "/tmp/" + dirName
    cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r $HOME/Library/Preferences/. %s", tmpDir))
    if out, err := cmd.CombinedOutput(); err != nil {
        log.Fatal(string(out))
    }
    return tmpDir
}

// MovePreferences Move the preferences dir
func MovePreferences(sourceDir string, targetDir string) {
    cmd := exec.Command("rm", "-rf", targetDir)
    if out, err := cmd.CombinedOutput(); err != nil {
        log.Fatal(string(out))
    }
    cmd = exec.Command("mv", sourceDir, targetDir)
    if out, err := cmd.CombinedOutput(); err != nil {
        log.Fatal(string(out))
    }
}

// WritePreferencesToFile Write preferences to a file.
func WritePreferencesToFile(fileName string, settings map[string]string) {
    // overwrite if file already exists
    file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    if err != nil {
        log.Fatal(err)
    }
    log.Infof("Writing preferences to %s", fileName)
    for _, command := range settings {
        file.WriteString(command + "\n")
    }
    if err := file.Close(); err != nil {
        log.Fatal(err)
    }
}
