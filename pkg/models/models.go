package models

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/go-clog/clog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/henson/proxypool/pkg/setting"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

// Engine represents a XORM engine or session.
type Engine interface {
	Delete(interface{}) (int64, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Exist(...interface{}) (bool, error)
	Find(interface{}, ...interface{}) error
	Get(interface{}) (bool, error)
	Id(interface{}) *xorm.Session
	In(string, ...interface{}) *xorm.Session
	Insert(...interface{}) (int64, error)
	InsertOne(interface{}) (int64, error)
	Iterate(interface{}, xorm.IterFunc) error
	Query(string, ...interface{}) (sql.Result, error)
	Sql(string, ...interface{}) *xorm.Session
	Table(interface{}) *xorm.Session
	Where(interface{}, ...interface{}) *xorm.Session
}

var (
	x      *xorm.Engine
	tables []interface{}
	// HasEngine .
	HasEngine bool

	DbCfg struct {
		Type, Host, Name, User, Passwd, Path, SSLMode string
	}
	// EnableSQLite3 bool
)

func init() {
	tables = append(tables,
		new(IP))
	gonicNames := []string{"SSL"}
	for _, name := range gonicNames {
		core.LintGonicMapper[name] = true
	}

}

// LoadDatabaseInfo .
func LoadDatabaseInfo() {
	sec := setting.Cfg.Section("database")
	DbCfg.Type = sec.Key("DB_TYPE").String()
	switch DbCfg.Type {
	case "sqlite":
		setting.UseSQLite = true
		// EnableSQLite3 = true
	case "mysql":
		setting.UseMySQL = true
	case "postgres":
		setting.UsePostgreSQL = true
	case "mssql":
		setting.UseMSSQL = true
	}
	DbCfg.Host = sec.Key("HOST").String()
	DbCfg.Name = sec.Key("NAME").String()
	DbCfg.User = sec.Key("USER").String()
	if len(DbCfg.Passwd) == 0 {
		DbCfg.Passwd = sec.Key("PASSWD").String()
	}
	DbCfg.SSLMode = sec.Key("SSL_MODE").String()
	DbCfg.Path = sec.Key("PATH").MustString("data/ProxyPool.db")
}

// parsePostgreSQLHostPort parses given input in various forms defined in
// https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING
// and returns proper host and port number.
func parsePostgreSQLHostPort(info string) (string, string) {
	host, port := "127.0.0.1", "5432"
	if strings.Contains(info, ":") && !strings.HasSuffix(info, "]") {
		idx := strings.LastIndex(info, ":")
		host = info[:idx]
		port = info[idx+1:]
	} else if len(info) > 0 {
		host = info
	}
	return host, port
}

func parseMSSQLHostPort(info string) (string, string) {
	host, port := "127.0.0.1", "1433"
	if strings.Contains(info, ":") {
		host = strings.Split(info, ":")[0]
		port = strings.Split(info, ":")[1]
	} else if strings.Contains(info, ",") {
		host = strings.Split(info, ",")[0]
		port = strings.TrimSpace(strings.Split(info, ",")[1])
	} else if len(info) > 0 {
		host = info
	}
	return host, port
}

func getEngine() (*xorm.Engine, error) {
	connStr := ""
	Param := "?"
	if strings.Contains(DbCfg.Name, Param) {
		Param = "&"
	}
	switch DbCfg.Type {
	case "mysql":
		if DbCfg.Host[0] == '/' { // looks like a unix socket
			connStr = fmt.Sprintf("%s:%s@unix(%s)/%s%scharset=utf8&parseTime=true",
				DbCfg.User, DbCfg.Passwd, DbCfg.Host, DbCfg.Name, Param)
		} else {
			connStr = fmt.Sprintf("%s:%s@tcp(%s)/%s%scharset=utf8&parseTime=true",
				DbCfg.User, DbCfg.Passwd, DbCfg.Host, DbCfg.Name, Param)
		}
	case "postgres":
		host, port := parsePostgreSQLHostPort(DbCfg.Host)
		if host[0] == '/' { // looks like a unix socket
			connStr = fmt.Sprintf("postgres://%s:%s@:%s/%s%ssslmode=%s&host=%s",
				url.QueryEscape(DbCfg.User), url.QueryEscape(DbCfg.Passwd), port, DbCfg.Name, Param, DbCfg.SSLMode, host)
		} else {
			connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s%ssslmode=%s",
				url.QueryEscape(DbCfg.User), url.QueryEscape(DbCfg.Passwd), host, port, DbCfg.Name, Param, DbCfg.SSLMode)
		}
	case "mssql":
		host, port := parseMSSQLHostPort(DbCfg.Host)
		connStr = fmt.Sprintf("server=%s; port=%s; database=%s; user id=%s; password=%s;", host, port, DbCfg.Name, DbCfg.User, DbCfg.Passwd)
	case "sqlite":
		if err := os.MkdirAll(path.Dir(DbCfg.Path), os.ModePerm); err != nil {
			return nil, fmt.Errorf("Fail to create directories: %v", err)
		}
		connStr = "file:" + DbCfg.Path + "?cache=shared&mode=rwc"
	default:
		return nil, fmt.Errorf("Unknown database type: %s", DbCfg.Type)
	}
	return xorm.NewEngine(DbCfg.Type, connStr)
}

// NewTestEngine .
func NewTestEngine(x *xorm.Engine) (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("Connect to database: %v", err)
	}

	x.SetMapper(core.GonicMapper{})
	return x.StoreEngine("InnoDB").Sync2(tables...)
}

// SetEngine .
func SetEngine() (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("Fail to connect to database: %v", err)
	}

	x.SetMapper(core.GonicMapper{})

	// WARNING: for serv command, MUST remove the output to os.stdout,
	// so use log file to instead print to stdout.
	sec := setting.Cfg.Section("log.xorm")
	logger, err := clog.NewFileWriter(path.Join(setting.LogRootPath, "xorm.log"),
		clog.FileRotationConfig{
			Rotate:  sec.Key("ROTATE").MustBool(true),
			Daily:   sec.Key("ROTATE_DAILY").MustBool(true),
			MaxSize: sec.Key("MAX_SIZE").MustInt64(100) * 1024 * 1024,
			MaxDays: sec.Key("MAX_DAYS").MustInt64(3),
		})
	if err != nil {
		return fmt.Errorf("Fail to create 'xorm.log': %v", err)
	}

	if !setting.DebugMode {
		x.SetLogger(log.NewSimpleLogger3(logger, log.DEFAULT_LOG_PREFIX, log.DEFAULT_LOG_FLAG, log.LOG_INFO))
	} else {
		x.SetLogger(log.NewSimpleLogger3(logger, log.DEFAULT_LOG_PREFIX, log.DEFAULT_LOG_FLAG, log.LOG_DEBUG))
	}
	x.ShowSQL(true)
	return nil
}

// NewEngine .
func NewEngine() (err error) {
	if err = SetEngine(); err != nil {
		return err
	}

	if err = x.StoreEngine("InnoDB").Sync2(tables...); err != nil {
		return fmt.Errorf("sync database struct error: %v", err)
	}

	return nil
}
