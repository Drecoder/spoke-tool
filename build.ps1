param(
    [string]$Target = "help",
    [string]$Arguments
)

function Invoke-BuildTool {
    param(
        [string]$Name,
        [string]$ExtraArgs = ""
    )
    
    Write-Host "Building $Name..." -ForegroundColor Green
    
    # Create bin directory if it doesn't exist
    if (-not (Test-Path "bin")) {
        New-Item -ItemType Directory -Path "bin" -Force | Out-Null
    }
    
    $version = "dev"
    $commit = "none"
    $date = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"
    
    $buildCmd = "go build -ldflags `"-X main.Version=$version -X main.Commit=$commit -X main.Date=$date`" -o `"bin/$Name.exe`" `"cmd/$Name/main.go`" $ExtraArgs"
    
    Write-Host "Running: $buildCmd" -ForegroundColor Gray
    
    $result = Invoke-Expression $buildCmd
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Built bin/$Name.exe" -ForegroundColor Green
        return $true
    } else {
        Write-Host "❌ Build failed with exit code $LASTEXITCODE" -ForegroundColor Red
        Write-Host $result -ForegroundColor Red
        return $false
    }
}

function Invoke-RunTool {
    param(
        [string]$Name,
        [string]$ExtraArgs
    )
    
    Write-Host "Running $Name..." -ForegroundColor Green
    
    $runCmd = "go run cmd/$Name/main.go $ExtraArgs"
    Write-Host "Running: $runCmd" -ForegroundColor Gray
    
    Invoke-Expression $runCmd
}

