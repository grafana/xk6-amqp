import Amqp from 'k6/x/amqp';
import Queue from 'k6/x/amqp/queue';

export default function () {
  console.log("K6 amqp extension enabled, version: " + Amqp.version)
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  const queueName = 'K6 queue'
  const consumerName = 'K6 consumer'

  Queue.declare({
    name: queueName,
    durable: false,
    delete_when_unused: false,
    exclusive: false,
    no_wait: false,
    args: null
  })

  console.log(queueName + " queue is ready")

  const publish = function(mark) {
    Amqp.publish({
      queue_name: queueName,
      exchange: '',
      mandatory: false,
      immediate: false,
      content_type: "text/plain",
      body: "Ping from k6 -> " + mark
    })
  }

  publish('A')
  publish('B')
  publish('C')

  const listener = function(data) { console.log('received data: ' + data) }
  Amqp.listen({
    queue_name: queueName,
    listener: listener,
    auto_ack: true,
    consumer: consumerName,
    // exclusive: false,
		// no_local: false,
		// no_wait: false,
    // args: null
  })
}
