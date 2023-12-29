package api

import (
	"slices"

	toml "github.com/BurntSushi/toml"
)

type RawConfig struct {
	Predicates []RawPredicate `toml:"sweep"`
}

type RawPredicate struct {
	Predicate []string
	Folder    string
}

type Config struct {
	Predicates []Predicate
}

func ParseConfig(content string) (*Config, error) {
	rawConfig := &RawConfig{}
	_, err := toml.Decode(content, &rawConfig)

	if err != nil {
		return nil, err
	}

	config := &Config{}

	for _, rawPredicate := range rawConfig.Predicates {
		predicate := Predicate{}

		predicate.Folder = rawPredicate.Folder
		predicate.Predicate = func(name string) bool {
			return slices.Contains(rawPredicate.Predicate, name)
		}

		config.Predicates = append(config.Predicates, predicate)
	}

	return config, nil
}
