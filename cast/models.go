package cast

// Caster is the top level struct representing a caster file
type Caster struct {
	Files   []File   `yaml:"files" json:"files" mapstructure:"files"`
	Folders []Folder `yaml:"folders" json:"folders" mapstructure:"folders"`
}

// File represents a file in the hierarchy
type File struct {
	Name    string `yaml:"name" json:"name" mapstructure:"name"`
	Content string `yaml:"content" json:"content" mapstructure:"content"`
}

// Folder represents a folder in the hierachy
type Folder struct {
	Name    string   `yaml:"name" json:"name" mapstructure:"name"`
	Files   []File   `yaml:"files" json:"files" mapstructure:"files"`
	Folders []Folder `yaml:"folders" json:"folders" mapstructure:"folders"`
}
