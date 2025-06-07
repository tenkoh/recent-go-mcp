# Release Data

This directory contains Go release data files that are embedded into the binary.

## Structure

- `releases/` - Individual JSON files for each Go version
  - `go1.23.json` - Go 1.23 release data
  - `go1.22.json` - Go 1.22 release data  
  - `go1.21.json` - Go 1.21 release data

## Adding New Versions

To add a new Go version:

1. Create a new JSON file in `releases/` following the naming pattern `go{version}.json`
2. Add the embed directive in `data.go`
3. Update the `releaseFiles` map in the `init()` function
4. Follow the existing JSON structure for consistency

## JSON Structure

Each file contains a single Go release object with:
- `version`: Go version string (e.g., "1.23")
- `release_date`: ISO 8601 date string
- `summary`: Brief description of the release
- `changes`: Array of general changes with category, description, and impact
- `packages`: Map of package names to arrays of package-specific changes

## Impact Types

- `new`: New feature or function
- `enhancement`: Improvement to existing functionality
- `performance`: Performance optimization
- `breaking`: Breaking change requiring migration
- `deprecation`: Deprecated feature