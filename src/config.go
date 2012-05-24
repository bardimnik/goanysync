package main

import (
    "github.com/wor/goconfig/config"
    "strings"
    "fmt"
    "log"
)

// configOptions to be read from the config file.
type ConfigOptions struct {
    tmpfsPath string
    syncPaths []string
    syncerBin string
}

func (self *ConfigOptions) Print() {
    const indent string = "  "
    fmt.Println("Config options:")
    fmt.Println(indent, "TMPFS:", self.tmpfsPath)
    fmt.Println(indent, "RSYNC_BIN:", self.syncerBin)
    fmt.Println(indent, "WHATTOSYNC:")
    for i, v := range self.syncPaths {
        fmt.Printf("%s%s %d: %s\n", indent, indent, i, v)
    }
}

// readConfigFile reads config file and checks that necessary information was
// given. After this it returns the read options in configOptions struct.
func ReadConfigFile(cfp string) (copts *ConfigOptions, err error) {
    const configError string = "Config error: "
    var c *config.Config
    c, err = config.Read(cfp, "# ", "=", true, true)
    if (err != nil) {
        return
    }

    // Read the config file
    tmpfsPath, _ := c.String("DEFAULT", "TMPFS")
    syncerBin, _ := c.String("DEFAULT", "RSYNC_BIN")
    syncPaths, _ := c.String("DEFAULT", "WHATTOSYNC")

    tmpfsPath = strings.TrimSpace(tmpfsPath)
    syncerBin = strings.TrimSpace(syncerBin)
    syncPaths = strings.TrimSpace(syncPaths)

    // Check that given options are valid to some degree
    if len(tmpfsPath) < 1 {
        log.Fatalln(configError, "Empty TMPFS path defined.")
    }
    if len(syncPaths) < 1 {
        log.Fatalln(configError, "Empty WHATTOSYNC paths defined.")
    }
    if len(syncerBin) < 1 {
        // TODO: only do this if rsync can be found from PATH
        syncerBin = "rsync"
    }

    // TODO: check that tmpfsPath is found and writable
    // TODO: check that syncerBin is found and executable

    // Parse WHATTOSYNC comma separated list of paths
    // XXX: if path names contain commas then though luck for now
    fieldFunc := func(r rune) bool {
        return r == ','
    }
    paths := strings.FieldsFunc(syncPaths, fieldFunc)
    if len(paths) < 1 {
        log.Fatalln(configError, "Empty WHATTOSYNC paths defined.")
    }
    for i, v := range paths {
        paths[i] = strings.TrimSpace(v)
    }

    copts = &ConfigOptions{tmpfsPath, paths, syncerBin}
    return
}
