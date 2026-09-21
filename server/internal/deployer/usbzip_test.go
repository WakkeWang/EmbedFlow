package deployer

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/WakkeWang/EmbedFlow/server/internal/store"
)

// USB stick zip generation (requirement 3.1.3 second half). The opener
// serves in-memory content so the tests assert zip layout + checksum
// content without touching disk.

func fakeOpener(contents map[string]string) opener {
	return func(a ArtifactView) (io.ReadCloser, error) {
		s, ok := contents[a.Name]
		if !ok {
			return nil, fs.ErrNotExist
		}
		return io.NopCloser(strings.NewReader(s)), nil
	}
}

func zipEntries(t *testing.T, raw []byte) map[string]string {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	out := map[string]string{}
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open entry %s: %v", f.Name, err)
		}
		data, _ := io.ReadAll(rc)
		rc.Close()
		out[f.Name] = string(data)
	}
	return out
}

var usbArts = []ArtifactView{
	{ID: 1, Name: "board.itb", Size: 8, Checksum: "aa11"},
	{ID: 2, Name: "recovery-p0133-sample-generic-2026-09-19-x1.sh", Size: 10, Checksum: "bb22"},
	{ID: 3, Name: "linxos-ipc-d2000-1.0.0.20250715-arm64.tar.gz", Size: 12, Checksum: "cc33"},
}

const (
	itbBody      = "ITBBYTES!"
	recoveryBody = "RECOVERYSH"
	rootfsBody   = "ROOTFSTAR "
)

func usbOpener() opener {
	return fakeOpener(map[string]string{
		"board.itb":                                      itbBody,
		"recovery-p0133-sample-generic-2026-09-19-x1.sh": recoveryBody,
		"linxos-ipc-d2000-1.0.0.20250715-arm64.tar.gz":   rootfsBody,
	})
}

