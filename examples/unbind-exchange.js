import Amqp from 'k6/x/amqp';
import Exchange from 'k6/x/amqp/exchange';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
   
  const sourceExchangeName = 'K6 exchange'
  const destinationExchangeName = 'destination K6 exchange'

  Exchange.unbind({
    destination_exchange_name: destinationExchangeName,
    routing_key: '',
    source_exchange_name: sourceExchangeName,
    no_wait: false,
    args: null
  })

  console.log(destinationExchangeName + ' exchange unbinded from ' + sourceExchangeName + ' exchange')
}
