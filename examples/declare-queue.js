import Amqp from 'k6/x/amqp';
import Queues from 'k6/x/amqp/queues';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  const queueName = 'K6 queue'
  
  Queues.declare({
    name: queueName,
    durable: false,
    delete_when_unused: false,
    exclusive: false,
    no_wait: false,
    args: null
  })

  console.log(queueName + " queue declared")
}
