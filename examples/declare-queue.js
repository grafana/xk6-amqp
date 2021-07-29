import Amqp from 'k6/x/amqp';
import Queue from 'k6/x/amqp/queue';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  const queueName = 'K6 queue'
  
  Queue.declare({
    name: queueName,
    durable: false,
    delete_when_unused: false,
    exclusive: false,
    no_wait: false,
    args: null
  })

  console.log(queueName + " queue declared")
}
