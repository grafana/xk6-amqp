import Amqp from 'k6/x/amqp';
import Queues from 'k6/x/amqp/queues';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  const queueName = 'K6 queue'
  const exchangeName = 'K6 exchange'

  Queues.bind({
    queue_name: queueName,
    routing_key: '',
    exchange_name: exchangeName,
    no_wait: false,
    args: null
  })

  console.log(queueName + " queue binded to " + exchangeName)
}
