package switch_api

type Switch interface {
	Update(string) error
}
