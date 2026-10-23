# ![](src/jetbrains/META-INF/pluginIcon.svg) Gopher Burrow

[![JetBrains Marketplace](https://img.shields.io/jetbrains/plugin/v/35081)](https://plugins.jetbrains.com/plugin/35081)
[![Downloads](https://img.shields.io/jetbrains/plugin/d/35081)](https://plugins.jetbrains.com/plugin/35081)
[![CI](https://github.com/mg0x7BE/gopher-burrow/actions/workflows/ci.yml/badge.svg)](https://github.com/mg0x7BE/gopher-burrow/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/mg0x7BE/gopher-burrow)](LICENSE)

Dark theme for JetBrains IDEs and Zed. The JetBrains version is built on Islands Dark.

![Gopher Burrow in GoLand](docs/screenshots/go.png)

- Colors tuned for Go and most popular languages.
- Made for GoLand, but I've also tried it in IntelliJ IDEA, Rider, RustRover and CLion, so it should look good there too.

## Install

- JetBrains IDEs: Settings > Plugins > Marketplace, search for Gopher Burrow.
- Zed: search for Gopher Burrow in Extensions, or download `gopher-burrow-zed-*.json` from [Releases](https://github.com/mg0x7BE/gopher-burrow/releases) and put it in the Zed themes folder (`~/.config/zed/themes` on macOS and Linux, `%APPDATA%\Zed\themes` on Windows).

## Build

```
go run ./cmd/build
```

## Release

Add a new version at the top of `CHANGELOG.md` in a pull request. After it's merged, GitHub Actions publishes a release with both files. Then upload the jar to JetBrains Marketplace and open a pull request in [zed-industries/extensions](https://github.com/zed-industries/extensions).

## License

Apache 2.0
