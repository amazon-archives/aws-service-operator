package version

import (
	"encoding/json"
	"fmt"
)

// Info creates a formattable struct for output
type Info struct {
	Version string
	Commit  string
	Date    string
}

// New will create a pointer to a new version object
func New(version string, commit string, date string) *Info {
	return &Info{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}

// ToJSON converts the Info into a JSON String
func (v *Info) ToJSON() string {
	bytes, _ := json.Marshal(v)
	return string(bytes) + "\n"
}

// ToShortened converts the Info into a JSON String
func (v *Info) ToShortened() string {
	str := fmt.Sprintf("Version: %v\nCommit: %v\nDate: %v\n", v.Version, v.Commit, v.Date)
	return str
}
