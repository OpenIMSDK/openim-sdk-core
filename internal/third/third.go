package third

import (
	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"
	"sync"
)

type Third struct {
	platformID    int32
	loginUserID   string
	systemType    string
	LogFilePath   string
	fileUploader  *file.File
	logUploadLock sync.Mutex
}

func NewThird(platformID int32, loginUserID, systemType, LogFilePath string, fileUploader *file.File) *Third {
	return &Third{platformID: platformID, loginUserID: loginUserID, systemType: systemType, LogFilePath: LogFilePath, fileUploader: fileUploader}
}
