package initial

import (
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	"github.com/henson/proxypool/pkg/models"
	"github.com/henson/proxypool/pkg/setting"
	"github.com/henson/proxypool/pkg/util"
	clog "unknwon.dev/clog/v2"
	"xorm.io/xorm"
)

// GlobalInit is for global configuration reload-able.
func GlobalInit() {
	setting.NewContext()
	setting.NewLogService()
	clog.Trace("Log path: %s", setting.LogRootPath)
	models.LoadDatabaseInfo()

	if setting.InstallLock {
		if err := models.NewEngine(); err != nil {
			clog.Fatal("Fail to initialize ORM engine: %v", err)
		}
		models.HasEngine = true
	}

	// if models.EnableSQLite3 {
	// 	clog.Info("SQLite Supported")
	// }
	if !setting.InstallLock {
		Database()
	}
}

// Database .
func Database() {
	//Set test engine
	var x *xorm.Engine
	if err := models.NewTestEngine(x); err != nil {
		clog.Fatal("Fail to set test ORM engine: %v", err)
	}
	// Save settings.
	cfg := ini.Empty()
	if util.IsFile(setting.ConfFile) {
		// Keeps custom settings if there is already something.
		if err := cfg.Append(setting.ConfFile); err != nil {
			clog.Error("Fail to load conf '%s': %v", setting.ConfFile, err)
		}
	}
	// Save App name
	cfg.Section("").Key("APP_NAME").SetValue(setting.AppName)
	// Save server config
	cfg.Section("server").Key("HTTP_ADDR").SetValue(setting.AppAddr)
	cfg.Section("server").Key("HTTP_PORT").SetValue(setting.AppPort)
	cfg.Section("server").Key("SESSION_EXPIRES").SetValue(setting.SessionExpires.String())
	// Save database config
	cfg.Section("database").Key("DB_TYPE").SetValue(models.DbCfg.Type)
	cfg.Section("database").Key("HOST").SetValue(models.DbCfg.Host)
	cfg.Section("database").Key("NAME").SetValue(models.DbCfg.Name)
	cfg.Section("database").Key("USER").SetValue(models.DbCfg.User)
	cfg.Section("database").Key("PASSWD").SetValue(models.DbCfg.Passwd)
	cfg.Section("database").Key("SSL_MODE").SetValue(models.DbCfg.SSLMode)
	cfg.Section("database").Key("PATH").SetValue(models.DbCfg.Path)
	// Change Installock value to true
	cfg.Section("security").Key("INSTALL_LOCK").SetValue("true")
	// Save log config
	cfg.Section("log").Key("MODE").SetValue("file")
	cfg.Section("log").Key("LEVEL").SetValue("Info")
	cfg.Section("log").Key("BUFFER_LEN").SetValue("100")
	cfg.Section("log").Key("ROOT_PATH").SetValue(setting.LogRootPath)
	// Save file setting
	os.MkdirAll(filepath.Dir(setting.ConfFile), os.ModePerm)
	if err := cfg.SaveTo(setting.ConfFile); err != nil {
		clog.Fatal("[Initial]Save config failed: %v", err)
	}
	clog.Info("[Initial]Initialize database completed.")
}
