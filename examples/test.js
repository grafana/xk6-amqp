import { sleep } from "k6";
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


	var o;
	var i = 0;
	const listener = function(data) { 
		console.log(`received data by VU ${__VU}, ITER ${__ITER}: ${data} `)
		i++
		if (i == 2) {
			console.log("stopping")
			o.stop() // stop on the  second message
		}
	}

	o = Amqp.listen({
		queue_name: queueName,
		listener: listener,
		auto_ack: true,
	})


	var p = Amqp.publish({
		queue_name: queueName,
		body: "Ping from k6 vu: " + __VU + ", iter:"+ __ITER,
	})

	p.then(() => {
		console.log("then resolved")
		sleep(2)
				Amqp.publish({
			queue_name: queueName,
			body: "Second ping from k6 vu: " + __VU + ", iter:"+ __ITER,
		}).then(() => {console.log("second send resolved")})
	}, (e) => {
				console.log("reject" + e);
	})
}
