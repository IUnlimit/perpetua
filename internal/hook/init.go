package hook

import (
	"github.com/IUnlimit/perpetua/internal/utils"
)

func Init() {
	aesKey := []byte("perpetua-aes-key-field-github-00")
	decryptAES := utils.DecryptAES(aesKey, TOKEN)
	headers["Authorization"] = "Bearer " + decryptAES
}
