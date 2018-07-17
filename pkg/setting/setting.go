package setting

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-clog/clog"
	"github.com/go-ini/ini"
	"github.com/henson/proxypool/pkg/util"
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
	clog.New(clog.CONSOLE, clog.ConsoleConfig{})

	var err error
	if AppPath, err = execPath(); err != nil {
		clog.Fatal(2, "Fail to get app path: %v\n", err)
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
		clog.Fatal(2, "Do not use '\\' or '\\\\' in paths, instead, please use '/' in all places")
	}
}

// NewContext initializes configuration context.
// NOTE: do not print any log except error.
func NewContext() {
	workDir, err := WorkDir()
	if err != nil {
		clog.Fatal(2, "Fail to get work directory: %v", err)
	}
	ConfFile = path.Join(workDir, "conf/app.ini")

	//Cfg, err = ini.Load("conf/example_app.ini")
	Cfg, err = ini.Load(ConfFile)
	if err != nil {
		clog.Fatal(2, "Fail to parse %s: %v", ConfFile, err)
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

func newLogService() {
	// Because we always create a console logger as primary logger before all settings are loaded,
	// thus if user doesn't set console logger, we should remove it after other loggers are created.
	hasConsole := false

	// Get the log mode.
	if DebugMode {
		LogModes = strings.Split("console", ",")
	} else {
		LogModes = strings.Split(Cfg.Section("log").Key("MODE").MustString("console"), ",")
	}
	LogConfigs = make([]interface{}, len(LogModes))
	levelNames := map[string]clog.LEVEL{
		"trace": clog.TRACE,
		"info":  clog.INFO,
		"warn":  clog.WARN,
		"error": clog.ERROR,
		"fatal": clog.FATAL,
	}
	for i, mode := range LogModes {
		mode = strings.ToLower(strings.TrimSpace(mode))
		sec, err := Cfg.GetSection("log." + mode)
		if err != nil {
			clog.Fatal(2, "Unknown logger mode: %s", mode)
		}

		validLevels := []string{"trace", "info", "warn", "error", "fatal"}
		name := Cfg.Section("log." + mode).Key("LEVEL").Validate(func(v string) string {
			v = strings.ToLower(v)
			if util.IsSliceContainsStr(validLevels, v) {
				return v
			}
			return "trace"
		})
		level := levelNames[name]

		// Generate log configuration.
		switch clog.MODE(mode) {
		case clog.CONSOLE:
			hasConsole = true
			LogConfigs[i] = clog.ConsoleConfig{
				Level:      level,
				BufferSize: Cfg.Section("log").Key("BUFFER_LEN").MustInt64(100),
			}

		case clog.FILE:
			logPath := path.Join(LogRootPath, "ProxyPool.log")
			if err = os.MkdirAll(path.Dir(logPath), os.ModePerm); err != nil {
				clog.Fatal(2, "Fail to create log directory '%s': %v", path.Dir(logPath), err)
			}

			LogConfigs[i] = clog.FileConfig{
				Level:      level,
				BufferSize: Cfg.Section("log").Key("BUFFER_LEN").MustInt64(100),
				Filename:   logPath,
				FileRotationConfig: clog.FileRotationConfig{
					Rotate:   sec.Key("LOG_ROTATE").MustBool(true),
					Daily:    sec.Key("DAILY_ROTATE").MustBool(true),
					MaxSize:  1 << uint(sec.Key("MAX_SIZE_SHIFT").MustInt(28)),
					MaxLines: sec.Key("MAX_LINES").MustInt64(1000000),
					MaxDays:  sec.Key("MAX_DAYS").MustInt64(7),
				},
			}
		}

		clog.New(clog.MODE(mode), LogConfigs[i])
		clog.Trace("Log Mode: %s (%s)", strings.Title(mode), strings.Title(name))
	}

	// Make sure everyone gets version info printed.
	clog.Info("%s %s", AppName, AppVer)
	if !hasConsole {
		clog.Delete(clog.CONSOLE)
	}
}

// NewServices .
func NewServices() {
	newLogService()
}
