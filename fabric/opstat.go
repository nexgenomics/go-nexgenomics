package fabric

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const magic_dir = "/tmp/opstat.d"

// GetOperationalStatus is for use by in-agent status monitors.
// It reads the files present in a "magic" directory and returns the
// contents as a flattened map.
// Use PutOperationalStatus to save values so they are visible to this
// function.
func GetOperationalStatus() (map[string]any, error) {
	out := map[string]any{}

	if e := os.MkdirAll(magic_dir, 0o755); e != nil {
		return out, e
	}

	filepath.WalkDir(magic_dir, func(path string, d fs.DirEntry, err error) error {
		if err == nil {
			if !d.IsDir() {
				if strings.HasSuffix(d.Name(), ".json") {
					if data, e := os.ReadFile(path); e == nil {
						var m map[string]any
						if e := json.Unmarshal(data, &m); e == nil {
							nm := strings.TrimSuffix(d.Name(), ".json")
							for k, v := range m {
								out[fmt.Sprintf("%s/%s", nm, k)] = v
							}
						}
					}
				}
			}
		}
		return nil
	})

	return out, nil
}
