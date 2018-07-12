package models

// IP struct
type IP struct {
	ID    int64  `xorm:"pk autoincr" json:"-"`
	Data  string `xorm:"NOT NULL" json:"ip"`
	Type1 string `xorm:"NOT NULL" json:"type1"`
	Type2 string `xorm:"NULL" json:"type2,omitempty"`
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

	count, _ := x.Where("id> ?", 0).Count(new(IP))
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
	switch value {
	case "http":
		err := x.Where("type1=?", "http").Find(&tmpIp)
		if err != nil {
			return tmpIp, err
		}
	case "https":
		err := x.Where("type2=?", "https").Find(&tmpIp)
		if err != nil {
			return tmpIp, err
		}
	default:
		return tmpIp, nil
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
