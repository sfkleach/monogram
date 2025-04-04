We are collaborating on a Go program, codename `monogram`. The `monogram` tool
translates from `monogram` notation into XML, JSON and other formats. The 
notation is designed to represent program-like texts. However it is just a
notation and not a programming language, although it does have an opinionated
grammar. Consequently it has no built-in variables, no built-in operators and
even the reserved words are dynamically discovered during the parse.

We have completed a good first version of that program.

The repo that it resides in is a monorepo with a parallel implementation,
written in some other programming language (Pop-11). As a consequence the
go folder is not at the top-level of the repo but in `~/go/monogram/`. The
application itself is in `~/go/monogram/cmd/monogram/main.go`.

When we publish a new release we use the following GitHub workflow, which 
is triggered on a tag push:

```yaml
name: Build and Release Monogram Binaries on Tag Push

on:
  push:
    tags:
      - "v*"  # Matches tags starting with "v"

jobs:
  build-and-release:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout Repository
        uses: actions/checkout@v3
        with:
          fetch-depth: 0  # Ensure full history is available, including tags

      - name: Check Out Main Branch
        run: |
          git fetch origin main
          git checkout main

      - name: Extract Git Tag
        id: get_tag
        run: |
          echo "TAG=${GITHUB_REF#refs/tags/}"
          echo "TAG=${GITHUB_REF#refs/tags/}" >> $GITHUB_ENV

      - name: Update version.go
        run: |
          echo "package lib" > go/monogram/lib/version.go
          echo "const Version = \"${TAG}\"" >> go/monogram/lib/version.go

      - name: Commit Version Update
        run: |
          git config --global user.name "github-actions[bot]"
          git config --global user.email "actions@github.com"
          git add go/monogram/lib/version.go
          git commit -m "Update version.go to ${TAG}"
          git push origin HEAD
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        

      - name: Set Up Go Environment
        uses: actions/setup-go@v4
        with:
          go-version: "1.23"

      - name: Build Binary for Linux x86_64
        run: |
          mkdir -p dist/linux/x86_64
          cd go/monogram/cmd/monogram
          GOOS=linux GOARCH=amd64 go build -o ../../../../dist/linux/x86_64/monogram
        shell: bash

      - name: Build Binary for Linux arm64
        run: |
          mkdir -p dist/linux/arm64
          cd go/monogram/cmd/monogram
          GOOS=linux GOARCH=arm64 go build -o ../../../../dist/linux/arm64/monogram
        shell: bash

      - name: Build Binary for Windows
        run: |
          mkdir -p dist/windows
          cd go/monogram/cmd/monogram
          GOOS=windows GOARCH=amd64 go build -o ../../../../dist/windows/monogram.exe
        shell: bash

      - name: Build Binary for macOS (Intel)
        run: |
          mkdir -p dist/macos/intel
          cd go/monogram/cmd/monogram
          GOOS=darwin GOARCH=amd64 go build -o ../../../../dist/macos/intel/monogram
        shell: bash

      - name: Build Binary for macOS (ARM)
        run: |
          mkdir -p dist/macos/arm64
          cd go/monogram/cmd/monogram
          GOOS=darwin GOARCH=arm64 go build -o ../../../../dist/macos/arm64/monogram
        shell: bash

      - name: Create GitHub Release
        id: create_release
        uses: actions/create-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tag_name: ${{ env.TAG }}
          release_name: "Release ${{ env.TAG }}"
          body: "Automatically generated release for version ${{ env.TAG }}."
          draft: true
          prerelease: false

      # Compression steps that also rename the files to the desired names
      - name: Compress Linux x86_64 Binary
        run: |
          tar -czvf dist/linux/x86_64/monogram-linux-x86_64.tar.gz -C dist/linux/x86_64 monogram
        shell: bash

      - name: Compress Linux arm64 Binary
        run: |
          tar -czvf dist/linux/arm64/monogram-linux-arm64.tar.gz -C dist/linux/arm64 monogram
        shell: bash

      - name: Compress Windows Binary
        run: |
          zip -j dist/windows/monogram-windows.zip dist/windows/monogram.exe
        shell: bash

      - name: Compress macOS Intel Binary
        run: |
          zip -j dist/macos/intel/monogram-macos-intel.zip dist/macos/intel/monogram
        shell: bash

      - name: Compress macOS ARM Binary
        run: |
          zip -j dist/macos/arm64/monogram-macos-arm64.zip dist/macos/arm64/monogram
        shell: bash

      # Create the GitHub release with gh CLI
      - name: Create GitHub Release
        run: |
          gh release create "$TAG" \
            -t "Release $TAG" \
            -n "Automatically generated release for version $TAG." \
            --repo ${{ github.repository }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      # Upload each compressed archive using the gh CLI. The file names are already set by the compression step!
      - name: Upload Linux x86_64 Archive with gh
        run: |
          gh release upload "$TAG" dist/linux/x86_64/monogram-linux-x86_64.tar.gz --clobber --repo ${{ github.repository }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Upload Linux arm64 Archive with gh
        run: |
          gh release upload "$TAG" dist/linux/arm64/monogram-linux-arm64.tar.gz --clobber --repo ${{ github.repository }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Upload Windows Archive with gh
        run: |
          gh release upload "$TAG" dist/windows/monogram-windows.zip --clobber --repo ${{ github.repository }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Upload macOS Intel Archive with gh
        run: |
          gh release upload "$TAG" dist/macos/intel/monogram-macos-intel.zip --clobber --repo ${{ github.repository }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Upload macOS ARM Archive with gh
        run: |
          gh release upload "$TAG" dist/macos/arm64/monogram-macos-arm64.zip --clobber --repo ${{ github.repository }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```


 

