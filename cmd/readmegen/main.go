package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"example.com/spoke-tool/api/types"
	"example.com/spoke-tool/cmd/shared"
	"example.com/spoke-tool/internal/config"
	"example.com/spoke-tool/internal/model"
	"example.com/spoke-tool/internal/test"
)

var (
	// CLI flags
	configPath  = flag.String("config", "config.yaml", "Path to config file")
	projectPath = flag.String("path", ".", "Path to project root")
	watch       = flag.Bool("watch", false, "Watch for changes")
	force       = flag.Bool("force", false, "Force regenerate all tests")
	runTests    = flag.Bool("run", true, "Run tests after generation")
	coverage    = flag.Bool("coverage", false, "Check coverage after tests")
	threshold   = flag.Float64("threshold", 80.0, "Coverage threshold percentage")
	verbose     = flag.Bool("verbose", false, "Verbose output")
	version     = flag.Bool("version", false, "Show version")
	logLevel    = flag.String("log-level", "info", "Log level (error, warn, info, debug, trace)")
	timeout     = flag.Duration("timeout", 10*time.Minute, "Timeout for operations")
	language    = flag.String("lang", "", "Specific language to target (go, nodejs, python)")

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
	log.Printf("Starting testgen %s (log level: %s)", versionInfo.Version, logLevel)

	// Load config
	cfg, err := loadTestGenConfig(flags.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override with CLI flags
	if flags.ProjectPath != "." {
		cfg.ProjectRoot = flags.ProjectPath
	}
	if *threshold > 0 {
		cfg.TestSpoke.CoverageThreshold = *threshold
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
					modelType, getModelDisplayName(modelType, cfg))
			}
		}
	}

	// Parse target language if specified
	targetLang := getTargetLanguage(*language)

	// Initialize test generator
	testGen := test.NewGenerator(test.GeneratorConfig{
		ModelClient:       modelClient,
		ProjectRoot:       cfg.ProjectRoot,
		AutoRunTests:      *runTests,
		CheckCoverage:     *coverage,
		CoverageThreshold: cfg.TestSpoke.CoverageThreshold,
		TargetLanguage:    targetLang,
		Verbose:           flags.Verbose,
	})

	// Track statistics
	stats := &shared.Stats{}
	startTime := time.Now()

	// Run in watch mode or one-shot
	if flags.Watch {
		log.Printf("Starting test generator in watch mode for %s", cfg.ProjectRoot)
		runWatch(ctx, testGen, cfg, stats)
	} else {
		log.Printf("Generating tests for %s", cfg.ProjectRoot)
		result, err := runOnce(ctx, testGen, cfg, flags.Force, stats, targetLang)
		if err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}

		stats.TotalDuration = time.Since(startTime)

		// Print stats if verbose
		if flags.Verbose {
			printStats(stats)
			printTestSummary(result)
		}

		log.Printf("Test generation complete: %s", result.Message)

		// Exit with appropriate code
		if !result.Success {
			os.Exit(int(result.ExitCode))
		}
	}

	return nil
}

// loadConfig loads and validates the configuration
func loadTestGenConfig(path string) (*types.Config, error) {
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
func getModelDisplayName(modelType model.ModelType, cfg *types.Config) string {
	switch modelType {
	case model.CodeLLamaEncoder:
		return cfg.Models.Encoder
	case model.Gemma2B:
		return cfg.Models.Fast
	case model.CodeLLamaDecoder:
		return cfg.Models.Decoder
	default:
		return "unknown"
	}
}

// getTargetLanguage parses the language flag
func getTargetLanguage(langFlag string) types.Language {
	if langFlag == "" {
		return ""
	}
	switch langFlag {
	case "go":
		return types.Go
	case "nodejs", "js":
		return types.NodeJS
	case "python", "py":
		return types.Python
	default:
		log.Printf("Warning: Unknown language %s, will auto-detect", langFlag)
		return ""
	}
}

// runOnce performs a single test generation
func runOnce(ctx context.Context, gen *test.Generator, cfg *types.Config, force bool, stats *shared.Stats, targetLang types.Language) (*test.GenerationResult, error) {
	// Analyze project for functions without tests
	analysisStart := time.Now()
	analysis, err := gen.AnalyzeProject(ctx)
	stats.AnalysisDuration = time.Since(analysisStart)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze project: %w", err)
	}
	stats.FilesAnalyzed = len(analysis.Files)
	stats.FunctionsFound = len(analysis.Functions)

	// Find functions that need tests
	needsTests := gen.FindUntestedFunctions(analysis)
	if len(needsTests) == 0 && !force {
		log.Println("All functions have tests, nothing to generate")
		return &test.GenerationResult{
			Success:         true,
			Message:         "All functions have tests",
			TestsGenerated:  0,
			FunctionsTested: 0,
		}, nil
	}

	log.Printf("Found %d functions needing tests", len(needsTests))

	// Generate tests
	genStart := time.Now()
	result, err := gen.GenerateTests(ctx, analysis, needsTests, force)
	stats.ModelDuration = time.Since(genStart)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tests: %w", err)
	}
	stats.TestsGenerated = len(result.GeneratedTests)
	stats.ModelsQueried = result.ModelsQueried

	// Write test files
	writeStart := time.Now()
	if err := gen.WriteTestFiles(ctx, result); err != nil {
		stats.WriteDuration = time.Since(writeStart)
		return nil, fmt.Errorf("failed to write test files: %w", err)
	}
	stats.WriteDuration = time.Since(writeStart)

	// Run tests if requested
	if gen.Config().AutoRunTests {
		log.Println("Running tests...")
		testResults, err := gen.RunTests(ctx, result)
		if err != nil {
			log.Printf("Warning: Tests failed to run: %v", err)
		} else {
			result.TestResults = testResults

			// Check if any tests failed
			if testResults.Failed > 0 {
				log.Printf("⚠️  %d tests failed", testResults.Failed)
				result.Success = false
				result.Message = fmt.Sprintf("%d tests failed", testResults.Failed)
				result.ExitCode = shared.ExitCodeGenerationError

				// Analyze failures by iterating through Results (not Failures)
				for _, testResult := range testResults.Results {
					if testResult.Status == types.TestStatusFailed {
						analysis, err := gen.AnalyzeFailure(ctx, &testResult)
						if err != nil {
							log.Printf("Failed to analyze failure: %v", err)
						} else {
							log.Printf("Failure analysis for %s:\n%s", testResult.Name, analysis)
						}
					}
				}
			} else {
				log.Printf("✅ All %d tests passed", testResults.Passed)
			}

			// Check coverage if requested
			if *coverage && testResults.Passed == testResults.Total {
				coverage, err := gen.CheckCoverage(ctx)
				if err != nil {
					log.Printf("Warning: Failed to check coverage: %v", err)
				} else {
					result.Coverage = coverage
					log.Printf("Coverage: %.1f%%", coverage.Overall)

					if coverage.Overall < cfg.TestSpoke.CoverageThreshold {
						log.Printf("⚠️  Coverage below threshold: %.1f%% < %.1f%%",
							coverage.Overall, cfg.TestSpoke.CoverageThreshold)
					}
				}
			}
		}
	}

	return result, nil
}

