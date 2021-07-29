import Amqp from 'k6/x/amqp';
import Queue from 'k6/x/amqp/queue';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  const queueName = 'K6 queue'
  console.log('Inspecting ' + queueName)
  console.log(JSON.stringify(Queue.inspect(queueName), null, 2))
}
