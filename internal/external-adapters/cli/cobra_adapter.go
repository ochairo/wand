// Package cli provides CLI adapter implementations.
package cli

import (
	"fmt"

	"github.com/ochairo/wand/internal/domain/interfaces"
	"github.com/spf13/cobra"
)

// CobraCLIAdapter wraps Cobra to implement the CLIAdapter interface
type CobraCLIAdapter struct {
	rootCmd                *cobra.Command
	installHandler         interfaces.CommandHandler
	listHandler            interfaces.CommandHandler
	switchHandler          interfaces.CommandHandler
	uninstallHandler       interfaces.CommandHandler
	initHandler            interfaces.CommandHandler
	addHandler             interfaces.CommandHandler
	removeHandler          interfaces.CommandHandler
	wandfileInstallHandler interfaces.CommandHandler
	wandfileCheckHandler   interfaces.CommandHandler
	wandfileDumpHandler    interfaces.CommandHandler
	searchHandler          interfaces.CommandHandler
	infoHandler            interfaces.CommandHandler
	doctorHandler          interfaces.CommandHandler
	updateHandler          interfaces.CommandHandler
	versionHandler         interfaces.CommandHandler
	outdatedHandler        interfaces.CommandHandler
}

// NewCobraCLIAdapter creates a new Cobra CLI adapter
func NewCobraCLIAdapter(
	installHandler interfaces.CommandHandler,
	listHandler interfaces.CommandHandler,
	switchHandler interfaces.CommandHandler,
	uninstallHandler interfaces.CommandHandler,
	initHandler interfaces.CommandHandler,
	addHandler interfaces.CommandHandler,
	removeHandler interfaces.CommandHandler,
	wandfileInstallHandler interfaces.CommandHandler,
	wandfileCheckHandler interfaces.CommandHandler,
	wandfileDumpHandler interfaces.CommandHandler,
	searchHandler interfaces.CommandHandler,
	infoHandler interfaces.CommandHandler,
	doctorHandler interfaces.CommandHandler,
	updateHandler interfaces.CommandHandler,
	versionHandler interfaces.CommandHandler,
	outdatedHandler interfaces.CommandHandler,
) *CobraCLIAdapter {
	adapter := &CobraCLIAdapter{
		installHandler:         installHandler,
		listHandler:            listHandler,
		switchHandler:          switchHandler,
		uninstallHandler:       uninstallHandler,
		initHandler:            initHandler,
		addHandler:             addHandler,
		removeHandler:          removeHandler,
		wandfileInstallHandler: wandfileInstallHandler,
		wandfileCheckHandler:   wandfileCheckHandler,
		wandfileDumpHandler:    wandfileDumpHandler,
		searchHandler:          searchHandler,
		infoHandler:            infoHandler,
		doctorHandler:          doctorHandler,
		updateHandler:          updateHandler,
		versionHandler:         versionHandler,
		outdatedHandler:        outdatedHandler,
	}
	adapter.rootCmd = &cobra.Command{
		Use:   "wand",
		Short: "Wand - A package version manager with shim-based version switching",
		Long: `Wand is a package version manager that allows you to install and manage
multiple versions of CLI tools and GUI applications. It uses shims for
transparent version switching per project.`,
	}

	adapter.setupCommands()
	return adapter
}

// Execute runs the CLI application
func (c *CobraCLIAdapter) Execute() error {
	return c.rootCmd.Execute()
}

// cobraCommandContext wraps a Cobra command to implement CommandContext
type cobraCommandContext struct {
	cmd  *cobra.Command
	args []string
}

func (c *cobraCommandContext) GetStringFlag(name string) (string, error) {
	return c.cmd.Flags().GetString(name)
}

func (c *cobraCommandContext) GetBoolFlag(name string) (bool, error) {
	return c.cmd.Flags().GetBool(name)
}

func (c *cobraCommandContext) GetArgs() []string {
	return c.args
}

func (c *cobraCommandContext) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(c.cmd.OutOrStdout(), format, args...)
}

func (c *cobraCommandContext) PrintError(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(c.cmd.ErrOrStderr(), format, args...)
}