switch ($Target) {
    "build" {
        Write-Host "=== Building all tools ===" -ForegroundColor Cyan
        $success = $true
        if (-not (Invoke-BuildTool "readmegen")) { $success = $false }
        if (-not (Invoke-BuildTool "testgen")) { $success = $false }
        
        if ($success) {
            Write-Host "`n✅ All tools built successfully!" -ForegroundColor Green
        } else {
            Write-Host "`n❌ Some builds failed" -ForegroundColor Red
            exit 1
        }
    }
    
    "build-readmegen" {
        Invoke-BuildTool "readmegen"
    }
    
    "build-testgen" {
        Invoke-BuildTool "testgen"
    }
    
    "build-readmegen-static" {
        Write-Host "Building readmegen with static linking..." -ForegroundColor Green
        Invoke-BuildTool "readmegen" "-tags netgo -installsuffix netgo"
    }
    
    "build-testgen-static" {
        Write-Host "Building testgen with static linking..." -ForegroundColor Green
        Invoke-BuildTool "testgen" "-tags netgo -installsuffix netgo"
    }
    
    "test" {
        Write-Host "=== Running all tests ===" -ForegroundColor Cyan
        Write-Host "Running: go test ./..." -ForegroundColor Gray
        go test ./...
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "`n✅ All tests passed!" -ForegroundColor Green
        } else {
            Write-Host "`n❌ Tests failed" -ForegroundColor Red
            exit 1
        }
    }
    
    "test-verbose" {
        Write-Host "=== Running all tests (verbose) ===" -ForegroundColor Cyan
        go test -v ./...
    }
    
    "test-coverage" {
        Write-Host "=== Running tests with coverage ===" -ForegroundColor Cyan
        go test -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out -o coverage.html
        Write-Host "✅ Coverage report generated: coverage.html" -ForegroundColor Green
    }
    
    "test-race" {
        Write-Host "=== Running tests with race detector ===" -ForegroundColor Cyan
        go test -race ./...
    }
    
    "bench" {
        Write-Host "=== Running benchmarks ===" -ForegroundColor Cyan
        go test -bench=. ./...
    }
    
    "clean" {
        Write-Host "Cleaning build artifacts..." -ForegroundColor Yellow
        if (Test-Path "bin") {
            Remove-Item -Recurse -Force "bin"
            Write-Host "✅ Removed bin directory" -ForegroundColor Green
        }
        Remove-Item -Recurse -Force "coverage.out" -ErrorAction SilentlyContinue
        Remove-Item -Recurse -Force "coverage.html" -ErrorAction SilentlyContinue
        Write-Host "✅ Clean complete" -ForegroundColor Green
    }
    
    "watch-readme" {
        Write-Host "Running readmegen in watch mode..." -ForegroundColor Green
        $watchArgs = "-watch -path ."
        if ($Arguments) {
            $watchArgs = "$watchArgs $Arguments"
        }
        Invoke-RunTool "readmegen" $watchArgs
    }
    
    "watch-test" {
        Write-Host "Running testgen in watch mode..." -ForegroundColor Green
        $watchArgs = "-watch -path ."
        if ($Arguments) {
            $watchArgs = "$watchArgs $Arguments"
        }
        Invoke-RunTool "testgen" $watchArgs
    }
    
    "run-readme" {
        Write-Host "Running readmegen..." -ForegroundColor Green
        Invoke-RunTool "readmegen" $Arguments
    }
    
    "run-test" {
        Write-Host "Running testgen..." -ForegroundColor Green
        Invoke-RunTool "testgen" $Arguments
    }
    
    "fmt" {
        Write-Host "Formatting code..." -ForegroundColor Cyan
        go fmt ./...
        Write-Host "✅ Code formatted" -ForegroundColor Green
    }
    
    "vet" {
        Write-Host "Running go vet..." -ForegroundColor Cyan
        go vet ./...
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ go vet passed" -ForegroundColor Green
        } else {
            Write-Host "❌ go vet failed" -ForegroundColor Red
            exit 1
        }
    }
    
    "lint" {
        Write-Host "Running golangci-lint..." -ForegroundColor Cyan
        golangci-lint run
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Linting passed" -ForegroundColor Green
        } else {
            Write-Host "❌ Linting failed" -ForegroundColor Red
            exit 1
        }
    }
    
    "tidy" {
        Write-Host "Tidying go modules..." -ForegroundColor Cyan
        go mod tidy
        Write-Host "✅ go.mod tidied" -ForegroundColor Green
    }
    
    "verify" {
        Write-Host "Verifying dependencies..." -ForegroundColor Cyan
        go mod verify
        Write-Host "✅ Dependencies verified" -ForegroundColor Green
    }
    
    "download" {
        Write-Host "Downloading dependencies..." -ForegroundColor Cyan
        go mod download
        Write-Host "✅ Dependencies downloaded" -ForegroundColor Green
    }
    
    "version" {
        Write-Host "spoke-tool build script v1.0.0" -ForegroundColor Cyan
        go version
    }
    
    "help" {
        Write-Host @"
╔════════════════════════════════════════════════════════════════╗
║                    spoke-tool build.ps1                        ║
╚════════════════════════════════════════════════════════════════╝

Usage: .\build.ps1 [target] [-Arguments "args"]

📦 BUILD TARGETS:
  build              Build all tools
  build-readmegen    Build only readmegen
  build-testgen      Build only testgen
  build-readmegen-static  Build readmegen with static linking
  build-testgen-static    Build testgen with static linking
  clean              Remove build artifacts

🧪 TEST TARGETS:
  test               Run all tests
  test-verbose       Run tests with verbose output
  test-coverage      Run tests with coverage report
  test-race          Run tests with race detector
  bench              Run benchmarks

🚀 RUN TARGETS:
  run-readme         Run readmegen (add args after -Arguments)
  run-test           Run testgen (add args after -Arguments)
  watch-readme       Run readmegen in watch mode
  watch-test         Run testgen in watch mode

🔧 UTILITY TARGETS:
  fmt                Format code with go fmt
  vet                Run go vet
  lint               Run golangci-lint
  tidy               Tidy go.mod
  verify             Verify dependencies
  download           Download dependencies
  version            Show version info
  help               Show this help message

📝 EXAMPLES:
  .\build.ps1 build
  .\build.ps1 test-coverage
  .\build.ps1 run-readme -Arguments "-path ..\myproject -verbose"
  .\build.ps1 watch-test
  .\build.ps1 clean

"@ -ForegroundColor Cyan
    }
    
    default {
        Write-Host "Unknown target: $Target" -ForegroundColor Red
        Write-Host "Run '.\build.ps1 help' for usage information" -ForegroundColor Yellow
        exit 1
    }
}