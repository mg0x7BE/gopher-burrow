// Command build packs src/ into a GoLand plugin jar and renders the design preview.
// Run it from the repository root: go run ./cmd/build
package main

import (
	"archive/zip"
	"bytes"
	"encoding/json/v2"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"text/template"
)

const (
	paletteFile     = "palette.json"
	srcDir          = "src"
	buildDir        = "build"
	pluginXML       = "src/META-INF/plugin.xml"
	previewTemplate = "design/preview.html.tmpl"
	previewOutput   = "design/preview.html"
)

var placeholder = regexp.MustCompile(`\$([a-z][a-z0-9-]*)`)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "build:", err)
		os.Exit(1)
	}
}

func run() error {
	palette, err := loadPalette()
	if err != nil {
		return err
	}
	version, err := pluginVersion()
	if err != nil {
		return err
	}

	jar := filepath.Join(buildDir, "gopher-burrow-"+version+".jar")
	if err := cleanBuildDir(); err != nil {
		return err
	}
	if err := writeJar(jar, palette); err != nil {
		return err
	}
	if err := writePreview(palette); err != nil {
		return err
	}

	fmt.Println("Built", jar)
	fmt.Println("Rendered", previewOutput)
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

func pluginVersion() (string, error) {
	data, err := os.ReadFile(pluginXML)
	if err != nil {
		return "", err
	}
	var plugin struct {
		Version string `xml:"version"`
	}
	if err := xml.Unmarshal(data, &plugin); err != nil {
		return "", fmt.Errorf("%s: %w", pluginXML, err)
	}
	if plugin.Version == "" {
		return "", fmt.Errorf("%s: missing <version>", pluginXML)
	}
	return plugin.Version, nil
}

func cleanBuildDir() error {
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return err
	}
	old, err := filepath.Glob(filepath.Join(buildDir, "*.jar"))
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

func writeJar(path string, palette map[string]string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, f.Close()) }()

	zw := zip.NewWriter(f)
	walkErr := filepath.WalkDir(srcDir, func(file string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(srcDir, file)
		if err != nil {
			return err
		}
		name := filepath.ToSlash(rel)

		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if ext := filepath.Ext(file); slices.Contains([]string{".xml", ".json", ".svg"}, ext) {
			if data, err = expand(data, palette); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			if err := validate(ext, data); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}

		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	})
	return errors.Join(walkErr, zw.Close())
}

func expand(data []byte, palette map[string]string) ([]byte, error) {
	var unknown []string
	out := placeholder.ReplaceAllFunc(data, func(m []byte) []byte {
		name := string(m[1:])
		hex, ok := palette[name]
		if !ok {
			unknown = append(unknown, name)
			return m
		}
		return []byte(hex)
	})
	if len(unknown) > 0 {
		slices.Sort(unknown)
		return nil, fmt.Errorf("unknown palette colors: %v", slices.Compact(unknown))
	}
	return out, nil
}

// validate catches broken files early, because GoLand silently ignores them.
func validate(ext string, data []byte) error {
	if ext == ".json" {
		var v any
		return json.Unmarshal(data, &v)
	}
	dec := xml.NewDecoder(bytes.NewReader(data))
	for {
		if _, err := dec.Token(); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}

func writePreview(palette map[string]string) error {
	tmpl, err := template.ParseFiles(previewTemplate)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, map[string]any{"Palette": palette}); err != nil {
		return err
	}
	return os.WriteFile(previewOutput, out.Bytes(), 0o644)
}