// runWatch watches for changes and regenerates tests
func runWatch(ctx context.Context, gen *test.Generator, cfg *types.Config, stats *shared.Stats) {
	// Create ticker for periodic checks (simplified - in production use fsnotify)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Track last analysis to avoid unnecessary updates
	var lastAnalysis *types.CodeAnalysis
	var lastResult *test.GenerationResult

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
			if lastAnalysis != nil && !testGenHasChanged(lastAnalysis, analysis) {
				continue // No changes
			}

			// Find functions that need tests
			needsTests := gen.FindUntestedFunctions(analysis)
			if len(needsTests) == 0 {
				if lastResult == nil || lastResult.TestsGenerated > 0 {
					log.Println("All functions have tests")
				}
				lastAnalysis = analysis
				continue
			}

			log.Printf("Found %d functions needing tests, regenerating...", len(needsTests))

			// Generate and run tests
			targetLang := gen.Config().TargetLanguage
			result, err := runOnce(ctx, gen, cfg, true, stats, targetLang)
			if err != nil {
				log.Printf("Generation failed: %v", err)
				continue
			}

			lastResult = result
			lastAnalysis = analysis
		}
	}
}

// hasChanged compares two analyses to determine if code has changed
func testGenHasChanged(last *types.CodeAnalysis, current *types.CodeAnalysis) bool {
	if last == nil || current == nil {
		return true
	}
	if len(last.Files) != len(current.Files) {
		return true
	}
	if len(last.Functions) != len(current.Functions) {
		return true
	}
	return false
}

// printStats prints command statistics
func printStats(stats *shared.Stats) {
	log.Printf("=== Statistics ===")
	log.Printf("Files analyzed:    %d", stats.FilesAnalyzed)
	log.Printf("Functions found:   %d", stats.FunctionsFound)
	log.Printf("Tests generated:   %d", stats.TestsGenerated)
	log.Printf("Models queried:    %d", stats.ModelsQueried)
	log.Printf("Analysis time:     %v", stats.AnalysisDuration)
	log.Printf("Model time:        %v", stats.ModelDuration)
	log.Printf("Write time:        %v", stats.WriteDuration)
	log.Printf("Total time:        %v", stats.TotalDuration)
}

// printTestSummary prints test results summary
func printTestSummary(result *test.GenerationResult) {
	if result.TestResults == nil {
		return
	}

	log.Printf("=== Test Results ===")
	log.Printf("Total:  %d", result.TestResults.Total)
	log.Printf("Passed: %d ✅", result.TestResults.Passed)
	log.Printf("Failed: %d ❌", result.TestResults.Failed)
	log.Printf("Skipped: %d ⏭️", result.TestResults.Skipped)

	if result.Coverage != nil {
		log.Printf("Coverage: %.1f%%", result.Coverage.Overall)
	}

	// Show failures by iterating through Results
	if result.TestResults.Failed > 0 {
		log.Printf("=== Failures ===")
		for _, test := range result.TestResults.Results {
			if test.Status == types.TestStatusFailed {
				log.Printf("  %s: %s", test.Name, test.Error)
			}
		}
	}
}
