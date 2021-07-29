package root

import (
	"fmt"
	"icotray/assets"

	"gopkg.in/yaml.v3"
)

type Version struct {
	Ref struct {
		Name string `yaml:"name"`
		Slug string `yaml:"slug"`
	}
	Commit struct {
		Sha      string `yaml:"sha"`
		ShortSha string `yaml:"shortSha"`
	}
}

func setVersion() {
	var version Version
	yaml.Unmarshal(assets.Version, &version)

	cmd.Version = fmt.Sprintf("%v-%v", version.Ref.Name, version.Commit.ShortSha)
}
