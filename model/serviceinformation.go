package model

type ServiceInfo struct {
	Name         string   `json:"NAME"`
	Version      string   `json:"VERSION"`
	VersionState string   `json:"VERSION_STATE"`
	VersionName  string   `json:"VERSION_NAME"`
	Author       string   `json:"AUTHOR"`
	Contributor  []string `json:"CONTRIBUTOR"`
}
