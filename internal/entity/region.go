package entity

type Region struct {
	ID         int64
	RegionName string
	District   District
	Country    Country
}
