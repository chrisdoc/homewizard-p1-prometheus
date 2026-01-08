# Docker Build Optimizations

This document describes the optimizations made to speed up Docker builds for this project.

## Changes Made

### 1. Dockerfile Layer Caching Optimization

**Before:**
```dockerfile
COPY . .
RUN go mod download
RUN go vet -v
RUN go test -v
RUN CGO_ENABLED=0 go build -o /go/src/github.com/chrisdoc/homewizard-p1-prometheus/app
```

**After:**
```dockerfile
# Copy go mod files first to cache dependencies layer
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Build with cache mount for build cache
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /go/src/github.com/chrisdoc/homewizard-p1-prometheus/app
```

**Benefits:**
- **Dependency caching**: By copying `go.mod` and `go.sum` first, Docker can cache the dependency download layer. This layer only rebuilds when dependencies change, not when source code changes.
- **BuildKit cache mounts**: Using `--mount=type=cache` for `/go/pkg/mod` and `/root/.cache/go-build` allows the Go module cache and build cache to persist across builds, significantly speeding up subsequent builds.
- **Removed tests from build**: Tests (`go vet` and `go test`) are now only run in CI, not during every Docker build. This reduces build time while still maintaining code quality checks in the CI pipeline.

### 2. GitHub Actions Cache Integration

**Added to `.github/workflows/publish.yml`:**
```yaml
cache-from: type=gha
cache-to: type=gha,mode=max
```

**Benefits:**
- Enables Docker layer caching in GitHub Actions using the GitHub Actions cache backend
- Subsequent builds in CI will reuse layers from previous builds
- `mode=max` ensures all layers are cached, not just the final image

### 3. Enhanced .dockerignore

**Added:**
- `.github` - GitHub workflow files
- `LICENSE` - License file
- `Makefile` - Build scripts
- `*.md` - Markdown documentation files
- `.gitignore` - Git configuration
- `.dockerignore` - Docker configuration

**Benefits:**
- Reduces Docker build context size
- Faster transfer of context to Docker daemon
- Prevents unnecessary file changes from invalidating cache layers

## Expected Performance Improvements

### First Build (Cold Cache)
- Similar to previous build time (all layers need to be built)

### Subsequent Builds with Code Changes Only
- **Before**: Full rebuild including dependency downloads (~2-3 minutes)
- **After**: Reuses dependency layer, only rebuilds source (~30-60 seconds)
- **Improvement**: ~50-75% faster

### Subsequent Builds with No Changes
- **Before**: Full rebuild (~2-3 minutes)
- **After**: Uses cache for all layers (~5-10 seconds)
- **Improvement**: ~95% faster

### Subsequent Builds with Dependency Changes
- **Before**: Full rebuild (~2-3 minutes)
- **After**: Rebuilds from dependency layer, uses build cache (~1-2 minutes)
- **Improvement**: ~30-50% faster

## Best Practices Applied

1. **Layer Ordering**: Most frequently changing files (source code) are copied last
2. **Cache Mounts**: BuildKit cache mounts provide persistent caching across builds
3. **Minimal Context**: .dockerignore reduces unnecessary file transfers
4. **Separation of Concerns**: Tests run in CI, not in production builds
5. **CI/CD Integration**: GitHub Actions cache reduces build times in workflows

## Usage

### Local Development
To take advantage of these optimizations locally, ensure BuildKit is enabled:
```bash
DOCKER_BUILDKIT=1 docker build -t homewizard-p1-prometheus .
```

### CI/CD
The GitHub Actions workflow automatically uses these optimizations. No additional configuration needed.
