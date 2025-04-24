package ports

type IPublisher interface {
	Publish(key string, data []byte, topic string)
}
