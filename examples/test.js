import Amqp from 'k6/x/amqp';

export default function () {
  console.log("K6 amqp extension enabled, version: " + Amqp.version)
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  console.log("Connection opened: " + url)
  Amqp.publish({
    queue_name: "general",
    body: "Hello AMQP from xk6"
  })
}
