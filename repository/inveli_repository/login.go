package invelirepository

type LoginRepositoryInterface interface {
	Login(username string, password string) (string, error)
}
