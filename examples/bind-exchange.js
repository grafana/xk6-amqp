import Amqp from 'k6/x/amqp';
import Exchanges from 'k6/x/amqp/exchanges';

export default function () {
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  
  const sourceExchangeName = 'K6 exchange'
  const destinationExchangeName = 'destination K6 exchange'

  Exchanges.bind({
    destination_exchange_name: destinationExchangeName,
    routing_key: '',
    source_exchange_name: sourceExchangeName,
    no_wait: false,
    args: null
  })

  console.log(destinationExchangeName + " exchange binded to " + sourceExchangeName + ' exchange')
}
