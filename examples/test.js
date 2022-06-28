import Amqp from 'k6/x/amqp';
import Queue from 'k6/x/amqp/queue';

export default function () {
  console.log("K6 amqp extension enabled, version: " + Amqp.version)
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  console.log("Connection opened: " + url)
  
  const queueName = 'K6 general'
  
  Queue.declare({
    name: queueName,
    // durable: false,
    // delete_when_unused: false,
    // exclusive: false,
    // no_wait: false,
    // args: null
  })

  console.log(queueName + " queue is ready")

  Amqp.publish({
    queue_name: queueName,
    body: "Ping from k6",
    content_type: "text/plain"
    // exchange: '',
    // mandatory: false,
    // immediate: false,
  })

  const listener = function(data) { console.log('received data: ' + data) }
  Amqp.listen({
    queue_name: queueName,
    listener: listener,
    auto_ack: true,
    // consumer: '',
    // exclusive: false,
		// no_local: false,
		// no_wait: false,
    // args: null
  })
}
