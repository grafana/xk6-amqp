import Amqp from 'k6/x/amqp'
import exec from 'k6/execution'

export const options = {
  vus: 10,
  duration: '30s',
}

const url = "amqp://guest:guest@localhost:5672/"
const connIds = new Map()

function getConnectionId (vuId) {
  if (!connIds.has(vuId)) {
    const connectionId = Amqp.start({ connection_url: url })
    connIds.set(vuId, connectionId)
    return connectionId
  }
  return connIds.get(vuId)
}

export default function () {
  console.log("K6 amqp extension enabled, version: " + Amqp.version)

  const connectionId = getConnectionId(exec.vu.idInInstance)
  const queueName = 'K6 queue'

  Amqp.publish({
    connection_id: connectionId,
    queue_name: queueName,
    exchange: '',
    content_type: 'text/plain',
    body: 'Ping from k6'
  })
}
