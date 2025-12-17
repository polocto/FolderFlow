
import "io/fs"

type ExampleStartegy struct {
	Custom string `yaml:"examplestartegy"`
}

func (s *ExampleStartegy) FinalDirPath(srcDir, destDir, filePath string, info fs.FileInfo) (string, error) {
	// Example: Custom file operation logic
	return "", nil
}

func (s *ExampleStartegy) Selector() string {
	return "ExampleStartegy"
}

// LoadConfig allows setting configuration for the strategy.
func (s *ExampleStartegy) LoadConfig(config map[string]interface{}) error {
	// No configuration needed for ExampleStartegy
	return nil
}

func init() {
	RegisterStrategy("ExampleStartegy", func() Strategy {
		return &ExampleStartegy{}
	})
}
