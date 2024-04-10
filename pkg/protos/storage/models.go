package storage

// FileAttrs - Атрибуты файла (редактируемые)
type FileAttrs struct {
	TTL       string // Количество дней в течение которого файл будет храниться (Необязательный)
	Filename  string // Наименование файла (Необязательный)
	Type      string // Бизнес-тип (Необязательный)
	SubType   string // Бизнес-подтип (Необязательный)
	Readonly  bool   // Флаг доступности файла только для чтения (Необязательный)
	Protected bool   // Флаг, включающий блокировку от удаления файла (Необязательный)
}

// AllFileAttrs - Атрибуты файла (редактируемые + сгенерированные)
type AllFileAttrs struct {
	Created     int64  // Дата создания (Unix)
	Expires     int64  // Срок хранения (Unix)
	Creator     string // Создатель
	CustomId    string // Кастомный Id, заданный создателем
	Filename    string // Имя файла
	Size        int64  // Размер файла в байтах
	StorageType int32  // Тип хранилища (0 - горячее, 1 - холодное)
	Type        string // Бизнес-тип
	SubType     string // Бизнес-подтип
	Readonly    bool   // Флаг доступности файла только для чтения
	Protected   bool   // Флаг, включающий блокировку от удаления файла
}

// ServiceAttrs - Атрибуты, установленные конкретным сервисом
type ServiceAttrs struct {
	TTL int64 // Срок хранения (Unix)
}
