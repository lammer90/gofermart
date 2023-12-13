package userstorage

type UserRepository interface {
	Save(login, authHash string) (err error)
	Find(login string) (authHash string, err error)
}
