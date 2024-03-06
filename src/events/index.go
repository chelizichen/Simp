package events

// 定义事件类型
type Event struct {
	Name string
	Data interface{}
}

// 事件发射器
type EventEmitter struct {
	events map[string]chan Event
}

// 新建一个事件发射器
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		events: make(map[string]chan Event),
	}
}

// 监听事件
func (ee *EventEmitter) On(eventName string, handler func(Event)) {
	if _, ok := ee.events[eventName]; !ok {
		ee.events[eventName] = make(chan Event)
	}

	go func() {
		for event := range ee.events[eventName] {
			handler(event)
		}
	}()
}

// 触发事件
func (ee *EventEmitter) Emit(eventName string, data interface{}) {
	event := Event{Name: eventName, Data: data}
	if channel, ok := ee.events[eventName]; ok {
		channel <- event
	}
}

// 移除事件监听
func (ee *EventEmitter) Off(eventName string) {
	if channel, ok := ee.events[eventName]; ok {
		close(channel)
		delete(ee.events, eventName)
	}
}
