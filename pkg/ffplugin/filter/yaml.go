package filter

import "github.com/mitchellh/mapstructure"

type FilterYAML struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config,omitempty"`
}

func (fy *FilterYAML) ToFilter() (Filter, error) {
	filter, err := NewFilter(fy.Name)
	if err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(fy.Config, filter); err != nil {
		return nil, err
	}
	return filter, nil
}
