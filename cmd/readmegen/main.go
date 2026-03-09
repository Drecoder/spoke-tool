cat > cmd/readmegen/main.go << 'EOF'
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/yourusername/spoke-tool/api/types"
	"github.com/yourusername/spoke-tool/cmd/shared"
	"github.com/yourusername/spoke-tool/internal/config"
	"github.com/yourusername/spoke-tool/internal/doc"
	"github.com/yourusername/spoke-tool/internal/model"
)

var (
	// CLI flags
	configPath  = flag.String("config", "config.yaml", "Path to config file")
	projectPath = flag.String("path", ".", "Path to project root")
	watch       = flag.Bool("watch", false, "Watch for changes")
	force       = flag.Bool("force", false, "Force regenerate all docs")
	verbose     = flag.Bool("verbose", false, "Verbose output")
	version     = flag.Bool("version", false, "Show version")
	logLevel    = flag.String("log-level", "info", "Log level (error, warn, info, debug, trace)")
	timeout     = flag.Duration("timeout", 5*time.Minute, "Timeout for operations")
	
	// Version info (set at build time)
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	flag.Parse()

	// Create version info
	versionInfo := shared.VersionInfo{
		Version:   Version,
		Commit:    Commit,
		BuildDate: Date,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	// Show version
	if *version {
		fmt.Println(versionInfo.String())
		os.Exit(int(shared.ExitCodeSuccess))
	}

	// Parse log level
	var level shared.LogLevel
	switch *logLevel {
	case "error":
		level = shared.LogLevelError
	case "warn":
		level = shared.LogLevelWarn
	case "info":
		level = shared.LogLevelInfo
	case "debug":
		level = shared.LogLevelDebug
	case "trace":
		level = shared.LogLevelTrace
	default:
		level = shared.LogLevelInfo
	}

	// Setup logging
	if *verbose {
		level = shared.LogLevelDebug
	}

	// Create flags struct
	flags := &shared.CommandFlags{
		ConfigPath:  *configPath,
		ProjectPath: *projectPath,
		Watch:       *watch,
		Force:       *force,
		Verbose:     *verbose,
		Version:     *version,
		LogLevel:    level,
		Timeout:     *timeout,
	}

	// Run the command
	if err := run(flags, versionInfo); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(flags *shared.CommandFlags, versionInfo shared.VersionInfo) error {
	// Setup logging based on level
	if flags.Verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	logLevel := flags.LogLevel.String()
	log.Printf("Starting readmegen %s (log level: %s)", versionInfo.Version, logLevel)

	// Load config
	cfg, err := loadConfig(flags.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override with CLI flags
	if flags.ProjectPath != "." {
		cfg.ProjectRoot = flags.ProjectPath
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), flags.Timeout)
	defer cancel()

	// Initialize model client
	modelClient, err := initModelClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize model client: %w", err)
	}

	// Check required models
	status, err := modelClient.CheckModels(ctx)
	if err != nil {
		log.Printf("Warning: Failed to check models: %v", err)
	} else {
		for modelType, available := range status {
			if !available {
				log.Printf("Warning: Model %s not available. Run: ollama pull %s", 
					modelType, getModelName(modelType, cfg))
			}
		}
	}

	// Initialize documentation generator
	docGen := doc.NewGenerator(doc.GeneratorConfig{
		ModelClient: modelClient,
		ProjectRoot: cfg.ProjectRoot,
		Sections:    cfg.ReadmeSpoke.Sections,
		Verbose:     flags.Verbose,
	})

	// Track statistics
	stats := &shared.Stats{}
	startTime := time.Now()

	// Run in watch mode or one-shot
	if flags.Watch {
		log.Printf("Starting README generator in watch mode for %s", cfg.ProjectRoot)
		runWatch(ctx, docGen, cfg, stats)
	} else {
		log.Printf("Generating README for %s", cfg.ProjectRoot)
		result, err := runOnce(ctx, docGen, cfg, flags.Force, stats)
		if err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}
		
		stats.TotalDuration = time.Since(startTime)
		
		// Print stats if verbose
		if flags.Verbose {
			printStats(stats)
		}
		
		log.Printf("README generation complete: %s", result.Message)
	}

	return nil
}

// loadConfig loads and validates the configuration
func loadConfig(path string) (*types.Config, error) {
	// If config doesn't exist, use defaults
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Config file %s not found, using defaults", path)
		return config.DefaultConfig(), nil
	}

	cfg, err := config.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}

// initModelClient creates and configures the model client
func initModelClient(cfg *types.Config) (*model.Client, error) {
	modelConfig := model.ClientConfig{
		OllamaHost: "http://localhost:11434",
		Timeout:    30 * time.Second,
		Models: map[model.ModelType]string{
			model.CodeBERT:   cfg.Models.Encoder,
			model.Gemma2B:    cfg.Models.Fast,
			model.DeepSeek7B: cfg.Models.Decoder,
		},
	}

	return model.NewClient(modelConfig)
}

// getModelName returns the actual model name for display
func getModelName(modelType model.ModelType, cfg *types.Config) string {
	switch modelType {
	case model.CodeBERT:
		return cfg.Models.Encoder
	case model.Gemma2B:
		return cfg.Models.Fast
	case model.DeepSeek7B:
		return cfg.Models.Decoder
	default:
		return "unknown"
	}
}

