// Package amqp contains AMQP API for a remote server.
package amqp

import (
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	amqpDriver "github.com/rabbitmq/amqp091-go"
	"github.com/vmihailenco/msgpack/v5"
	"go.k6.io/k6/js/modules"
	"sync"
)

const version = "v0.4.0"

// RootModule is the root module for the AMQP API.
type RootModule struct {
	Queue    *Queue
	Exchange *Exchange
}

// AMQP type holds connection to a remote AMQP server.
type AMQP struct {
	vu modules.VU

	Version    string
	Connection *amqpDriver.Connection
	Queue      *Queue
	Exchange   *Exchange
}

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &AMQP{}
)

// New returns a pointer to a new RootModule instance
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance returns a new instance of the module given a virtual user.
func (r *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &AMQP{
		vu:       vu,
		Version:  version,
		Queue:    &Queue{},
		Exchange: &Exchange{},
	}
}

// Exports implements modules.Instance function to expose objects to JavaScript.
func (amqp *AMQP) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"Amqp":     amqp,
			"Exchange": amqp.Exchange,
			"Queue":    amqp.Queue,
			"version":  amqp.Version,
			"start":    amqp.Start,
			"publish":  amqp.Publish,
			"listen":   amqp.Listen,
		},
	}
}

// Options defines configuration options for an AMQP session.
type Options struct {
	ConnectionURL string
}

// PublishOptions defines a message payload with delivery options.
type PublishOptions struct {
	QueueName   string
	Body        string
	Headers     amqpDriver.Table
	Exchange    string
	ContentType string
	Mandatory   bool
	Immediate   bool
	Persistent  bool
}

// ConsumeOptions defines options for use when consuming a message.
type ConsumeOptions struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqpDriver.Table
}

// ListenerType is the message handler implemented within JavaScript.
type ListenerType func(string) error

// ListenOptions defines options for subscribing to message(s) within a queue.
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

const messagepack = "application/x-msgpack"

// Start establishes a session with an AMQP server given the provided options.
func (amqp *AMQP) Start(options Options) error {
	conn, err := amqpDriver.Dial(options.ConnectionURL)
	amqp.Connection = conn
	amqp.Queue.Connection = conn
	amqp.Exchange.Connection = conn
	return err
}

// Publish delivers the payload using options provided.
func (amqp *AMQP) Publish(options PublishOptions) (*goja.Promise, error) {
	ch, err := amqp.Connection.Channel()
	if err != nil {
		return nil, err
	}
	promise, resolve, reject := amqp.vu.Runtime().NewPromise()
	callback := amqp.vu.RegisterCallback()

	go func() {
		defer func() {
			_ = ch.Close()
		}()

		message := amqpDriver.Publishing{
			Headers:     options.Headers,
			ContentType: options.ContentType,
		}

		if options.ContentType == messagepack {
			var jsonParsedBody interface{}

			if err = json.Unmarshal([]byte(options.Body), &jsonParsedBody); err != nil {
				callback(func() error {
					reject(err)
					return nil
				})
			}

			message.Body, err = msgpack.Marshal(jsonParsedBody)
			if err != nil {
				callback(func() error {
					reject(err)
					return nil
				})
			}
		} else {
			message.Body = []byte(options.Body)
		}

		if options.Persistent {
			message.DeliveryMode = amqpDriver.Persistent
		}

		err = ch.PublishWithContext(
			amqp.vu.Context(),
			options.Exchange,
			options.QueueName,
			options.Mandatory,
			options.Immediate,
			message,
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
	return promise, nil
}

// Listen binds to an AMQP queue in order to receive message(s) as they are received.

func (amqp *AMQP) Listen(options ListenOptions) (*goja.Object, error) {
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

	callback := amqp.vu.RegisterCallback()
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
				callback = amqp.vu.RegisterCallback() // get a new reserve
				return options.Listener(string(d.Body))
			})
		}
	}()

	return o, nil
}

func init() {
	rootModule := &RootModule{Queue: &Queue{}, Exchange: &Exchange{}}

	modules.Register("k6/x/amqp", rootModule)
	modules.Register("k6/x/amqp/queue", rootModule.Queue)
	modules.Register("k6/x/amqp/exchange", rootModule.Exchange)
}
