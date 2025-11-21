// main is the entry point for the wand CLI application.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	domainadapters "github.com/ochairo/wand/internal/domain-adapters"
	domainorchestrators "github.com/ochairo/wand/internal/domain-orchestrators"
	"github.com/ochairo/wand/internal/domain/services"
	externaladapters "github.com/ochairo/wand/internal/external-adapters"
	"github.com/ochairo/wand/internal/external-adapters/cli"
)

var (
	// Version information (set by build flags)
	Version   = "dev"
	BuildTime = "unknown"
	Commit    = "unknown"
)

func main() {
	// Get paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get home directory: %v\n", err)
		os.Exit(1)
	}
	wandDir := filepath.Join(homeDir, ".wand")
	formulasDir := filepath.Join(wandDir, "formulas")

	// Initialize domain adapters
	fs := domainadapters.NewFileSystemAdapter()
	downloader := domainadapters.NewDownloaderAdapter()
	extractor := domainadapters.NewExtractorAdapter(fs)
	shellExecutor := domainadapters.NewShellExecutorAdapter()

	// Initialize repositories
	registryRepo := domainadapters.NewRegistryRepository(fs, wandDir)
	formulaRepo := domainadapters.NewFormulaRepository(fs, formulasDir)
	wandrcRepo := domainadapters.NewWandRCRepository(fs)
	wandfileRepo := domainadapters.NewWandfileRepository(fs)
	dotfileRepo := domainadapters.NewDotfileRepository(fs, wandDir)

	// Initialize external adapters
	githubClient := externaladapters.NewGitHubAdapter("")

	// Initialize domain services
	versionService := services.NewVersionService(githubClient, formulaRepo)
	shimService := services.NewShimService(registryRepo, wandrcRepo, formulaRepo, fs, wandDir)
	installerService := services.NewInstallerService(
		formulaRepo,
		registryRepo,
		downloader,
		extractor,
		fs,
		shellExecutor,
		versionService,
		wandDir,
		homeDir,
	)
	wandfileService := services.NewWandfileService(
		wandfileRepo,
		registryRepo,
		installerService,
		versionService,
		dotfileRepo,
		fs,
		shellExecutor,
		homeDir,
	)

	// Initialize orchestrators
	installOrchestrator := domainorchestrators.NewInstallOrchestrator(
		installerService,
		shimService,
		versionService,
		formulaRepo,
	)

	// Initialize command handlers
	installHandler := domainorchestrators.NewInstallCommandHandler(
		installOrchestrator,
		registryRepo,
		wandrcRepo,
	)
	listHandler := domainorchestrators.NewListCommandHandler(
		registryRepo,
		wandrcRepo,
		versionService,
	)
	switchHandler := domainorchestrators.NewSwitchCommandHandler(
		registryRepo,
		wandrcRepo,
	)
	infoHandler := domainorchestrators.NewInfoCommandHandler(
		registryRepo,
		wandrcRepo,
	)
	uninstallHandler := domainorchestrators.NewUninstallCommandHandler(
		installOrchestrator,
	)
	initHandler := domainorchestrators.NewInitCommandHandler(
		wandrcRepo,
	)
	addHandler := domainorchestrators.NewAddCommandHandler(
		wandrcRepo,
		registryRepo,
	)
	removeHandler := domainorchestrators.NewRemoveCommandHandler(
		wandrcRepo,
	)
	wandfileInstallHandler := domainorchestrators.NewWandfileInstallCommandHandler(
		wandfileRepo,
		wandfileService,
	)
	wandfileCheckHandler := domainorchestrators.NewWandfileCheckCommandHandler(
		wandfileRepo,
		wandfileService,
	)
	wandfileDumpHandler := domainorchestrators.NewWandfileDumpCommandHandler(
		wandfileRepo,
		wandfileService,
	)
	searchHandler := domainorchestrators.NewSearchCommandHandler(
		formulaRepo,
	)
	doctorHandler := domainorchestrators.NewDoctorCommandHandler(
		registryRepo,
		formulaRepo,
		fs,
		wandDir,
	)
	updateHandler := domainorchestrators.NewUpdateCommandHandler(
		installOrchestrator,
		registryRepo,
	)
	versionHandler := domainorchestrators.NewVersionCommandHandler(
		Version,
		BuildTime,
		Commit,
	)
	outdatedHandler := domainorchestrators.NewOutdatedCommandHandler(
		installOrchestrator,
		registryRepo,
	)

	// Initialize CLI adapter with handlers
	cliAdapter := cli.NewCobraCLIAdapter(
		installHandler,
		listHandler,
		switchHandler,
		uninstallHandler,
		initHandler,
		addHandler,
		removeHandler,
		wandfileInstallHandler,
		wandfileCheckHandler,
		wandfileDumpHandler,
		searchHandler,
		infoHandler,
		doctorHandler,
		updateHandler,
		versionHandler,
		outdatedHandler,
	)

	if err := cliAdapter.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
