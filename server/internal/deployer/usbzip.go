// USB-stick zip generation (requirement 3.1.3 second half): the rule
// declares how a build record's artifacts must be laid out on the stick,
// and the server packages them (streamed, nothing lands on disk) plus a
// generated checksum file and an operator README.
//
// Three rule-configured pieces (the deploy payload's usb_disk section,
// store.USBDisk):
//   - RootName: the zip's top folder, a Render template, so the stick path
//     the device script expects (<ver>/<version folder>/) is reproduced by
//     unzipping.
//   - Layout: artifact glob -> directory inside that root.
//   - Checksums: custom fields (rootfsfile / patchfile / itbfile, ...),
//     each bound to one artifact glob. The field's "<name>=<actual file>"
//     and "<name>_md5sum=<live md5>" lines are generated at package time --
//     artifact names carry versions and change per build, so nothing is
//     stored and md5 is always computed fresh (the build-time checksum
//     setting may be sha256, so it cannot serve here).
package deployer

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// USBZipError reports which configured patterns failed to resolve -- the
// operator fixes the rule or picks another build record; a half package is
// never produced.
type USBZipError struct {
	Missing []string
}

func (e *USBZipError) Error() string {
	return "artifacts not found for patterns: " + strings.Join(e.Missing, ", ")
}

// opener returns a reader for one artifact's archived file; the cmd layer
// binds it to <data>/artifacts/<recordID>/<name>.
type opener func(a ArtifactView) (io.ReadCloser, error)

// ResolveUSB pre-checks a rule's stick configuration against one build
// record's artifact list: every layout glob and every checksum glob must
// match exactly one artifact. Call it before any bytes are written so a
// bad configuration fails as a clean HTTP error, never as a half zip.
func ResolveUSB(disk store.USBDisk, arts []ArtifactView) error {
	var missing []string
	matchOne := func(glob string) bool { return len(matchArtifacts(arts, glob)) == 1 }
	for _, l := range disk.Layout {
		if !matchOne(l.Glob) {
			missing = append(missing, l.Glob)
		}
	}
	for _, f := range disk.Checksums {
		if !matchOne(f.Glob) {
			missing = append(missing, f.Glob)
		}
	}
	if len(missing) > 0 {
		return &USBZipError{Missing: missing}
	}
	return nil
}

