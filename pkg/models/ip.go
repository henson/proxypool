package models

// IP struct
type IP struct {
	ID   int64  `xorm:"pk autoincr" json:"-"`
	Data string `xorm:"NOT NULL" json:"ip"`
	Type string `xorm:"NOT NULL" json:"type"`
}

// NewIP .
func NewIP() *IP {
	return &IP{}
}

// SaveIps save ips info to database
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

	count, _ := x.Where("id> ?", 1).Count(new(IP))
	return count
}

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

func GetOne(ip string) *IP {
	return getOne(ip)
}

func getAll() ([]*IP, error) {
	tmpIp := make([]*IP, 0)
	err := x.Find(&tmpIp)
	if err != nil {
		return nil, err
	}
	return tmpIp, nil
}

func GetAll() ([]*IP, error) {
	return getAll()
}

func findAll(value string) ([]*IP, error) {
	tmpIp := make([]*IP, 0)
	err := x.Where("type=?", value).Find(&tmpIp)
	if err != nil {
		return nil, err
	}
	return tmpIp, nil
}

func FindAll(value string) ([]*IP, error) {
	return findAll(value)
}

func update(ip IP) error {
	_, err := x.Id(1).Update(ip)
	if err != nil {
		return err
	}
	return nil
}

func Update(ip IP) error {
	return update(ip)
}
