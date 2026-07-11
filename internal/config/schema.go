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
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	Actions     []Action `yaml:"actions"`
}

type Action struct {
	ID          string         `yaml:"id"`
	Description string         `yaml:"description"`
	Instance    string         `yaml:"instance"`
	Timeout     int            `yaml:"timeout"` // ms, 0 = no timeout
	ExecuteWith map[string]any `yaml:"execute_with"`
	MutateOn    MutateOn       `yaml:"mutate_on"`
}

type MutateOn struct {
	Done string `yaml:"done"`
	Fail string `yaml:"fail"`
}
