package main

import (
    "fmt"
    "log"
    "os/exec"
    "regexp"
    "strings"
)

// Plist files we should ignore
var blackList = map[string]struct{}{
    "ByHost":                            {},
    "com.apple.systempreferences.plist": {},
    "com.apple.spaces.plist":            {},
    "jetbrains.jetprofile.asset.plist":  {},
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
            if verbose {
                log.Println("Could not parse %s", line)
            }
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
                    key = getValueOfTag(lines[index-2])
                } else {
                    // name of key is located two lines before
                    //         <key>autohide</key>
                    // -       <false/>
                    // +       <true/>
                    key = getValueOfTag(lines[index-2])
                }
                value := getValueOfTag(diff)

                if key != "" && value != "" {
                    originalPlist := fmt.Sprintf("$HOME/Library/Preferences/%s", oldPlist[14:])
                    command := fmt.Sprintf("defaults write %s %s %s", originalPlist, key, value)

                    if verbose {
                        fmt.Println("")
                        fmt.Println(originalPlist)
                    }
                    fmt.Println(command)
                    if verbose {
                        fmt.Println(string(out))
                    }

                    newSettings[key] = command
                }
            }
        }
    }
}

// Diff two versions of preferences.
func diffPreferences(oldPrefs string, newPrefs string, newSettings map[string]string) {
    cmd := exec.Command("diff", "-ur", oldPrefs, newPrefs)
    out, err := cmd.CombinedOutput()
    if err != nil {
        // exit code 1 = diffs
        for _, diff := range strings.Split(string(out), "\n") {
            if strings.HasPrefix(diff, "Binary files") {
                fields := strings.Fields(diff)
                // Remove the prefix "/tmp/oldPrefs/" and subdirs
                if _, ok := blackList[strings.Split(fields[2][14:], "/")[0]]; !ok {
                    diffPlistFile(fields[2], fields[4], newSettings)
                }
            }
        }
    }
}