// setupCommands configures all CLI commands
func (c *CobraCLIAdapter) setupCommands() {
	c.rootCmd.AddCommand(c.createInstallCommand())
	c.rootCmd.AddCommand(c.createListCommand())
	c.rootCmd.AddCommand(c.createSwitchCommand())
	c.rootCmd.AddCommand(c.createUninstallCommand())
	c.rootCmd.AddCommand(c.createInitCommand())
	c.rootCmd.AddCommand(c.createAddCommand())
	c.rootCmd.AddCommand(c.createRemoveCommand())
	c.rootCmd.AddCommand(c.createWandfileCommand())
	c.rootCmd.AddCommand(c.createSearchCommand())
	c.rootCmd.AddCommand(c.createInfoCommand())
	c.rootCmd.AddCommand(c.createDoctorCommand())
	c.rootCmd.AddCommand(c.createUpdateCommand())
	c.rootCmd.AddCommand(c.createVersionCommand())
	c.rootCmd.AddCommand(c.createOutdatedCommand())
}

// createInstallCommand creates the install command
func (c *CobraCLIAdapter) createInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <package>[@version]",
		Short: "Install a package",
		Long: `Install a package at the specified version.
If no version is specified, installs the latest version.

Examples:
  wand install node@18.0.0
  wand install terraform
  wand install visual-studio-code`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.installHandler.Handle(ctx)
		},
	}

	cmd.Flags().BoolP("global", "g", false, "Install globally (system-wide)")
	cmd.Flags().Bool("force", false, "Force reinstall if already installed")

	return cmd
}

// createListCommand creates the list command
func (c *CobraCLIAdapter) createListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [package]",
		Short: "List package versions",
		Long: `List available or installed versions of a package.

By default, shows locally installed packages. Use --remote to show versions available from the registry.

Examples:
  wand list                    # List all installed packages
  wand list kubectl            # List installed versions of kubectl (local)
  wand list kubectl --local    # Explicitly list local versions
  wand list kubectl --remote   # List available versions from registry`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.listHandler.Handle(ctx)
		},
	}

	cmd.Flags().Bool("remote", false, "List available versions from remote registry")
	cmd.Flags().Bool("local", false, "List installed versions on local system (default)")
	cmd.MarkFlagsMutuallyExclusive("remote", "local")

	return cmd
}

// createSwitchCommand creates the switch command
func (c *CobraCLIAdapter) createSwitchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch <package>@<version>",
		Short: "Switch to a different version of a package",
		Long: `Switch the active version of an installed package.

The version must already be installed. Use 'wand install' to install new versions.
By default, switches the version for the current project. Use --global to switch system-wide.

Examples:
  wand switch kubectl@1.30.0
  wand switch terraform@1.5.0 --global`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.switchHandler.Handle(ctx)
		},
	}

	cmd.Flags().BoolP("global", "g", false, "Switch version globally (system-wide)")

	return cmd
}

// createInfoCommand creates the info command
func (c *CobraCLIAdapter) createInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <package>",
		Short: "Show information about the active version",
		Long: `Display detailed information about the currently active version of a package.

Shows which version is active globally and/or in the current project.

Examples:
  wand info kubectl
  wand info terraform`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.infoHandler.Handle(ctx)
		},
	}

	return cmd
}

// createUninstallCommand creates the uninstall command
func (c *CobraCLIAdapter) createUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall <package>[@version]",
		Short: "Uninstall a package",
		Long: `Uninstall a package version.
If no version is specified, uninstalls all versions.

Examples:
  wand uninstall node@18.0.0
  wand uninstall terraform`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.uninstallHandler.Handle(ctx)
		},
	}

	return cmd
}

// createInitCommand creates the init command
func (c *CobraCLIAdapter) createInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a .wandrc file in the current directory",
		Long: `Initialize a .wandrc file for managing project-specific package versions.

After initializing, use 'wand add <package>@<version>' to pin versions.

Example:
  wand init`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.initHandler.Handle(ctx)
		},
	}

	return cmd
}

// createAddCommand creates the add command
func (c *CobraCLIAdapter) createAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <package>@<version>",
		Short: "Add a package version to .wandrc",
		Long: `Add a package version to the project's .wandrc file.
The package and version must already be installed.

If .wandrc doesn't exist, it will be created automatically.

Examples:
  wand add nano@8.7.0
  wand add zsh@5.9.0`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.addHandler.Handle(ctx)
		},
	}

	return cmd
}

