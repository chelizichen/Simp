package main

import (
	"Simp/src/events"
	"fmt"
	"time"
)

func main() {
	emitter := events.NewEventEmitter()

	// 监听事件
	emitter.On("message", func(event events.Event) {
		fmt.Printf("Received message: %s with data %v\n", event.Name, event.Data)
	})

	// 触发事件
	emitter.Emit("message", "Hello, World!")

	// 移除事件监听
	emitter.Off("message")

	// 再次触发事件，由于监听器已被移除，所以这次不会有输出
	emitter.Emit("message", "Hello again!")
	time.Sleep(time.Second * 4)
}
