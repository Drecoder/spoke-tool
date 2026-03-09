// Create runner
runner := NewRunner(RunnerConfig{
	WorkingDir:      "./myproject",
	Timeout:         2 * time.Minute,
	CollectCoverage: true,
})

// Run specific test files
run, err := runner.RunTests(ctx, []string{"math_test.go"}, types.Go)

// Run all tests
results, err := runner.RunAllTests(ctx)

// Get summary
fmt.Println(runner.GetSummary(run))

// Get failure details
failures := runner.GetFailureDetails(run)
for _, f := range failures {
	fmt.Printf("❌ %s: %s\n", f.TestName, f.ErrorMsg)
}