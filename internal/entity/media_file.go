package entity

type MediaFile struct {
	ID     int64
	Link   string // Ссылка на фото в магазине
	TypeID int64
	CardID int64
}

type MediaFileTypeEnum struct {
	ID   int64
	Type string
}