// BuildUSBZip streams the stick zip into w. rootName must be pre-rendered
// by the caller (it carries build-time variables) and ResolveUSB must have
// passed: with that precondition every write below can only fail on IO.
// With no layout entries the artifacts go flat into the root.
func BuildUSBZip(w io.Writer, rootName string, disk store.USBDisk, arts []ArtifactView, open opener, now time.Time) error {
	type placed struct {
		art ArtifactView
		dir string
	}
	var placements []placed
	fieldHits := map[string]ArtifactView{}
	for _, f := range disk.Checksums {
		m := matchArtifacts(arts, f.Glob)
		// ResolveUSB guarantees exactly one; the guard keeps the builder
		// safe when called without the pre-check (tests, future paths).
		if len(m) == 1 {
			fieldHits[f.Glob] = m[0]
		}
	}
	if len(disk.Layout) == 0 {
		for _, a := range arts {
			placements = append(placements, placed{art: a, dir: "."})
		}
	} else {
		for _, l := range disk.Layout {
			m := matchArtifacts(arts, l.Glob)
			if len(m) != 1 {
				// ResolveUSB guarantees exactly one; skip otherwise (tests,
				// future paths) rather than index out of range.
				continue
			}
			placements = append(placements, placed{art: m[0], dir: cleanZipDir(l.Dir)})
		}
	}

	root := sanitizeRoot(rootName)
	zw := zip.NewWriter(w)
	seen := map[string]bool{}
	writeEntry := func(name string, content func(io.Writer) error) error {
		if seen[name] {
			return fmt.Errorf("usbzip: duplicate entry %q", name)
		}
		seen[name] = true
		hdr := &zip.FileHeader{Name: name, Method: zip.Deflate, Modified: now}
		hdr.SetMode(0o755)
		fw, err := zw.CreateHeader(hdr)
		if err != nil {
			return err
		}
		return content(fw)
	}

	// Artifacts into their mapped directories.
	for _, p := range placements {
		art, dir := p.art, p.dir
		name := path.Join(root, dir, art.Name)
		if err := writeEntry(name, func(fw io.Writer) error {
			rc, err := open(art)
			if err != nil {
				return fmt.Errorf("usbzip: open %s: %w", art.Name, err)
			}
			defer rc.Close()
			_, err = io.Copy(fw, rc)
			return err
		}); err != nil {
			return err
		}
	}

	// The checksum file (rule-configured custom fields, live md5) in the root.
	if len(disk.Checksums) > 0 {
		var lines []string
		for _, f := range disk.Checksums {
			a := fieldHits[f.Glob]
			sum, err := md5Artifact(open, a)
			if err != nil {
				return err
			}
			lines = append(lines, f.Name+"="+a.Name, f.Name+"_md5sum="+sum)
		}
		content := strings.Join(lines, "\n") + "\n"
		if err := writeEntry(path.Join(root, "checksum.md5"), func(fw io.Writer) error {
			_, err := fw.Write([]byte(content))
			return err
		}); err != nil {
			return err
		}
	}

	// Operator guide: unzip, then copy the root folder (the folder itself)
	// to the stick's top level, keeping the device script's expected path.
	guide := []string{
		"EmbedFlow USB deployment package",
		"",
		"1. Unzip this archive.",
		"2. Copy the folder \"" + root + "\" (the whole folder) to the top level of the USB stick.",
		"3. Insert the stick into the device and follow the deployment rule's manual steps.",
		"",
		"Generated " + now.Format("2006-01-02 15:04:05"),
		"",
	}
	gcontent := strings.Join(guide, "\n")
	if err := writeEntry("README.txt", func(fw io.Writer) error {
		_, err := fw.Write([]byte(gcontent))
		return err
	}); err != nil {
		return err
	}

	return zw.Close()
}

// md5Artifact computes the artifact's md5 fresh from its archived file.
func md5Artifact(open opener, a ArtifactView) (string, error) {
	rc, err := open(a)
	if err != nil {
		return "", fmt.Errorf("usbzip: open %s: %w", a.Name, err)
	}
	defer rc.Close()
	h := md5.New()
	if _, err := io.Copy(h, rc); err != nil {
		return "", fmt.Errorf("usbzip: md5 %s: %w", a.Name, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// matchArtifacts returns the artifacts whose Name matches the glob
// (filepath.Match semantics on the file name: board.itb exact,
// *-generic-*.sh).
func matchArtifacts(arts []ArtifactView, glob string) []ArtifactView {
	if glob == "" {
		return nil
	}
	var out []ArtifactView
	for _, a := range arts {
		ok, err := filepath.Match(glob, a.Name)
		if err != nil || !ok {
			continue
		}
		out = append(out, a)
	}
	return out
}

// sanitizeRoot strips path separators from the rendered root name: the
// stick layout is one folder, not a tree the operator cannot eyeball.
func sanitizeRoot(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "usb-deploy"
	}
	return strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			return '-'
		}
		return r
	}, name)
}

// cleanZipDir normalizes a layout dir to a slash path, rejecting escapes
// from the zip root (".." and absolute paths degrade to the root itself).
func cleanZipDir(dir string) string {
	dir = strings.TrimSpace(dir)
	if dir == "" || dir == "." || dir == "/" {
		return "."
	}
	p := path.Clean(strings.ReplaceAll(dir, "\\", "/"))
	if strings.HasPrefix(p, "..") || path.IsAbs(p) {
		return "."
	}
	return p
}
