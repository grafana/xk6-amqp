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
    name: queueName
  })

  console.log(queueName + " queue is ready")

  let body = {
    metadata: {
      header1: "Performance Test Message"
    },
    body: {
      field1: "some value"
    }
  }

  Amqp.publish({
    queue_name: queueName,
    body: JSON.stringify(body),
    content_type: "application/x-msgpack"
  })

  const listener = function(data) { console.log('received data: ' + data) }
  Amqp.listen({
    queue_name: queueName,
    listener: listener,
    auto_ack: true
  })
}
