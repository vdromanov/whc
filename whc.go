package whc

type ExternalOption struct {
	Name  string
	Value string
}

type ExternalOptionsSet interface {
	CheckExist(options []*ExternalOption) error
	PrintRequired(options []*ExternalOption)
	FillValues(options []*ExternalOption)
}

type DataStorage struct {
	URI  string
	Auth string
	Conn interface{}
}

type DataStorageConn interface {
	Connect(uri, auth string) (*DataStorage, error)
}

const (
	ShortDateFmt = "2 January"
	DateFmt      = "02.01.2006"
	TimeFmt      = "15:04"
	DatetimeFmt  = "2006-01-02 15:04"
)
