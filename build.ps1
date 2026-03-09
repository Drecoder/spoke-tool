# From PowerShell (not Git Bash), run:
cd C:\WS\spoke-tool

# Create the build.ps1 file
cat > build.ps1 << 'EOF'
param(
    [string]$Target = "help",
    [string]$Args
)

function Build-Tool {
    param([string]$Name)
    Write-Host "Building $Name..." -ForegroundColor Green
    $version = "dev"
    $commit = "none"
    $date = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"
    go build -ldflags "-X main.Version=$version -X main.Commit=$commit -X main.Date=$date" -o "bin/$Name.exe" "cmd/$Name/main.go"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Built bin/$Name.exe" -ForegroundColor Green
    } else {
        Write-Host "❌ Build failed" -ForegroundColor Red
    }
}

switch ($Target) {
    "build" {
        Build-Tool "readmegen"
        Build-Tool "testgen"
    }
    "build-readmegen" {
        Build-Tool "readmegen"
    }
    "build-testgen" {
        Build-Tool "testgen"
    }
    "test" {
        Write-Host "Running tests..." -ForegroundColor Green
        go test ./...
    }
    "watch-readme" {
        Write-Host "Running readmegen in watch mode..." -ForegroundColor Green
        go run cmd/readmegen/main.go -watch -path .
    }
    "watch-test" {
        Write-Host "Running testgen in watch mode..." -ForegroundColor Green
        go run cmd/testgen/main.go -watch -path .
    }
    "run-readme" {
        Write-Host "Running readmegen..." -ForegroundColor Green
        $runArgs = $Args -split ' '
        go run cmd/readmegen/main.go $runArgs
    }
    "run-test" {
        Write-Host "Running testgen..." -ForegroundColor Green
        $runArgs = $Args -split ' '
        go run cmd/testgen/main.go $runArgs
    }
    default {
        Write-Host @"
Usage: .\build.ps1 [target]

Targets:
  build              Build all tools
  build-readmegen    Build only readmegen
  build-testgen      Build only testgen
  test               Run tests
  watch-readme       Run readmegen in watch mode
  watch-test         Run testgen in watch mode
  run-readme         Run readmegen (add args after)
  run-test           Run testgen (add args after)

Examples:
  .\build.ps1 run-test -path ..\myproject -verbose
  .\build.ps1 run-readme -path ..\myproject -force
  .\build.ps1 build
  .\build.ps1 test
"@ -ForegroundColor Cyan
    }
}
EOF