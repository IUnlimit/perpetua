package handle

type Handler struct {
	Redirect func(data []byte)
}
