package amqp

import (
	"fmt"
	"sync"

	"github.com/dop251/goja"
	amqpDriver "github.com/streadway/amqp"
	"go.k6.io/k6/js/modules"
)

const version = "v0.0.1"

type Amqp struct {
	vu modules.VU

	Version    string
	Connection *amqpDriver.Connection
	Queue      *Queue
	Exchange   *Exchange
}

type AmqpRoot struct {
	Queue    *Queue
	Exchange *Exchange
}

func (amqp *AmqpRoot) NewModuleInstance(vu modules.VU) modules.Instance {
	return &Amqp{
		vu:       vu,
		Version:  version,
		Queue:    amqp.Queue,
		Exchange: amqp.Exchange,
	}
}

// GetExports returns the exports of the metrics module
func (amqp *Amqp) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"start":   amqp.Start,
			"listen":  amqp.Listen,
			"publish": amqp.Publish,
			"version": amqp.Version,
		},
	}
}

var _ modules.Module = &AmqpRoot{}

type AmqpOptions struct {
	ConnectionUrl string
}

type PublishOptions struct {
	QueueName string
	Body      string
	Exchange  string
	Mandatory bool
	Immediate bool
}

type ConsumeOptions struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqpDriver.Table
}

type ListenerType func(string) error

type ListenOptions struct {
	Listener  ListenerType
	QueueName string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqpDriver.Table
}

func (amqp *Amqp) Start(options AmqpOptions) error {
	conn, err := amqpDriver.Dial(options.ConnectionUrl)
	amqp.Connection = conn
	amqp.Queue.Connection = conn
	amqp.Exchange.Connection = conn
	return err
}

func (amqp *Amqp) Publish(options PublishOptions) (*goja.Promise, error) {
	ch, err := amqp.Connection.Channel()
	if err != nil {
		return nil, err
	}
	callback := amqp.vu.Reserve()
	p, resolve, reject := amqp.vu.Runtime().NewPromise()
	go func() {
		defer ch.Close()
		err := ch.Publish(
			options.Exchange,
			options.QueueName,
			options.Mandatory,
			options.Immediate,
			amqpDriver.Publishing{
				ContentType: "text/plain",
				Body:        []byte(options.Body),
			},
		)
		if err != nil {
			callback(func() error {
				fmt.Println("reject err", err)
				reject(err)
				return nil
			})
		} else {
			callback(func() error {
				fmt.Println("resolve")
				resolve(nil)
				return nil
			})
		}
	}()
	return p, nil
}

func (amqp *Amqp) Listen(options ListenOptions) (*goja.Object, error) {
	ch, err := amqp.Connection.Channel()
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		options.QueueName,
		options.Consumer,
		options.AutoAck,
		options.Exclusive,
		options.NoLocal,
		options.NoWait,
		options.Args,
	)
	if err != nil {
		return nil, err
	}

	callback := amqp.vu.Reserve()
	rt := amqp.vu.Runtime()
	stop := make(chan struct{})
	o := rt.NewObject()
	o.Set("stop", func() {
		select {
		case <-stop:
		case <-amqp.vu.Context().Done():
		}
		<-stop
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case stop <- struct{}{}:
		case <-amqp.vu.Context().Done():
		}
		err := ch.Close()
		if err != nil {
			amqp.vu.State().Logger.WithError(err).Warnf("closing amqp Listen connection")
		}

		wg.Wait()
		callback(func() error { return nil })
		close(stop)
	}()

	go func() {
		defer wg.Done()
		for d := range msgs {
			callback(func() error {
				callback = amqp.vu.Reserve() // get a new reserve
				return options.Listener(string(d.Body))
			})
		}
	}()

	return o, nil
}

func init() {
	queue := &Queue{}
	exchange := &Exchange{}
	modules.Register("k6/x/amqp", &AmqpRoot{Queue: queue, Exchange: exchange})
	modules.Register("k6/x/amqp/queue", queue)
	modules.Register("k6/x/amqp/exchange", exchange)
}
