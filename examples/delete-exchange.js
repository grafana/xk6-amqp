import Amqp from 'k6/x/amqp';
import Exchanges from 'k6/x/amqp/exchanges';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  console.log("Connection opened: " + url)

  const exchangeName = 'K6 exchange'

  Exchanges.delete(exchangeName)

  console.log(exchangeName + " exchange deleted")
}
