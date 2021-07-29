import Amqp from 'k6/x/amqp';
import Queue from 'k6/x/amqp/queue';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  const queueName = 'K6 queue'
  const exchangeName = 'K6 exchange'

  Queue.unbind({
    queue_name: queueName,
    routing_key: '',
    exchange_name: exchangeName,
    args: null
  })

  console.log(queueName + " queue unbinded from " + exchangeName)
}
