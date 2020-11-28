package setting

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-ini/ini"
	clog "unknwon.dev/clog/v2"
)

var (
	//App settings
	AppVer  string
	AppName string
	AppURL  string
	AppPath string
	AppAddr string
	AppPort string

	//Global setting objects
	Cfg       *ini.File
	DebugMode bool
	IsWindows bool
	ConfFile  string

	// Database settings
	UseSQLite3    bool
	UseMySQL      bool
	UsePostgreSQL bool
	UseMSSQL      bool

	// Log settings
	LogRootPath string
	LogModes    []string
	LogConfigs  []interface{}

	//Security settings
	InstallLock bool // true mean installed

	// OAuth2
	SessionExpires time.Duration
)

// execPath returns the executable path.
func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}

func init() {
	IsWindows = runtime.GOOS == "windows"
	var err error
	if AppPath, err = execPath(); err != nil {
		clog.Fatal("Fail to get app path: %v\n", err)
	}

	// Note: we don't use path.Dir here because it does not handle case
	//	which path starts with two "/" in Windows: "//psf/Home/..."
	AppPath = strings.Replace(AppPath, "\\", "/", -1)
}

// WorkDir returns absolute path of work directory.
func WorkDir() (string, error) {
	wd := os.Getenv("ALIGN_WORK_DIR")
	if len(wd) > 0 {
		return wd, nil
	}

	i := strings.LastIndex(AppPath, "/")
	if i == -1 {
		return AppPath, nil
	}
	return AppPath[:i], nil
}

func forcePathSeparator(path string) {
	if strings.Contains(path, "\\") {
		clog.Fatal("Do not use '\\' or '\\\\' in paths, instead, please use '/' in all places")
	}
}

// NewContext initializes configuration context.
// NOTE: do not print any log except error.
func NewContext() {
	workDir, err := WorkDir()
	if err != nil {
		clog.Fatal("Fail to get work directory: %v", err)
	}
	ConfFile = path.Join(workDir, "conf/app.ini")

	//Cfg, err = ini.Load("conf/example_app.ini")
	Cfg, err = ini.Load(ConfFile)
	if err != nil {
		clog.Fatal("Fail to parse %s: %v", ConfFile, err)
	}

	Cfg.NameMapper = ini.AllCapsUnderscore

	// Load security config
	InstallLock = Cfg.Section("security").Key("INSTALL_LOCK").MustBool(false)

	// Load server config
	sec := Cfg.Section("server")
	AppName = Cfg.Section("").Key("APP_NAME").MustString("ProxyPool")
	AppURL = sec.Key("ROOT_URL").MustString("http://localhost:3000/")
	if AppURL[len(AppURL)-1] != '/' {
		AppURL += "/"
	}
	AppAddr = sec.Key("HTTP_ADDR").MustString("0.0.0.0")
	AppPort = sec.Key("HTTP_PORT").MustString("3001")
	SessionExpires = sec.Key("SESSION_EXPIRES").MustDuration(time.Hour * 24 * 7)
}

func NewLogService() {
	err := clog.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}
