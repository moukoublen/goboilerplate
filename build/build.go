package build

// Build time variables
var (
	Version     = ""
	Branch      = ""
	Commit      = ""
	CommitShort = ""
	Tag         = ""
)

type buildInfo struct {
	Version     string `json:"version"`
	Branch      string `json:"branch"`
	Commit      string `json:"commit"`
	CommitShort string `json:"commit_short"`
	Tag         string `json:"tag"`
}

func BuildInfo() buildInfo {
	return buildInfo{
		Version:     Version,
		Branch:      Branch,
		Commit:      Commit,
		CommitShort: CommitShort,
		Tag:         Tag,
	}
}
