package file

import "fmt"

type UploadFileCallback interface {
	Open(size int64)                                                    // file opening size
	PartSize(partSize int64, num int)                                   // shard size, number
	HashPartProgress(index int, size int64, partHash string)            // hash value of each shard
	HashPartComplete(partsHash string, fileHash string)                 // sharding is complete, server marks hash and file final hash
	UploadID(uploadID string)                                           // upload ID
	UploadPartComplete(index int, partSize int64, partHash string)      // upload shard progress
	UploadComplete(fileSize int64, streamSize int64, storageSize int64) // overall progress
	Complete(size int64, url string, typ int)                           // upload completed
}

type emptyUploadCallback struct{}

func (e emptyUploadCallback) Open(size int64) {
	fmt.Println("Callback Open:", size)
}

func (e emptyUploadCallback) PartSize(partSize int64, num int) {
	fmt.Println("Callback PartSize:", partSize, num)
}

func (e emptyUploadCallback) HashPartProgress(index int, size int64, partHash string) {
	//fmt.Println("Callback HashPartProgress:", index, size, partHash)
}

func (e emptyUploadCallback) HashPartComplete(partsHash string, fileHash string) {
	fmt.Println("Callback HashPartComplete:", partsHash, fileHash)
}

func (e emptyUploadCallback) UploadID(uploadID string) {
	fmt.Println("Callback UploadID:", uploadID)
}

func (e emptyUploadCallback) UploadPartComplete(index int, partSize int64, partHash string) {
	fmt.Println("Callback UploadPartComplete:", index, partSize, partHash)
}

func (e emptyUploadCallback) UploadComplete(fileSize int64, streamSize int64, storageSize int64) {
	fmt.Println("Callback UploadComplete:", fileSize, streamSize, storageSize)
}

func (e emptyUploadCallback) Complete(size int64, url string, typ int) {
	fmt.Println("Callback Complete:", size, url, typ)
}
