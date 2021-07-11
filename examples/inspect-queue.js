import Amqp from 'k6/x/amqp';
import Queues from 'k6/x/amqp/queues';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  const queueName = 'K6 general'
  console.log('Inspecting ' + queueName)
  console.log(JSON.stringify(Queues.inspect(queueName), null, 2))
}
