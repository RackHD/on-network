package switch_operations

type Switch interface {
	Update(string, string) error
}
