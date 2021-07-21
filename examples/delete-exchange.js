import Amqp from 'k6/x/amqp';
import Exchange from 'k6/x/amqp/exchange';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  const exchangeName = 'K6 exchange'

  Exchange.delete(exchangeName)

  console.log(exchangeName + " exchange deleted")
}
