package factory

var (
	_ Factory = &MockEC2Factory{}
	_ Pod     = &MockEC2{}
)

type MockEC2Factory struct {
}

func (m MockEC2Factory) Create(number int) ([]Pod, error) {
	var pods []Pod
	for i := 0; i < number; i++ {
		pods = append(pods, &MockEC2{
			name: string(i),
		})
	}
	return pods, nil
}

type MockEC2 struct {
	name string
	EC2
}

func (m MockEC2) Target() string {
	return "localhost"
}

func (m MockEC2) Name() string {
	return m.name
}

func (m MockEC2) Ready() (bool, error) {
	return true, nil
}

func (m MockEC2) Delete() error {
	return nil
}
