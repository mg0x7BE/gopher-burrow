<p align="center">
  <img src="src/META-INF/pluginIcon.svg" width="120" alt="Gopher Burrow icon">
</p>

<h1 align="center">Gopher Burrow</h1>

<p align="center">Dark theme for GoLand, built on Islands Dark and tuned for OLED screens.</p>

![Gopher Burrow in GoLand](docs/screenshots/go.png)

- Real black around the panels, soft violet where you read.
- Go code you can scan at a glance: types, functions and fields each get their own color.
- The same colors follow you into YAML, SQL, Dockerfiles and the rest of the repo.

## Install

Grab the jar from [Releases](../../releases), open Settings > Plugins > gear icon > Install Plugin from Disk, then pick Gopher Burrow in Settings > Appearance > Theme. Needs GoLand 2026.2 or newer.

## Build

```
go run ./cmd/build
```

Colors live in [palette.json](palette.json). The [samples](samples) folder has a file for every language the theme covers.

## Credits

Based on Islands Dark from the [IntelliJ Platform](https://github.com/JetBrains/intellij-community) (Apache 2.0). The icon is our take on the Go gopher by Renee French (CC BY 4.0). Not affiliated with Google, the Go project or JetBrains.

## License

Apache 2.0, see [LICENSE](LICENSE) and [NOTICE](NOTICE).