// runOnce performs a single documentation generation
func runOnce(ctx context.Context, gen *doc.Generator, cfg *types.Config, force bool, stats *shared.Stats) (*shared.CommandResult, error) {
	// Analyze project
	analysisStart := time.Now()
	analysis, err := gen.AnalyzeProject(ctx)
	stats.AnalysisDuration = time.Since(analysisStart)
	if err != nil {
		return &shared.CommandResult{
			Success:  false,
			Message:  "Analysis failed",
			ExitCode: shared.ExitCodeAnalysisError,
			Error:    err.Error(),
		}, fmt.Errorf("failed to analyze project: %w", err)
	}
	stats.FilesAnalyzed = len(analysis.Files)
	stats.FunctionsFound = len(analysis.Functions)

	// Check if README needs update
	needsUpdate := force
	if !needsUpdate {
		needsUpdate, err = gen.NeedsUpdate(ctx, analysis)
		if err != nil {
			log.Printf("Warning: Failed to check if update needed: %v", err)
			needsUpdate = true // Update if we can't determine
		}
	}

	if !needsUpdate {
		return &shared.CommandResult{
			Success:  true,
			Message:  "README is up to date, skipping generation",
			ExitCode: shared.ExitCodeSuccess,
		}, nil
	}

	// Generate README sections
	genStart := time.Now()
	sections, err := gen.GenerateSections(ctx, analysis)
	stats.ModelDuration = time.Since(genStart)
	if err != nil {
		return &shared.CommandResult{
			Success:  false,
			Message:  "Section generation failed",
			ExitCode: shared.ExitCodeGenerationError,
			Error:    err.Error(),
		}, fmt.Errorf("failed to generate sections: %w", err)
	}
	stats.DocsGenerated = len(sections)

	// Assemble final README
	writeStart := time.Now()
	readme, err := gen.AssembleReadme(ctx, analysis, sections)
	if err != nil {
		return &shared.CommandResult{
			Success:  false,
			Message:  "README assembly failed",
			ExitCode: shared.ExitCodeGenerationError,
			Error:    err.Error(),
		}, fmt.Errorf("failed to assemble README: %w", err)
	}

	// Write README file
	readmePath := filepath.Join(cfg.ProjectRoot, "README.md")
	if err := os.WriteFile(readmePath, []byte(readme.Content), 0644); err != nil {
		return &shared.CommandResult{
			Success:  false,
			Message:  "Write failed",
			ExitCode: shared.ExitCodeWriteError,
			Error:    err.Error(),
		}, fmt.Errorf("failed to write README: %w", err)
	}
	stats.WriteDuration = time.Since(writeStart)

	log.Printf("Successfully updated %s", readmePath)

	return &shared.CommandResult{
		Success:  true,
		Message:  fmt.Sprintf("Updated %s", readmePath),
		ExitCode: shared.ExitCodeSuccess,
	}, nil
}

// runWatch watches for changes and regenerates documentation
func runWatch(ctx context.Context, gen *doc.Generator, cfg *types.Config, stats *shared.Stats) {
	// Create ticker for periodic checks (simplified - in production use fsnotify)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Track last analysis to avoid unnecessary updates
	var lastAnalysis *types.CodeAnalysis

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Analyze project
			analysis, err := gen.AnalyzeProject(ctx)
			if err != nil {
				log.Printf("Analysis failed: %v", err)
				continue
			}

			// Check if anything changed
			if lastAnalysis != nil && !hasChanged(lastAnalysis, analysis) {
				continue // No changes
			}

			// Check if README needs update
			needsUpdate, err := gen.NeedsUpdate(ctx, analysis)
			if err != nil {
				log.Printf("Failed to check update: %v", err)
				continue
			}

			if needsUpdate {
				log.Println("Changes detected, regenerating README...")
				
				// Generate and write README
				sections, err := gen.GenerateSections(ctx, analysis)
				if err != nil {
					log.Printf("Failed to generate sections: %v", err)
					continue
				}

				readme, err := gen.AssembleReadme(ctx, analysis, sections)
				if err != nil {
					log.Printf("Failed to assemble README: %v", err)
					continue
				}

				readmePath := filepath.Join(cfg.ProjectRoot, "README.md")
				if err := os.WriteFile(readmePath, []byte(readme.Content), 0644); err != nil {
					log.Printf("Failed to write README: %v", err)
					continue
				}

				log.Printf("README updated at %s", time.Now().Format(time.RFC3339))
			}

			lastAnalysis = analysis
		}
	}
}

// hasChanged checks if the code analysis has changed significantly
func hasChanged(old, new *types.CodeAnalysis) bool {
	if old == nil || new == nil {
		return true
	}

	// Check number of files
	if len(old.Files) != len(new.Files) {
		return true
	}

	// Check number of functions
	if len(old.Functions) != len(new.Functions) {
		return true
	}

	// Check timestamps (simplified)
	return old.Timestamp != new.Timestamp
}

// printStats prints command statistics
func printStats(stats *shared.Stats) {
	log.Printf("=== Statistics ===")
	log.Printf("Files analyzed:    %d", stats.FilesAnalyzed)
	log.Printf("Functions found:   %d", stats.FunctionsFound)
	log.Printf("Docs generated:    %d", stats.DocsGenerated)
	log.Printf("Models queried:    %d", stats.ModelsQueried)
	log.Printf("Analysis time:     %v", stats.AnalysisDuration)
	log.Printf("Model time:        %v", stats.ModelDuration)
	log.Printf("Write time:        %v", stats.WriteDuration)
	log.Printf("Total time:        %v", stats.TotalDuration)
}
EOF