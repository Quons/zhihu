package course_handler

type Courser interface {
	Speek(word string) (s string)
}

type TestCourse struct {
}

func (t *TestCourse) Speek(word string) (s string) {
	return s
}
func GetTest(c Courser) (s string) {
	return c.Speek("hello")
}
