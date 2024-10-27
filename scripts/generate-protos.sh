#!/bin/bash

set -e

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "protoc is not installed. Please install Protocol Buffers compiler first."
    exit 1
fi

# Check if required Go plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Directory setup
PROJECT_ROOT=$(git rev-parse --show-toplevel || pwd)
PROTO_DIR="${PROJECT_ROOT}/proto"
GO_OUT_DIR="${PROJECT_ROOT}"

# Function to generate proto files for a specific service
generate_proto() {
    local service=$1
    local proto_file="${PROTO_DIR}/${service}/${service}.proto"
    local out_dir="${GO_OUT_DIR}/proto/${service}"
    
    if [ -f "$proto_file" ]; then
        echo "Generating protobuf and gRPC code for ${service}..."
        
        # Create output directory if it doesn't exist
        mkdir -p "${out_dir}"
        
        # Generate protobuf code
        protoc \
            --proto_path="${PROTO_DIR}" \
            --go_out="${GO_OUT_DIR}" \
            --go_opt=paths=source_relative \
            --go-grpc_out="${GO_OUT_DIR}" \
            --go-grpc_opt=paths=source_relative \
            "${proto_file}"
            
        echo "✓ Generated ${service} proto files"
    else
        echo "⚠ Proto file not found for ${service}"
    fi
}

# Generate for all services
SERVICES=("auth" "wizard" "mana" "spell" "realm")

echo "🚀 Starting proto generation..."

for service in "${SERVICES[@]}"; do
    generate_proto "$service"
done

echo "✨ Proto generation complete!"