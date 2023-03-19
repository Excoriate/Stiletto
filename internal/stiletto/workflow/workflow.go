package workflow

type WorkflowFile struct {
	Name      string
	Extension string
	Path      string
	Directory string
}

type WorkflowSchemaV1 struct {
	Pipeline string `yaml:"pipeline"`
	Global   struct {
		Workdir string   `yaml:"workdir"`
		Env     []string `yaml:"env"`
	} `yaml:"global"`
	Tasks []struct {
		Name   string   `yaml:"name"`
		Action string   `yaml:"action"`
		With   []string `yaml:"with"`
		Env    []string `yaml:"env,omitempty"`
	} `yaml:"tasks"`
}
