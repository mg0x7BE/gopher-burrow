// Command build turns src/jetbrains into a JetBrains plugin jar and src/zed into a Zed extension,
// filling in colors from palette.json and the version and release notes from the newest
// entry in CHANGELOG.md.
// Run it from the repository root: go run ./cmd/build
package main

import (
	"archive/zip"
	"bytes"
	"encoding/json/v2"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"html"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const (
	paletteFile    = "palette.json"
	changelogFile  = "CHANGELOG.md"
	jetbrainsSrc   = "src/jetbrains"
	jetbrainsBuild = "build-jetbrains"
	zedSrc         = "src/zed"
	zedBuild       = "build-zed"
)

var (
	placeholder   = regexp.MustCompile(`\$([a-z][a-z0-9-]*)`)
	versionFormat = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)

	printVersion = flag.Bool("version", false, "print the newest version from CHANGELOG.md and exit")
	printNotes   = flag.Bool("notes", false, "print the newest release notes from CHANGELOG.md and exit")
)

type release struct {
	version string
	notes   string
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "build:", err)
		os.Exit(1)
	}
}

func run() error {
	r, err := loadRelease()
	if err != nil {
		return err
	}
	switch {
	case *printVersion:
		fmt.Println(r.version)
		return nil
	case *printNotes:
		fmt.Println(r.notes)
		return nil
	}

	palette, err := loadPalette()
	if err != nil {
		return err
	}
	vars := maps.Clone(palette)
	vars["version"] = r.version
	vars["release-notes"] = notesHTML(r.notes)

	jar := filepath.Join(jetbrainsBuild, "gopher-burrow-jetbrains-"+r.version+".jar")
	if err := cleanBuildDir(); err != nil {
		return err
	}
	if err := writeJar(jar, vars); err != nil {
		return err
	}

	fmt.Println("Built", jar)

	if err := writeZed(vars); err != nil {
		return err
	}
	fmt.Println("Built", zedBuild)
	return nil
}

func loadPalette() (map[string]string, error) {
	data, err := os.ReadFile(paletteFile)
	if err != nil {
		return nil, err
	}
	var palette map[string]string
	if err := json.Unmarshal(data, &palette); err != nil {
		return nil, fmt.Errorf("%s: %w", paletteFile, err)
	}
	return palette, nil
}

func loadRelease() (release, error) {
	data, err := os.ReadFile(changelogFile)
	if err != nil {
		return release{}, err
	}
	r, err := parseChangelog(string(data))
	if err != nil {
		return release{}, fmt.Errorf("%s: %w", changelogFile, err)
	}
	return r, nil
}

// parseChangelog returns the newest entry: the first "## 1.2.3" heading and the text under it.
func parseChangelog(text string) (release, error) {
	_, rest, ok := strings.Cut("\n"+text, "\n## ")
	if !ok {
		return release{}, errors.New("no version heading like ## 1.2.3")
	}
	heading, body, _ := strings.Cut(rest, "\n")
	version := strings.TrimSpace(heading)
	if !versionFormat.MatchString(version) {
		return release{}, fmt.Errorf("%q is not a version like 1.2.3", version)
	}
	notes, _, _ := strings.Cut(body, "\n## ")
	notes = strings.TrimSpace(notes)
	if notes == "" {
		return release{}, fmt.Errorf("%s has no release notes", version)
	}
	return release{version, notes}, nil
}

// notesHTML turns plain-text paragraphs into the HTML that JetBrains shows as change notes.
func notesHTML(notes string) string {
	var b strings.Builder
	for _, p := range strings.Split(notes, "\n\n") {
		if p = strings.Join(strings.Fields(p), " "); p != "" {
			b.WriteString("<p>" + html.EscapeString(p) + "</p>")
		}
	}
	return b.String()
}

func cleanBuildDir() error {
	if err := os.MkdirAll(jetbrainsBuild, 0o755); err != nil {
		return err
	}
	old, err := filepath.Glob(filepath.Join(jetbrainsBuild, "*.jar"))
	if err != nil {
		return err
	}
	for _, f := range old {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func writeJar(path string, vars map[string]string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, f.Close()) }()

	zw := zip.NewWriter(f)
	walkErr := filepath.WalkDir(jetbrainsSrc, func(file string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(jetbrainsSrc, file)
		if err != nil {
			return err
		}
		data, err := render(file, vars)
		if err != nil {
			return err
		}

		w, err := zw.Create(filepath.ToSlash(rel))
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	})
	return errors.Join(walkErr, zw.Close())
}

// writeZed mirrors src/zed into build-zed, which is committed because Zed's
// registry reads extensions straight from the repository.
func writeZed(vars map[string]string) error {
	written := map[string]bool{}
	err := filepath.WalkDir(zedSrc, func(file string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(zedSrc, file)
		if err != nil {
			return err
		}
		data, err := render(file, vars)
		if err != nil {
			return err
		}
		out := filepath.Join(zedBuild, rel)
		written[out] = true
		return writeFile(out, data)
	})
	if err != nil {
		return err
	}

	// Zed only accepts a license inside the extension directory.
	license, err := os.ReadFile("LICENSE")
	if err != nil {
		return err
	}
	out := filepath.Join(zedBuild, "LICENSE")
	written[out] = true
	if err := writeFile(out, license); err != nil {
		return err
	}
	return removeStale(zedBuild, written)
}

// removeStale deletes what an earlier build left behind, but keeps the directory
// itself so a Zed dev extension pointing at it stays valid.
func removeStale(dir string, keep map[string]bool) error {
	var dirs []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		switch {
		case err != nil:
			return err
		case d.IsDir():
			dirs = append(dirs, path)
		case !keep[path]:
			return os.Remove(path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, d := range slices.Backward(dirs[1:]) {
		entries, err := os.ReadDir(d)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			if err := os.Remove(d); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// render reads a source file and fills in its placeholders.
func render(file string, vars map[string]string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ext := filepath.Ext(file)
	if !slices.Contains([]string{".xml", ".json", ".svg", ".toml"}, ext) {
		return data, nil
	}
	if data, err = expand(data, vars); err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	if err := validate(ext, data); err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	return data, nil
}

func expand(data []byte, vars map[string]string) ([]byte, error) {
	var unknown []string
	out := placeholder.ReplaceAllFunc(data, func(m []byte) []byte {
		name := string(m[1:])
		value, ok := vars[name]
		if !ok {
			unknown = append(unknown, name)
			return m
		}
		return []byte(value)
	})
	if len(unknown) > 0 {
		slices.Sort(unknown)
		return nil, fmt.Errorf("unknown placeholders: %v", slices.Compact(unknown))
	}
	return out, nil
}

// validate catches broken files early, because the editors silently ignore them.
// TOML is left unchecked, since the standard library has no parser for it.
func validate(ext string, data []byte) error {
	switch ext {
	case ".json":
		var v any
		return json.Unmarshal(data, &v)
	case ".xml", ".svg":
		dec := xml.NewDecoder(bytes.NewReader(data))
		for {
			if _, err := dec.Token(); err == io.EOF {
				return nil
			} else if err != nil {
				return err
			}
		}
	}
	return nil
}
