package strategy

import (
	"log/slog"

	"github.com/mitchellh/mapstructure"
)

type StrategyYAML struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config,omitempty"`
}

func (sy *StrategyYAML) ToStrategy() (Strategy, error) {
	slog.Debug("Loading strategy config")
	strategy, err := NewStrategy(sy.Name)
	if err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(sy.Config, strategy); err != nil {
		return nil, err
	}
	return strategy, nil
}
