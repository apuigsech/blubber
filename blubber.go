package blubber

type Blubber struct {
	ProviderList	[]Provider
}

func NewBlubber() *Blubber {
	bl := &Blubber{}
	return bl
}