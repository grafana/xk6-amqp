import Amqp from 'k6/x/amqp';
import { sleep } from 'k6';

export default function () {
  console.log("K6 amqp extension enabled, version: " + Amqp.version)
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  console.log("Connection opened: " + url)

  const listener = function(data) { console.log('received data: ' + data) }
  Amqp.listen({
    queue_name: "general",
    listener: listener
  })
  const publish = function(mark) {
    Amqp.publish({
      queue_name: "general",
      body: "Ping from k6: " + mark
    })
  }

  publish('A')
  sleep('3s')
  publish('B')
  sleep('5s')
  publish('C')
}
