#
# Makefile for Phantom - Fast build without heavy downloads
#

GO ?= go
CGO_ENABLED ?= 0
TAGS = -tags go_sqlite

# Auto-detect OS and Architecture
UNAME_S := $(shell uname -s 2>/dev/null || echo Windows)
UNAME_M := $(shell uname -m 2>/dev/null || echo x86_64)

# Detect if running in Termux
IS_TERMUX := $(shell [ -d /data/data/com.termux ] && echo 1 || echo 0)

# Set OS defaults
ifeq ($(UNAME_S),Linux)
	GOOS ?= linux
	EXT =
	ifeq ($(IS_TERMUX),1)
		GOOS = android
	endif
endif

ifeq ($(UNAME_S),Darwin)
	GOOS ?= darwin
	EXT =
endif

ifeq ($(findstring MINGW,$(UNAME_S)),MINGW)
	GOOS ?= windows
	EXT = .exe
endif

ifeq ($(findstring MSYS,$(UNAME_S)),MSYS)
	GOOS ?= windows
	EXT = .exe
endif

ifeq ($(UNAME_S),Windows_NT)
	GOOS ?= windows
	EXT = .exe
endif

# Auto-detect architecture
ifeq ($(UNAME_M),x86_64)
	GOARCH ?= amd64
endif

ifeq ($(UNAME_M),aarch64)
	GOARCH ?= arm64
endif

ifeq ($(UNAME_M),arm64)
	GOARCH ?= arm64
endif

ifeq ($(UNAME_M),armv7l)
	GOARCH ?= arm
endif

ifeq ($(UNAME_M),i686)
	GOARCH ?= 386
endif

# Build flags
LDFLAGS = -ldflags "-s -w \
	-X github.com/cryptdefender3232/phantom/client/command/update.PhantomPublicKey=RWTZPg959v3b7tLG7VzKHRB1/QT+d3c71Uzetfa44qAoX5rH7mGoQTTR \
	-X github.com/cryptdefender3232/phantom/client/assets.DefaultArmoryPublicKey=RWSBpxpRWDrD7Fe+VvRE3c2VEDC2NK80rlNCj+BX0gz44Xw07r6KQD9L \
	-X github.com/cryptdefender3232/phantom/client/assets.DefaultArmoryRepoURL=https://api.github.com/repos/phantomarmory/armory/releases"

#
# Main targets
#
.PHONY: default
default: clean
	@echo "🔍 Detected: $(GOOS)/$(GOARCH)"
	@echo "🔨 Building Phantom C2..."
	@$(MAKE) server
	@$(MAKE) client
	@echo ""
	@echo "✅ Build complete!"
	@echo "   📦 phantom-server$(EXT)"
	@echo "   📦 phantom-client$(EXT)"
	@echo ""
	@echo "Run: ./phantom-server$(EXT)"

.PHONY: server
server:
	@echo "   Building server..."
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
		$(GO) build -mod=vendor -trimpath $(TAGS),server $(LDFLAGS) \
		-o phantom-server$(EXT) ./server || (echo "❌ Server build failed!" && exit 1)

.PHONY: client
client:
	@echo "   Building client..."
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 \
		$(GO) build -mod=vendor -trimpath $(TAGS),client $(LDFLAGS) \
		-o phantom-client$(EXT) ./client || (echo "❌ Client build failed!" && exit 1)

#
# Cross-compilation targets
#
.PHONY: linux
linux:
	@$(MAKE) GOOS=linux GOARCH=amd64 default

.PHONY: linux-arm64
linux-arm64:
	@$(MAKE) GOOS=linux GOARCH=arm64 default

.PHONY: macos
macos:
	@$(MAKE) GOOS=darwin GOARCH=arm64 default

.PHONY: macos-amd64
macos-amd64:
	@$(MAKE) GOOS=darwin GOARCH=amd64 default

.PHONY: windows
windows:
	@$(MAKE) GOOS=windows GOARCH=amd64 default

.PHONY: windows-arm64
windows-arm64:
	@$(MAKE) GOOS=windows GOARCH=arm64 default

#
# Build all platforms
#
.PHONY: all-platforms
all-platforms: clean
	@echo "🔨 Building for all platforms..."
	@$(MAKE) GOOS=linux GOARCH=amd64 server client
	@mv phantom-server phantom-server-linux-amd64
	@mv phantom-client phantom-client-linux-amd64
	@$(MAKE) GOOS=darwin GOARCH=arm64 server client
	@mv phantom-server phantom-server-macos-arm64
	@mv phantom-client phantom-client-macos-arm64
	@$(MAKE) GOOS=windows GOARCH=amd64 server client
	@mv phantom-server.exe phantom-server-windows-amd64.exe
	@mv phantom-client.exe phantom-client-windows-amd64.exe
	@echo "✅ All platforms built!"

#
# Protobuf generation
#
.PHONY: pb
pb:
	@echo "🔨 Generating protobuf files..."
	@protoc -I protobuf/ protobuf/commonpb/common.proto --go_out=paths=source_relative:protobuf/
	@protoc -I protobuf/ protobuf/phantompb/phantom.proto --go_out=paths=source_relative:protobuf/
	@protoc -I protobuf/ protobuf/clientpb/client.proto --go_out=paths=source_relative:protobuf/
	@protoc -I protobuf/ protobuf/dnspb/dns.proto --go_out=paths=source_relative:protobuf/
	@protoc -I protobuf/ protobuf/rpcpb/services.proto --go_out=paths=source_relative:protobuf/ --go-grpc_out=protobuf/ --go-grpc_opt=paths=source_relative
	@echo "✅ Protobuf files generated"

#
# Debug build (with symbols)
#
.PHONY: debug
debug: clean
	@echo "🔨 Building debug version..."
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
		$(GO) build -mod=vendor $(TAGS),server -o phantom-server$(EXT) ./server
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 \
		$(GO) build -mod=vendor $(TAGS),client -o phantom-client$(EXT) ./client
	@echo "✅ Debug build complete"

#
# Clean
#
.PHONY: clean
clean:
	@rm -f phantom-server phantom-server.exe phantom-server-* 2>/dev/null || true
	@rm -f phantom-client phantom-client.exe phantom-client-* 2>/dev/null || true

.PHONY: clean-all
clean-all: clean
	@rm -rf ./server/assets/fs/darwin/amd64 2>/dev/null || true
	@rm -rf ./server/assets/fs/darwin/arm64 2>/dev/null || true
	@rm -rf ./server/assets/fs/windows/amd64 2>/dev/null || true
	@rm -rf ./server/assets/fs/linux/amd64 2>/dev/null || true
	@rm -f ./server/assets/fs/*.zip 2>/dev/null || true
	@rm -f ./.downloaded_assets 2>/dev/null || true

#
# Help
#
.PHONY: help
help:
	@echo "Phantom C2 Framework - Build System"
	@echo ""
	@echo "Quick Start:"
	@echo "  make              # Auto-detect OS and build"
	@echo "  make linux        # Build for Linux"
	@echo "  make macos        # Build for macOS"
	@echo "  make windows      # Build for Windows"
	@echo ""
	@echo "Advanced:"
	@echo "  make all-platforms    # Build for all platforms"
	@echo "  make debug            # Build with debug symbols"
	@echo "  make pb               # Regenerate protobuf files"
	@echo "  make clean            # Remove build artifacts"
	@echo "  make clean-all        # Deep clean including assets"
	@echo ""
	@echo "Current system: $(GOOS)/$(GOARCH)"
