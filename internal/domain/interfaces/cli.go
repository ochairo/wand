package interfaces

// CLIAdapter provides an abstraction over CLI frameworks (Cobra, etc.)
// This allows the domain to be independent of specific CLI implementations
type CLIAdapter interface {
	// Execute runs the CLI application
	Execute() error
}

// CommandContext provides context for command execution
type CommandContext interface {
	// GetStringFlag returns a string flag value
	GetStringFlag(name string) (string, error)

	// GetBoolFlag returns a boolean flag value
	GetBoolFlag(name string) (bool, error)

	// GetArgs returns positional arguments
	GetArgs() []string

	// Printf prints formatted output
	Printf(format string, args ...interface{})

	// PrintError prints formatted error output
	PrintError(format string, args ...interface{})
}

// CommandHandler handles command execution using domain logic
type CommandHandler interface {
	// Handle executes the command with the given context
	Handle(ctx CommandContext) error
}
