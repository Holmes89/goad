package mail

type MailRepo interface {
	GetMail(username string, id string) (*Mail, error)
	GetAllMail(username string, all bool) ([]*Mail, error)
	CreateMail(*Mail) error
	UpdateMail(*Mail) error
	DeleteMail(id string, username string) error
}
