package config

const (
	RabbitAsyncTransEnable = false
	RabbitUrl = "amqp://guest:guest@192.168.47.131:5672/"
	RabbitExchangeName = "uploadserver.trans"
	RabbitCephQueueName = "uploadserver.trans.ceph"
	RabbitCephErrQueueName = "uploadserver.trans.ceph.err"
	RabbitCephRoutingKey = "ceph"
	RabbitCephErrRoutingKey = "ceph.err"
)
