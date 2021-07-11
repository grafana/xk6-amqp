import Amqp from 'k6/x/amqp';
import Queues from 'k6/x/amqp/queues';
import { sleep } from 'k6';

export default function () {
  console.log("K6 amqp extension enabled, version: " + Amqp.version)
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  console.log("Connection opened: " + url)

  const queueName = 'K6 general'
  const consumerName = 'K6 consumer'

  Queues.declare({
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
    consumer: consumerName,
    // auto_ack: true,
    // exclusive: false,
		// no_local: false,
		// no_wait: false,
    // args: null
  })
}