func md5Of(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func TestUSBZip_LayoutAndChecksumFile(t *testing.T) {
	disk := store.USBDisk{
		RootName: "P0133-sample-9.9.9",
		Layout: []store.USBLayout{
			{Glob: "board.itb", Dir: "."},
			{Glob: "*-generic-*.sh", Dir: "scripts"},
			{Glob: "*.tar.gz", Dir: "."},
		},
		Checksums: []store.USBChecksumField{
			{Name: "rootfsfile", Glob: "*.tar.gz"},
			{Name: "patchfile", Glob: "*-generic-*.sh"},
			{Name: "itbfile", Glob: "board.itb"},
		},
	}
	var buf bytes.Buffer
	if err := BuildUSBZip(&buf, "P0133-sample-9.9.9", disk, usbArts, usbOpener(), time.Now()); err != nil {
		t.Fatalf("BuildUSBZip: %v", err)
	}
	entries := zipEntries(t, buf.Bytes())

	root := "P0133-sample-9.9.9"
	if got := entries[root+"/board.itb"]; got != itbBody {
		t.Fatalf("board.itb content = %q", got)
	}
	if got := entries[root+"/scripts/recovery-p0133-sample-generic-2026-09-19-x1.sh"]; got != recoveryBody {
		t.Fatalf("recovery.sh content = %q", got)
	}
	csum := entries[root+"/checksum.md5"]
	want := "rootfsfile=linxos-ipc-d2000-1.0.0.20250715-arm64.tar.gz\n" +
		"rootfsfile_md5sum=" + md5Of(rootfsBody) + "\n" +
		"patchfile=recovery-p0133-sample-generic-2026-09-19-x1.sh\n" +
		"patchfile_md5sum=" + md5Of(recoveryBody) + "\n" +
		"itbfile=board.itb\n" +
		"itbfile_md5sum=" + md5Of(itbBody) + "\n"
	if csum != want {
		t.Fatalf("checksum.md5 =\n%q\nwant\n%q", csum, want)
	}
	if _, ok := entries["README.txt"]; !ok {
		t.Fatal("README.txt missing")
	}
}

func TestUSBZip_NoLayout_GoesFlatIntoRoot(t *testing.T) {
	disk := store.USBDisk{RootName: "stick"}
	var buf bytes.Buffer
	if err := BuildUSBZip(&buf, "stick", disk, usbArts, usbOpener(), time.Now()); err != nil {
		t.Fatalf("BuildUSBZip: %v", err)
	}
	entries := zipEntries(t, buf.Bytes())
	for _, name := range []string{"stick/board.itb", "README.txt"} {
		if _, ok := entries[name]; !ok {
			t.Fatalf("entry %s missing (have %v)", name, keysOf(entries))
		}
	}
	// No Checksums configured -> no checksum file.
	if _, ok := entries["stick/checksum.md5"]; ok {
		t.Fatal("checksum.md5 present without configured fields")
	}
}

func TestUSBZip_MissingPattern_ListedInError(t *testing.T) {
	disk := store.USBDisk{
		RootName:  "r",
		Checksums: []store.USBChecksumField{{Name: "x", Glob: "nope-*"}},
	}
	err := ResolveUSB(disk, usbArts)
	ze, ok := err.(*USBZipError)
	if !ok {
		t.Fatalf("want USBZipError, got %v", err)
	}
	if len(ze.Missing) != 1 || ze.Missing[0] != "nope-*" {
		t.Fatalf("missing = %v", ze.Missing)
	}
}

func TestUSBZip_AmbiguousPattern_ListedInError(t *testing.T) {
	arts := append([]ArtifactView{}, usbArts...)
	arts = append(arts, ArtifactView{ID: 9, Name: "board2.itb", Size: 4})
	disk := store.USBDisk{
		RootName:  "r",
		Checksums: []store.USBChecksumField{{Name: "itbfile", Glob: "board*.itb"}},
	}
	err := ResolveUSB(disk, arts)
	ze, ok := err.(*USBZipError)
	if !ok {
		t.Fatalf("want USBZipError, got %v", err)
	}
	if len(ze.Missing) != 1 || ze.Missing[0] != "board*.itb" {
		t.Fatalf("missing = %v", ze.Missing)
	}
}

func TestUSBZip_RootSanitized(t *testing.T) {
	disk := store.USBDisk{RootName: ""}
	var buf bytes.Buffer
	if err := BuildUSBZip(&buf, "a/b\\c:d", disk, usbArts[:1], usbOpener(), time.Now()); err != nil {
		t.Fatalf("BuildUSBZip: %v", err)
	}
	entries := zipEntries(t, buf.Bytes())
	if _, ok := entries["a-b-c-d/board.itb"]; !ok {
		t.Fatalf("sanitized root entry missing (have %v)", keysOf(entries))
	}
}

func TestUSBZip_LayoutDirEscapeDegraded(t *testing.T) {
	disk := store.USBDisk{
		RootName: "r",
		Layout:   []store.USBLayout{{Glob: "board.itb", Dir: "../escape"}},
	}
	var buf bytes.Buffer
	if err := BuildUSBZip(&buf, "r", disk, usbArts[:1], usbOpener(), time.Now()); err != nil {
		t.Fatalf("BuildUSBZip: %v", err)
	}
	entries := zipEntries(t, buf.Bytes())
	if _, ok := entries["r/board.itb"]; !ok {
		t.Fatalf("escaped dir not degraded to root (have %v)", keysOf(entries))
	}
	if _, ok := entries["r/../escape/board.itb"]; ok {
		t.Fatal("zip escape entry exists")
	}
}

func TestResolveUSB_Clean(t *testing.T) {
	disk := store.USBDisk{
		RootName: "r",
		Layout:   []store.USBLayout{{Glob: "board.itb", Dir: "."}},
		Checksums: []store.USBChecksumField{
			{Name: "rootfsfile", Glob: "*.tar.gz"},
		},
	}
	if err := ResolveUSB(disk, usbArts); err != nil {
		t.Fatalf("ResolveUSB: %v", err)
	}
}

func keysOf(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
