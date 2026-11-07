package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirectoryHelpers(t *testing.T) {
	if !IsRunningInContainer() {
		t.Skip("container detection depends on test runtime")
	}
	dirs := GetDirectories()
	if dirs.Config != "/config" || dirs.Data != "/data" || dirs.Logs != "/logs" {
		t.Fatalf("container directories = %#v", dirs)
	}
	root := t.TempDir()
	dirs = Directories{Config: filepath.Join(root, "config"), Data: filepath.Join(root, "data"), Logs: filepath.Join(root, "logs")}
	if err := EnsureDirectories(dirs); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{dirs.Config, dirs.Data, dirs.Logs} {
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			t.Fatalf("directory %q: %v", dir, err)
		}
	}
	extra := filepath.Join(root, "extra")
	if err := EnsureDir(extra); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(extra); err != nil || !info.IsDir() {
		t.Fatalf("extra directory: %v", err)
	}
	if GetBackupDir() != "/mnt/Backups/apimgr/quotes" {
		t.Fatalf("backup directory = %q", GetBackupDir())
	}
}

func TestDefaultDirs(t *testing.T) {
	cfg, data, logs := GetDefaultDirs("demo")
	if cfg == "" || data == "" || logs == "" {
		t.Fatalf("empty default directories: %q %q %q", cfg, data, logs)
	}
}

func TestDefaultDirsForOperatingSystems(t *testing.T) {
	values := map[string]string{}
	getenv := func(key string) string { return values[key] }

	config, data, logs := defaultDirsForOS("demo", "linux", true, "/home/user", getenv)
	if config != "/etc/apimgr/demo" || data != "/var/lib/apimgr/demo" || logs != "/var/log/apimgr/demo" {
		t.Fatalf("root Linux directories = %q %q %q", config, data, logs)
	}

	values["ProgramData"] = `D:\ProgramData`
	config, data, logs = defaultDirsForOS("demo", "windows", true, `C:\Users\user`, getenv)
	if config != filepath.Join(`D:\ProgramData`, OrgName, "demo") || data != filepath.Join(`D:\ProgramData`, OrgName, "demo", "data") || logs != filepath.Join(`D:\ProgramData`, OrgName, "demo", "logs") {
		t.Fatalf("root Windows directories = %q %q %q", config, data, logs)
	}
	values["ProgramData"] = ""
	config, _, _ = defaultDirsForOS("demo", "windows", true, "", getenv)
	if config != filepath.Join(`C:\ProgramData`, OrgName, "demo") {
		t.Fatalf("default ProgramData config = %q", config)
	}

	values["APPDATA"] = `C:\Roaming`
	values["LOCALAPPDATA"] = `C:\Local`
	config, data, logs = defaultDirsForOS("demo", "windows", false, `C:\Users\user`, getenv)
	if config != filepath.Join(`C:\Roaming`, OrgName, "demo") || data != filepath.Join(`C:\Local`, OrgName, "demo") || logs != filepath.Join(`C:\Local`, OrgName, "demo", "logs") {
		t.Fatalf("user Windows directories = %q %q %q", config, data, logs)
	}
	values["APPDATA"], values["LOCALAPPDATA"] = "", ""
	config, data, _ = defaultDirsForOS("demo", "windows", false, `C:\Users\user`, getenv)
	if config != filepath.Join(`C:\Users\user`, "AppData", "Roaming", OrgName, "demo") || data != filepath.Join(`C:\Users\user`, "AppData", "Local", OrgName, "demo") {
		t.Fatalf("default user Windows directories = %q %q", config, data)
	}

	config, data, logs = defaultDirsForOS("demo", "darwin", false, "/Users/user", getenv)
	if config != "/Users/user/.config/apimgr/demo" || data != "/Users/user/.local/share/apimgr/demo" || logs != "/Users/user/.local/share/apimgr/demo/logs" {
		t.Fatalf("user macOS directories = %q %q %q", config, data, logs)
	}

	values["XDG_CONFIG_HOME"], values["XDG_DATA_HOME"] = "/tmp/config", "/tmp/data"
	config, data, logs = defaultDirsForOS("demo", "linux", false, "/home/user", getenv)
	if config != "/tmp/config/apimgr/demo" || data != "/tmp/data/apimgr/demo" || logs != "/tmp/data/apimgr/demo/logs" {
		t.Fatalf("XDG Linux directories = %q %q %q", config, data, logs)
	}
	values["XDG_CONFIG_HOME"], values["XDG_DATA_HOME"] = "", ""
	config, data, logs = defaultDirsForOS("demo", "linux", false, "/home/user", getenv)
	if config != "/home/user/.config/apimgr/demo" || data != "/home/user/.local/share/apimgr/demo" || logs != "/home/user/.local/share/apimgr/demo/logs" {
		t.Fatalf("default Linux directories = %q %q %q", config, data, logs)
	}
}
