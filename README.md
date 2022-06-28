> ### ⚠️ This is a proof of concept
>
> As this is a proof of concept,  it won't be supported by the k6 team.
> It may also break in the future as xk6 evolves. USE AT YOUR OWN RISK!
> Any issues with the tool should be raised [here](https://github.com/grafana/xk6-amqp/issues).

</br>
</br>

<div align="center">

# xk6-amqp
AMQP xk6 plugin. Built for [k6](https://go.k6.io/k6) using [xk6](https://github.com/grafana/xk6).

</div>

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download `xk6`:
  ```bash
  $ go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```bash
  $ xk6 build --with github.com/grafana/xk6-amqp@latest
  ```
## Development

To make development a little smoother, you may run the build script provided in the root folder. It will create a k6 binary with your local code rather than from GitHub.

```bash
$ ./build.sh && ./k6 run my-test-script.js
```


## Example

```javascript
import Amqp from 'k6/x/amqp';
import Queue from 'k6/x/amqp/queue';

export default function () {
  console.log("K6 amqp extension enabled, version: " + Amqp.version)
  const url = "amqp://guest:guest@localhost:5672/"
  Amqp.start({
    connection_url: url
  })
  console.log("Connection opened: " + url)
  
  const queueName = 'K6 general'
  
  Queue.declare({
    name: queueName,
    // durable: false,
    // delete_when_unused: false,
    // exclusive: false,
    // no_wait: false,
    // args: null
  })

  console.log(queueName + " queue is ready")

  Amqp.publish({
    queue_name: queueName,
    body: "Ping from k6",
    content_type: "text/plain"
    // exchange: '',
    // mandatory: false,
    // immediate: false,
  })

  const listener = function(data) { console.log('received data: ' + data) }
  Amqp.listen({
    queue_name: queueName,
    listener: listener,
    // consumer: '',
    // auto_ack: true,
    // exclusive: false,
		// no_local: false,
		// no_wait: false,
    // args: null
  })
}

```

Result output:

```plain
$ ./k6 run script.js

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: ../xk6-amqp/examples/test.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)

INFO[0000] K6 amqp extension enabled, version: v0.0.1    source=console
INFO[0000] Connection opened: amqp://guest:guest@localhost:5672/  source=console
INFO[0000] K6 general queue is ready                     source=console
INFO[0000] received data: Ping from k6                   source=console

running (00m00.0s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m00.0s/10m0s  1/1 iters, 1 per VU

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=31.37ms min=31.37ms med=31.37ms max=31.37ms p(90)=31.37ms p(95)=31.37ms
     iterations...........: 1   30.855627/s

```

Inspect examples folder for more details.