// createRemoveCommand creates the remove command
func (c *CobraCLIAdapter) createRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <package>",
		Short: "Remove a package from .wandrc",
		Long: `Remove a package from the project's .wandrc file.

This does not uninstall the package, only removes it from the project configuration.
The package will fall back to the global version.

Examples:
  wand remove nano
  wand remove node`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.removeHandler.Handle(ctx)
		},
	}

	return cmd
}

// createWandfileCommand creates the wandfile command with subcommands
func (c *CobraCLIAdapter) createWandfileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wandfile",
		Short: "Manage system configuration with wandfile",
		Long: `Manage declarative system configuration using a wandfile.

A wandfile defines all packages (CLI and GUI) and dotfiles for your system.
Use wandfile to install, check, or export your system configuration.`,
	}

	// Add subcommands
	cmd.AddCommand(c.createWandfileInstallCommand())
	cmd.AddCommand(c.createWandfileCheckCommand())
	cmd.AddCommand(c.createWandfileDumpCommand())

	return cmd
}

// createWandfileInstallCommand creates the wandfile install command
func (c *CobraCLIAdapter) createWandfileInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [wandfile]",
		Short: "Install all packages from wandfile",
		Long: `Install all packages and configure dotfiles from a wandfile.

If no path is specified, looks for './wandfile' in the current directory.

Examples:
  wand wandfile install
  wand wandfile install my-system.wandfile`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.wandfileInstallHandler.Handle(ctx)
		},
	}

	return cmd
}

// createWandfileCheckCommand creates the wandfile check command
func (c *CobraCLIAdapter) createWandfileCheckCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check [wandfile]",
		Short: "Check if all packages in wandfile are installed",
		Long: `Verify that all packages defined in wandfile are currently installed.

If no path is specified, looks for './wandfile' in the current directory.

Examples:
  wand wandfile check
  wand wandfile check my-system.wandfile`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.wandfileCheckHandler.Handle(ctx)
		},
	}

	return cmd
}

// createWandfileDumpCommand creates the wandfile dump command
func (c *CobraCLIAdapter) createWandfileDumpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump [wandfile]",
		Short: "Generate wandfile from installed packages",
		Long: `Create a wandfile from all currently installed packages.

If no path is specified, saves to './wandfile' in the current directory.

Examples:
  wand wandfile dump
  wand wandfile dump my-system.wandfile`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.wandfileDumpHandler.Handle(ctx)
		},
	}

	return cmd
}

// createSearchCommand creates the search command
func (c *CobraCLIAdapter) createSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <term>",
		Short: "Search for packages",
		Long: `Search for packages in the formula repository.
Searches package names, descriptions, and tags.

Examples:
  wand search editor
  wand search nano
  wand search cli`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.searchHandler.Handle(ctx)
		},
	}

	return cmd
}

// createDoctorCommand creates the doctor command
func (c *CobraCLIAdapter) createDoctorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check system health",
		Long: `Run diagnostic checks on your Wand installation.
Verifies:
  - Wand directory structure
  - Package registry integrity
  - Formula repository access
  - Shims directory

Examples:
  wand doctor`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.doctorHandler.Handle(ctx)
		},
	}

	return cmd
}

// createUpdateCommand creates the update command
func (c *CobraCLIAdapter) createUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <package>",
		Short: "Update a package to the latest version",
		Long: `Update an installed package to its latest version.
This will reinstall the package with the latest available version.

Examples:
  wand update nano
  wand update node`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.updateHandler.Handle(ctx)
		},
	}

	return cmd
}

// createVersionCommand creates the version command
func (c *CobraCLIAdapter) createVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long: `Display the version of Wand, along with build information.

Examples:
  wand version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.versionHandler.Handle(ctx)
		},
	}

	return cmd
}

// createOutdatedCommand creates the outdated command
func (c *CobraCLIAdapter) createOutdatedCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outdated",
		Short: "Check for outdated packages",
		Long: `List all installed packages that have newer versions available.
Compares installed versions with the latest versions in the repository.

Examples:
  wand outdated`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := &cobraCommandContext{cmd: cmd, args: args}
			return c.outdatedHandler.Handle(ctx)
		},
	}

	return cmd
}
