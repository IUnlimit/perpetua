package perp

func Bootstrap() {
	Configure()
	go Start()
	EnableAgent()
}
