package test

import "fmt"

type TestSendMsg struct {
}

func (TestSendMsg) OnSuccess(data string) {
	fmt.Println("testSendImg, OnSuccess, output: ", data)
}

func (TestSendMsg) OnError(code int32, msg string) {
	fmt.Println("testSendImg, OnError, ", code, msg)
}

func (TestSendMsg) OnProgress(progress int) {
	fmt.Println("progress: ", progress)
}
