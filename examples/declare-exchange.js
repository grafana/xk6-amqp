import Amqp from 'k6/x/amqp';
import Exchange from 'k6/x/amqp/exchange';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  console.log("Connection opened: " + url)

  const exchangeName = 'K6 exchange'

  Exchange.declare({
    name: exchangeName,
  	kind: 'direct',
    durable: false,
    auto_delete: false,
    internal: false,
    no_wait: false,
	  args: null
  })

  console.log(exchangeName + " exchange is ready")
}
