package storage

import (
	"github.com/henson/ProxyPool/models"
	"github.com/henson/ProxyPool/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Config 全局配置文件
var Config = util.NewConfig()

// GlobalMgoSession 全局连接Session
var GlobalMgoSession, _ = mgo.Dial(string(Config.Mongo.Addr))

// Storage struct is used for storeing persistent data of alerts
type Storage struct {
	database string
	table    string
	session  *mgo.Session
}

// NewStorage creates and returns new Storage instance
func NewStorage() *Storage {
	return &Storage{database: Config.Mongo.DB, table: Config.Mongo.Table, session: GlobalMgoSession}
}

// GetDBSession returns a new connection from the pool
func (s *Storage) GetDBSession() *mgo.Session {
	return s.session.Clone()
}

// Create insert new item
func (s *Storage) Create(item interface{}) error {
	ses := s.GetDBSession()
	defer ses.Close()
	err := ses.DB(s.database).C(s.table).Insert(item)
	if err != nil {
		return err
	}
	return nil
}

// GetOne Finds and returns one data from storage
func (s *Storage) GetOne(value string) (*models.IP, error) {
	ses := s.GetDBSession()
	defer ses.Close()
	t := models.NewIP()
	err := ses.DB(s.database).C(s.table).Find(bson.M{"data": value}).One(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Count all collections
func (s *Storage) Count() int {
	ses := s.GetDBSession()
	defer ses.Close()
	num, err := ses.DB(s.database).C(s.table).Count()
	if err != nil {
		num = 0
	}
	return num
}

// Delete .
func (s *Storage) Delete(ip *models.IP) error {
	ses := s.GetDBSession()
	defer ses.Close()
	err := ses.DB(s.database).C(s.table).RemoveId(ip.ID)
	if err != nil {
		return err
	}
	return nil
}

// Update .
func (s *Storage) Update(ip *models.IP) error {
	ses := s.GetDBSession()
	defer ses.Close()
	err := ses.DB(s.database).C(s.table).Update(bson.M{"_id": ip.ID}, ip)
	if err != nil {
		return err
	}
	return nil
}

// GetAll .
func (s *Storage) GetAll() ([]*models.IP, error) {
	ses := s.GetDBSession()
	defer ses.Close()
	var ips []*models.IP
	err := ses.DB(s.database).C(s.table).Find(nil).All(&ips)
	if err != nil {
		return nil, err
	}
	return ips, nil
}

// FindAll .
func (s *Storage) FindAll(value string) ([]*models.IP, error) {
	ses := s.GetDBSession()
	defer ses.Close()
	var ips []*models.IP
	err := ses.DB(s.database).C(s.table).Find(bson.M{"type": bson.M{"$regex": value, "$options": "$i"}}).All(&ips)
	if err != nil {
		return nil, err
	}
	return ips, nil
}
