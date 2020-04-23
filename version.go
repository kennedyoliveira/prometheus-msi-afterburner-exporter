package main

import "fmt"

var (
	branch    string
	sha1      string
	buildDate string
	tag       string
)

type VersionInfo struct {
	Branch    string
	sha1      string
	buildDate string
	tag       string
}

func GetVersionInfo() *VersionInfo {
	return &VersionInfo{
		Branch:    branch,
		sha1:      sha1,
		buildDate: buildDate,
		tag:       tag,
	}
}

func (v *VersionInfo) String() string {
	return fmt.Sprintf("Branch=%s, Sha1=%s, BuildDate=%s, tag=%s", branch, sha1, buildDate, tag)
}
