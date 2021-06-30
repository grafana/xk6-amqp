import Amqp from 'k6/x/amqp';

export default function () {
  console.log("K6 amqp extension enabled, version: " + Amqp.version)
}
