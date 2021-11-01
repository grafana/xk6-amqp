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
	modules.InstanceCore

	Version    string
	Connection *amqpDriver.Connection
	Queue      *Queue
	Exchange   *Exchange
}

type AmqpRoot struct {
	Queue    *Queue
	Exchange *Exchange
}

func (amqp *AmqpRoot) NewModuleInstance(core modules.InstanceCore) modules.Instance {
	return &Amqp{
		InstanceCore: core,
		Version:      version,
		Queue:        amqp.Queue,
		Exchange:     amqp.Exchange,
	}
}

// GetExports returns the exports of the metrics module
func (amqp *Amqp) GetExports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"start":   amqp.Start,
			"listen":  amqp.Listen,
			"publish": amqp.Publish,
			"version": amqp.Version,
		},
	}
}

var _ modules.IsModuleV2 = &AmqpRoot{}

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
	p, resolve, reject := amqp.MakeHandledPromise()
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
			fmt.Println("reject err", err)
			reject(err)
		} else {
			fmt.Println("resolve")
			resolve(nil)
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

	_, resolve, _ := amqp.MakeHandledPromise() // we only care that the loop won't let the iteration end here
	rt := amqp.GetRuntime()
	stop := make(chan struct{})
	o := rt.NewObject()
	o.Set("stop", func() {
		stop <- struct{}{}
		<-stop
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		<-stop
		ch.Close()
		wg.Wait()
		resolve(nil) // we release the loop to end the iteration
		close(stop)
	}()

	go func() {
		defer wg.Done()
		for d := range msgs {
			amqp.AddToEventLoop(func() { options.Listener(string(d.Body)) })
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
