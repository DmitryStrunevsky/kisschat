# GitHub Actions Workflow Summary

## Overview
This GitHub Actions workflow (`Build and Release`) automatically builds binaries for Linux and Windows, packages them, and creates GitHub releases on every commit to the main branch.

## Workflow Details

### Triggers
- **Push to main branch**: Creates a release with binaries
- **Pull requests to main branch**: Builds and tests the binaries (no release created)

### Build Process
1. **Matrix Build**: Builds for both Linux (amd64) and Windows (amd64) simultaneously
2. **Binary Creation**: 
   - Linux: Creates `kisschat` executable
   - Windows: Creates `kisschat.exe` executable
3. **Packaging**:
   - Linux: Compressed as `kisschat-linux-amd64.tar.gz`
   - Windows: Compressed as `kisschat-windows-amd64.zip`

### Release Process (Main Branch Only)
1. **Artifact Download**: Downloads the built and packaged binaries
2. **Tag Generation**: Creates a unique tag using timestamp and commit hash (e.g., `v20241209-143022-abc1234`)
3. **Release Creation**: Creates a GitHub release with:
   - Descriptive release notes
   - Automated changelog information
   - Both Linux and Windows binaries as downloadable assets

## File Structure
- **Linux Package**: `kisschat-linux-amd64.tar.gz` containing the `kisschat` executable
- **Windows Package**: `kisschat-windows-amd64.zip` containing the `kisschat.exe` executable

## Usage
Once the workflow runs successfully:
1. Go to the repository's "Releases" page
2. Download the appropriate binary for your platform
3. Extract the archive
4. Run the executable according to the README instructions

## Technical Notes
- Uses Go 1.18 for building
- Uses modern GitHub Actions (v4 artifacts, softprops/action-gh-release@v2)
- Includes proper error handling and artifact management
- Follows GitHub Actions best practices for security and efficiency