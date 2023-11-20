> ### ⚠️ Deprecated!!!
>
> This extension was originally created as a _proof of concept_.
> At this time, there are no maintainers available to support this extension.
>
> USE AT YOUR OWN RISK!

<br />

# xk6-amqp

A k6 extension for publishing and consuming messages from queues and exchanges.
This project utilizes [AMQP 0.9.1](https://www.rabbitmq.com/tutorials/amqp-concepts.html), the most common AMQP protocol in use today.

> :warning: This project is not compatible with [AMQP 1.0](http://docs.oasis-open.org/amqp/core/v1.0/os/amqp-core-overview-v1.0-os.html).
> A list of AMQP 1.0 brokers and other AMQP 1.0 resources may be found at [github.com/xinchen10/awesome-amqp](https://github.com/xinchen10/awesome-amqp).

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download [xk6](https://github.com/grafana/xk6):
  ```bash
  $ go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. [Build the k6 binary](https://github.com/grafana/xk6#command-usage):
  ```bash
  $ xk6 build --with github.com/grafana/xk6-amqp@latest
  ```

## Development
To make development a little smoother, use the `Makefile` in the root folder. The default target will format your code, run tests, and create a `k6` binary with your local code rather than from GitHub.

```shell
git clone git@github.com:grafana/xk6-amqp.git
cd xk6-amqp
make
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
    // timestamp: Math.round(Date.now() / 1000)
    // exchange: '',
    // mandatory: false,
    // immediate: false,
    // headers: {
    //   'header-1': '',
    // },
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

# Testing Locally

This repository includes a [docker-compose.yml](./docker-compose.yml) file that starts RabbitMQ with Management Plugin for testing the extension locally.

> :warning: This environment is intended for testing only and should not be used for production purposes.

1. Start the docker compose environment.
   ```bash
   docker compose up -d
   ```
   Output should appear similar to the following:
   ```shell
   ✔ Network xk6-amqp_default       Created               ...    0.0s
   ✔ Container xk6-amqp-rabbitmq-1  Started               ...    0.2s
   ```
2. Use your [custom k6 binary](#build) to run a k6 test script connecting to your RabbitMQ server started in the previous step.
   ```bash
   ./k6 run examples/test.js
   ```
3. Use the RabbitMQ admin console by accessing [http://localhost:15672/](http://localhost:15672/), then login using `guest` for both the Username and Password.
   This will allow you to monitor activity within your messaging server.
