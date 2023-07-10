package models

// Caster is the top level struct representing a caster file
type Caster struct {
	Files   []File   `yaml:"files,omitempty" json:"files" mapstructure:"files"`
	Folders []Folder `yaml:"folders,omitempty" json:"folders" mapstructure:"folders"`
}

// File represents a file in the hierarchy
type File struct {
	Name    string `yaml:"name,omitempty" json:"name" mapstructure:"name"`
	Content string `yaml:"content,omitempty" json:"content" mapstructure:"content"`
	Ref     string `yaml:"ref,omitempty" json:"ref" mapstructure:"ref"`
}

// Folder represents a folder in the hierachy
type Folder struct {
	Name    string   `yaml:"name,omitempty" json:"name" mapstructure:"name"`
	Files   []File   `yaml:"files,omitempty" json:"files" mapstructure:"files"`
	Folders []Folder `yaml:"folders,omitempty" json:"folders" mapstructure:"folders"`
}

// Variable represents a variable file, key value or environment variable
type Variable struct {
	File  string `yaml:"omitempty"`
	Key   string `yaml:"omitempty"`
	Value string `yaml:"omitempty"`
	Env   string `yaml:"omitempty"`
}
