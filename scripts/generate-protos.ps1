function Test-Command {
    param($Command)
    try { Get-Command $Command -ErrorAction Stop | Out-Null; return $true }
    catch { return $false }
}

if (-not (Test-Command "protoc")) {
    Write-Error "protoc is not installed"
    exit 1
}

if (-not (Test-Command "go")) {
    Write-Error "go is not installed"
    exit 1
}

# Ensure Go proto plugins are installed
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

$SCRIPT_DIR = $PSScriptRoot
$PROJECT_ROOT = (Get-Item $SCRIPT_DIR).Parent.FullName

$SERVICES = @("auth", "wizard", "mana", "spell", "realm")

foreach ($SERVICE in $SERVICES) {
    $PROTO_PATH = Join-Path $PROJECT_ROOT "proto\$SERVICE"
    $PROTO_FILE = Join-Path $PROTO_PATH "$SERVICE.proto"
    
    if (Test-Path $PROTO_FILE) {
        Write-Host "Generating protobuf and gRPC code for $SERVICE..."
        
        $PROTO_REL_PATH = "proto\$SERVICE\$SERVICE.proto"
        
        protoc `
            --proto_path="$PROJECT_ROOT\proto" `
            --go_out="$PROJECT_ROOT" `
            --go_opt=paths=source_relative `
            --go-grpc_out="$PROJECT_ROOT" `
            --go-grpc_opt=paths=source_relative `
            "$SERVICE\$SERVICE.proto"
            
        Write-Host "✓ Generated $SERVICE proto files"
    }
    else {
        Write-Host "⚠ Proto file not found: $PROTO_FILE"
    }
}

Write-Host "✨ Proto generation complete!"