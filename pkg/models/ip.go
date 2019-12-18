package models

import (
	"fmt"
	"time"
)

// IP struct
type IP struct {
	ID         int64     `xorm:"pk autoincr" json:"-"`
	Data       string    `xorm:"NOT NULL" json:"ip"`
	Type1      string    `xorm:"NOT NULL" json:"type1"`
	Type2      string    `xorm:"NULL" json:"type2,omitempty"`
	Speed      int64     `xorm:"NOT NULL" json:"speed,omitempty"`
	CreateTime time.Time `xorm:"NOT NULL" json:"-"`
	UpdateTime time.Time `xorm:"NOT NULL" json:"-"`
}

// NewIP .
func NewIP() *IP {
	//init the speed to 100 Sec
	return &IP{
		Speed:      100,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}

//InsertIps SaveIps save ips info to database
func InsertIps(ip *IP) (err error) {

	ses := x.NewSession()
	defer ses.Close()
	if err := ses.Begin(); err != nil {
		return err
	}
	if _, err = ses.Insert(ip); err != nil {
		return err
	}

	return ses.Commit()
}

func countIps() int64 {
	// set id >= 0, fix bug: when this is nothing in the database
	count, _ := x.Where("id>= ?", 0).Count(new(IP))
	return count
}

// CountIPs .
func CountIPs() int64 {
	return countIps()
}

func deleteIP(ip *IP) error {
	_, err := x.Delete(ip)
	if err != nil {
		return err
	}
	return nil
}

// DeleteIP .
func DeleteIP(ip *IP) error {
	return deleteIP(ip)
}

func getOne(ip string) *IP {
	var tmpIp IP
	result, _ := x.Where("data=?", ip).Get(tmpIp)
	if result {
		return &tmpIp
	}

	return NewIP()

}

// GetOne .
func GetOne(ip string) *IP {
	return getOne(ip)
}

func getAll() ([]*IP, error) {
	tmpIp := make([]*IP, 0)

	err := x.Where("speed <= 1000").Find(&tmpIp)
	if err != nil {
		return nil, err
	}
	return tmpIp, nil
}

// GetAll .
func GetAll() ([]*IP, error) {
	return getAll()
}

func findAll(value string) ([]*IP, error) {
	tmpIp := make([]*IP, 0)
	switch value {
	case "http":
		err := x.Where("speed <= 1000 and type1=?", "http").Find(&tmpIp)
		if err != nil {
			return tmpIp, err
		}
	case "https":
		//test has https proxy on databases or not
		HasHttps := TestHttps()
		if HasHttps == false {
			return tmpIp, nil
		}
		err := x.Where("speed <= 1000 and type1=?", "https").Find(&tmpIp)
		if err != nil {
			fmt.Println(err.Error())
			return tmpIp, err
		}
	default:
		return tmpIp, nil
	}

	return tmpIp, nil
}

// FindAll .
func FindAll(value string) ([]*IP, error) {
	return findAll(value)
}

func update(ip *IP) error {
	tmp := ip
	tmp.UpdateTime = time.Now()
	_, err := x.ID(1).Update(tmp)
	if err != nil {
		return err
	}
	return nil
}

// Update the Ip
func Update(ip *IP) error {
	return update(ip)
}

//Test if have https proxy in database
//just test on mysql database
// dbName: ProxyPool
// dbTableName: ip
// select distinct if(exists(select * from ip where type2='https'),1,0) as a from ip;
func TestHttps() bool {
	has, err := x.Exist(&IP{Type1: "https"})
	if err != nil {
		return false
	}

	return has
}
