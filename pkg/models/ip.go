package models

// IP struct
type IP struct {
	ID    int64  `xorm:"pk autoincr" json:"-"`
	Data  string `xorm:"NOT NULL" json:"ip"`
	Type1 string `xorm:"NOT NULL" json:"type1"`
	Type2 string `xorm:"NULL" json:"type2,omitempty"`
	Speed int64  `xorm:"NOT NULL" json:"speed,omitempty"`
}

// NewIP .
func NewIP() *IP {
	return &IP{Speed: 999999}
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

	count, _ := x.Where("id> ?", 0).Count(new(IP))
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
		err := x.Where("speed <= 1000 and type2=?", "https").Find(&tmpIp)
		if err != nil {
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

func update(ip IP) error {
	_, err := x.Id(1).Update(ip)
	if err != nil {
		return err
	}
	return nil
}

// Update .
func Update(ip IP) error {
	return update(ip)
}
