package internal

import (
    "bufio"
    "os"
    "os/exec"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestDeletePreferences(t *testing.T) {
    prefDir := SavePreferences("test")
    DeletePreferences(prefDir)

    cmd := exec.Command("/bin/bash", "-c", "[ -d /tmp/test ]")
    if _, err := cmd.CombinedOutput(); err == nil {
        assert.Fail(t, "/tmp/test is not deleted")
    }

    // Clean up temporary test plist dirs
    cmd = exec.Command("/bin/bash", "-c", "rm -rf /tmp/test")
    if out, err := cmd.CombinedOutput(); err != nil {
        assert.Fail(t, string(out))
    }
}

func TestSavePreferences(t *testing.T) {
    prefDir := SavePreferences("/tmp/test")
    assert.Equal(t, "/tmp/test", prefDir)

    cmd := exec.Command("/bin/bash", "-c", "[ -d /tmp/test ]")
    if _, err := cmd.CombinedOutput(); err != nil {
        assert.Fail(t, "/tmp/test does not exist")
    }

    // Clean up temporary test plist dirs
    cmd = exec.Command("/bin/bash", "-c", "rm -rf /tmp/test")
    if out, err := cmd.CombinedOutput(); err != nil {
        assert.Fail(t, string(out))
    }
}

func TestUpdatePreferences(t *testing.T) {
    // Setup plist dir so we dont modify initial test data
    cmd := exec.Command("cp", "-r", "testdata/old/.", "testdata/old_tmp")
    if out, err := cmd.CombinedOutput(); err != nil {
        assert.Fail(t, string(out))
    }

    settingsMap := make(map[string]string)
    diffPlistFile(
        "testdata/old_tmp/.GlobalPreferences.plist.valueChanged",
        "testdata/new/.GlobalPreferences.plist.valueChanged",
        settingsMap,
    )
    assert.Equal(t, 1, len(settingsMap))

    UpdatePreferences("testdata/old_tmp/", "testdata/new/.GlobalPreferences.plist.valueChanged")

    settingsMap = make(map[string]string)
    diffPlistFile(
        "testdata/old_tmp/.GlobalPreferences.plist.valueChanged",
        "testdata/new/.GlobalPreferences.plist.valueChanged",
        settingsMap,
    )
    assert.Equal(t, 0, len(settingsMap)) // updated

    // Clean up temporary test plist dirs
    cmd = exec.Command("/bin/bash", "-c", "rm -rf testdata/*_tmp")
    if out, err := cmd.CombinedOutput(); err != nil {
        assert.Fail(t, string(out))
    }
}

func TestWritePreferencesToFile(t *testing.T) {
    settingsMap := make(map[string]string)
    settingsMap["_HIHideMenuBar"] =
        "defaults write $HOME/Library/Preferences/.GlobalPreferences.plist.valueChanged \"_HIHideMenuBar\" \"false\""
    settingsMap["AppleFontSmoothing"] =
        "defaults write $HOME/Library/Preferences/.GlobalPreferences.plist.valueAdded \"AppleFontSmoothing\" \"AppleInterfaceStyle\""

    WritePreferencesToFile("macos.sh", settingsMap)

    file, err := os.Open("macos.sh")
    if err != nil {
        assert.Fail(t, "Could not find 'macos.sh'")
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    assert.Contains(t, lines,
        "defaults write $HOME/Library/Preferences/.GlobalPreferences.plist.valueChanged \"_HIHideMenuBar\" \"false\"",
    )
    assert.Contains(t, lines,
        "defaults write $HOME/Library/Preferences/.GlobalPreferences.plist.valueAdded \"AppleFontSmoothing\" \"AppleInterfaceStyle\"",
    )

    // Clean up output file
    cmd := exec.Command("/bin/bash", "-c", "rm macos.sh")
    if out, err := cmd.CombinedOutput(); err != nil {
        assert.Fail(t, string(out))
    }
}
