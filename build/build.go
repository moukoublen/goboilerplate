package build

// Build time variables.
//
//nolint:gochecknoglobals
var (
	Version     = ""
	Branch      = ""
	Commit      = ""
	CommitShort = ""
	Tag         = ""
)

type Info struct {
	Version     string `json:"version"`
	Branch      string `json:"branch"`
	Commit      string `json:"commit"`
	CommitShort string `json:"commit_short"`
	Tag         string `json:"tag"`
}

func GetInfo() Info {
	return Info{
		Version:     Version,
		Branch:      Branch,
		Commit:      Commit,
		CommitShort: CommitShort,
		Tag:         Tag,
	}
}
