package commandselector

// Go doesn't have enums, so we use a type alias and a const block to simulate them
type PackageManager string

const (
	Npm  PackageManager = "npm"
	Yarn PackageManager = "yarn"
	Pnpm PackageManager = "pnpm"
)

type Project struct {
	Manager PackageManager
}

type Package struct {
	Path    string
	IsRoot  bool
	Manager PackageManager
	Json    PackageJson
}

type PackageJson struct {
	Name    string            `json:"name"`
	Scripts map[string]string `json:"scripts"`
	Dir     string
}

type ExecOptions struct {
	WithRunner bool
}

type Command struct {
	PackageName  string
	CommandName  string
	CommandValue string
	PackageDir   string
	ExecOptions  ExecOptions
}
