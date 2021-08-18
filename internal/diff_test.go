package internal

import (
    "os/exec"
    "testing"

    "github.com/spf13/viper"
    "github.com/stretchr/testify/assert"
)

func Test_convertPlistToXmlString(t *testing.T) {
    assert.Equal(t,
        "<(plutil -convert xml1 -o /dev/stdout $HOME/Library/Preferences/.GlobalPreferences.plist)",
        convertPlistToXmlString("$HOME/Library/Preferences/.GlobalPreferences.plist"),
    )
}

func Test_convertPlistToXmlString_emptyPlistFileName(t *testing.T) {
    assert.Equal(t,
        "<(plutil -convert xml1 -o /dev/stdout )",
        convertPlistToXmlString(""),
    )
}

func Test_getValueOfTag_withInnerValue(t *testing.T) {
    assert.Equal(t, "genie", getValueOfTag("<string>genie</string>"))
}

func Test_getValueOfTag_collapsedValue(t *testing.T) {
    assert.Equal(t, "true", getValueOfTag("<true/>"))
}

func Test_getValueOfTag_invalidTag(t *testing.T) {
    assert.Equal(t, "", getValueOfTag("<tagWithNoClosingTag"))
}

func Test_getTagName(t *testing.T) {
    assert.Equal(t, "string", getTagName("<string>genie</string>"))
}

func Test_getTagName_invalidTag(t *testing.T) {
    assert.Equal(t, "", getTagName("<tagWithNoClosingTag"))
}

func Test_diffPlistFile_valueChanged(t *testing.T) {
    viper.SetDefault("PreferencesDir", "$HOME/Library/Preferences")
    viper.SetDefault("TmpPreferencesDir", "testdata/old")
    settingsMap := make(map[string]string)
    diffPlistFile(
        "testdata/old/.GlobalPreferences.plist.valueChanged",
        "testdata/new/.GlobalPreferences.plist.valueChanged",
        settingsMap,
    )
    assert.Contains(t, settingsMap, "_HIHideMenuBar")
    assert.Equal(t,
        "defaults write $HOME/Library/Preferences/.GlobalPreferences.plist.valueChanged \"_HIHideMenuBar\" \"false\"",
        settingsMap["_HIHideMenuBar"],
    )
}

func Test_diffPlistFile_valueAdded(t *testing.T) {
    viper.SetDefault("PreferencesDir", "$HOME/Library/Preferences")
    viper.SetDefault("TmpPreferencesDir", "testdata/old")
    settingsMap := make(map[string]string)
    diffPlistFile(
        "testdata/old/.GlobalPreferences.plist.valueAdded",
        "testdata/new/.GlobalPreferences.plist.valueAdded",
        settingsMap,
    )
    assert.Contains(t, settingsMap, "AppleFontSmoothing")
    assert.Equal(t,
        "defaults write $HOME/Library/Preferences/.GlobalPreferences.plist.valueAdded \"AppleFontSmoothing\" \"AppleInterfaceStyle\"",
        settingsMap["AppleFontSmoothing"],
    )
}

func Test_diffPlistFile_plistDoesNotExist(t *testing.T) {
    settingsMap := make(map[string]string)
    diffPlistFile("nonexistent.plist", "nonexistent.plist", settingsMap)
    assert.Empty(t, settingsMap)
}

func TestDiffPreferences_valueChanged(t *testing.T) {
    viper.SetDefault("PreferencesDir", "$HOME/Library/Preferences")
    viper.SetDefault("TmpPreferencesDir", "testdata/old_tmp")
    // Setup plist dirs so we dont modify initial test data
    cmd := exec.Command("cp", "-r", "testdata/old/.", "testdata/old_tmp")
    if out, err := cmd.CombinedOutput(); err != nil {
        assert.Fail(t, string(out))
    }
    cmd = exec.Command("cp", "-r", "testdata/new/.", "testdata/new_tmp")
    if out, err := cmd.CombinedOutput(); err != nil {
        assert.Fail(t, string(out))
    }

    settingsMap := make(map[string]string)
    DiffPreferences("testdata/old_tmp/", "testdata/new_tmp/", settingsMap)
    assert.Contains(t, settingsMap, "_HIHideMenuBar")
    assert.Equal(t,
        "defaults write $HOME/Library/Preferences/.GlobalPreferences.plist.valueChanged \"_HIHideMenuBar\" \"false\"",
        settingsMap["_HIHideMenuBar"],
    )
    assert.Contains(t, settingsMap, "AppleFontSmoothing")
    assert.Equal(t,
        "defaults write $HOME/Library/Preferences/.GlobalPreferences.plist.valueAdded \"AppleFontSmoothing\" \"AppleInterfaceStyle\"",
        settingsMap["AppleFontSmoothing"],
    )

    // Clean up temporary test plist dirs
    cmd = exec.Command("/bin/bash", "-c", "rm -rf testdata/*_tmp")
    if out, err := cmd.CombinedOutput(); err != nil {
        assert.Fail(t, string(out))
    }
}

func TestDiffPreferences_plistDirDoesNotExist(t *testing.T) {
    settingsMap := make(map[string]string)
    DiffPreferences("nonexistent", "nonexistent", settingsMap)
    assert.Empty(t, settingsMap)
}
