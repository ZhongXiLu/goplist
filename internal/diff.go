package internal

import (
    "fmt"
    "os/exec"
    "regexp"
    "strings"

    log "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
)

// Plist files we should ignore
var blackList = map[string]struct{}{
    "ByHost":                            {},
    "jetbrains.jetprofile.asset.plist":  {},
    "com.cisco.anyconnect.gui.plist":    {},
    "com.apple.AddressBook.plist":       {},
    "com.apple.systempreferences.plist": {},
    "com.apple.spaces.plist":            {},
    "com.apple.studentd.plist":          {},
    "com.apple.remindd.plist":           {},
}

// Get the bash command that outputs the equivalent xml representation of a plist file.
func convertPlistToXmlString(plistFile string) string {
    return fmt.Sprintf("<(plutil -convert xml1 -o /dev/stdout %s)", plistFile)
}

// Get the inner value of an XML tag.
// Example: "+       <date>2021-08-14T10:54:26Z</date>" -> 2021-08-14T10:54:26Z
func getValueOfTag(line string) string {
    regex := regexp.MustCompile(`<.*>(.*?)</.*>`)
    matches := regex.FindStringSubmatch(line)
    if len(matches) == 2 {
        return matches[1]
    } else {
        regex := regexp.MustCompile(`<(.*?)/>`) // for collapsed values (e.g. <value/>)
        matches := regex.FindStringSubmatch(line)
        if len(matches) == 2 {
            return matches[1]
        } else {
            log.Debugf("Could not parse %s", line)
        }
    }
    return ""
}

// Get the name of the tag.
func getTagName(line string) string {
    regex := regexp.MustCompile(`<(.*?)>`)
    matches := regex.FindStringSubmatch(line)
    if len(matches) == 2 {
        return matches[1]
    }
    return ""
}

// Diff two plist files.
func diffPlistFile(oldPlist string, newPlist string, newSettings map[string]string) {
    cmd := exec.Command("/bin/bash", "-c",
        fmt.Sprintf("diff -u %s %s", convertPlistToXmlString(oldPlist), convertPlistToXmlString(newPlist)),
    )
    out, err := cmd.CombinedOutput()
    if err != nil {
        // exit code 1 = diffs

        // parse diff output
        lines := strings.Split(string(out), "\n")
        for index, diff := range lines {
            if strings.HasPrefix(diff, "+") && !strings.HasPrefix(diff, "+++") {
                // found diff line
                var key string
                if getTagName(diff) == "date" {
                    continue
                }
                if strings.HasPrefix(lines[index-1], "+") {
                    // name of key is located one line before
                    // +       <key>AppleInterfaceStyle</key>
                    // +       <string>Dark</string>
                    key = getValueOfTag(lines[index-1])
                } else {
                    // name of key is located two lines before
                    //         <key>autohide</key>
                    // -       <false/>
                    // +       <true/>
                    key = getValueOfTag(lines[index-2])
                } // TODO: there might be other cases with nested xml tags...
                value := getValueOfTag(diff)

                if key != "" && value != "" {
                    plistFileName := oldPlist[len(viper.GetString("TmpPreferencesDir"))+1:]
                    originalPlist := fmt.Sprintf("%s/%s", viper.GetString("PreferencesDir"), plistFileName)
                    command := fmt.Sprintf("defaults write %s \"%s\" \"%s\"", originalPlist, key, value)

                    log.Debug("")
                    log.Debug(originalPlist)
                    log.Info(command)
                    log.Debug(string(out))

                    newSettings[key] = command
                }
            }
        }
    }
}

// DiffPreferences Diff two versions of preferences.
func DiffPreferences(oldPrefs string, newPrefs string, newSettings map[string]string) {
    cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("diff -ur %s %s", oldPrefs, newPrefs))
    out, err := cmd.CombinedOutput()
    if err != nil {
        // exit code 1 = diffs
        for _, diff := range strings.Split(string(out), "\n") {
            if strings.HasPrefix(diff, "Binary files") {
                fields := strings.Fields(diff)
                oldPlist := fields[2]
                newPlist := fields[4]
                // Remove the prefix "/tmp/prefs/" and subsequent subdirs
                plistFileName := strings.Split(oldPlist[len(viper.GetString("TmpPreferencesDir"))+1:], "/")[0]
                if _, ok := blackList[plistFileName]; !ok {
                    diffPlistFile(oldPlist, newPlist, newSettings)
                    UpdatePreferences(oldPrefs, newPlist)
                }
            }
        }
    }
}
