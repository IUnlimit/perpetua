package hook

import "encoding/base64"

func Init() {
	initGithub()
}

func initGithub() {
	decode, _ := base64.StdEncoding.DecodeString(TOKEN)
	headers["Authorization"] = "Bearer " + string(decode)
}
