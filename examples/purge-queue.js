import Amqp from 'k6/x/amqp';
import Queues from 'k6/x/amqp/queues';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  const queueName = 'K6 queue'
  const count = Queues.purge(queueName, false)

  console.log(queueName + " purge: " + count + " messages deleted")
}
