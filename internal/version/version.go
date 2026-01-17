package version

import "fmt"

var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

type Info struct {
	Version   string
	Commit    string
	BuildTime string
}

func Info() Info {
	return Info{Version: Version, Commit: Commit, BuildTime: BuildTime}
}

func (i Info) String() string {
	return fmt.Sprintf("%s (commit %s, built %s)", i.Version, i.Commit, i.BuildTime)
}
