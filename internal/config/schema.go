package config

type Config struct {
	Version   string              `yaml:"version"`
	Info      Info                `yaml:"info"`
	Vars      map[string]string   `yaml:"vars"`
	Instances map[string]Instance `yaml:"instances"`
	Flows     []Flow              `yaml:"flows"`
}

type Info struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Instance struct {
	Driver string         `yaml:"driver"`
	Config map[string]any `yaml:"config"`
}

type Flow struct {
	Name    string   `yaml:"name"`
	Actions []Action `yaml:"actions"`
}

type Action struct {
	ID          string         `yaml:"id"`
	Description string         `yaml:"description"`
	Use         string         `yaml:"use"`
	Timeout     int            `yaml:"timeout"` // ms, 0 = no timeout
	RunWith     map[string]any `yaml:"run_with"`
	MutateOn    MutateOn       `yaml:"mutate_on"`
}

type MutateOn struct {
	Done string `yaml:"done"`
	Fail string `yaml:"fail"`
}
