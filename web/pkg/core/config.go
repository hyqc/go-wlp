package core

type Configure interface {
	GetFilePath() string
	Init() error
	Handle() error
}

func Init(conf Configure) error {
	if err := conf.Init(); err != nil {
		return err
	}
	if err := conf.Handle(); err != nil {
		return err
	}
	return nil
}
