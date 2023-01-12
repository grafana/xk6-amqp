//go:build integration
// +build integration

package amqp

import (
	"context"
	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/eventloop"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/testutils"
	"go.k6.io/k6/metrics"
	"io/ioutil"
	"testing"
)

// setupTestEnv should be called from each test to build the execution environment for the test
func setupTestEnv(t *testing.T) *goja.Runtime {
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper{})

	testLog := logrus.New()
	testLog.AddHook(&testutils.SimpleLogrusHook{
		HookedLevels: []logrus.Level{logrus.WarnLevel},
	})
	testLog.SetOutput(ioutil.Discard)

	state := &lib.State{
		Options: lib.Options{
			SystemTags: metrics.NewSystemTagSet(metrics.TagVU),
		},
		Logger: testLog,
		Tags:   lib.NewVUStateTags(metrics.NewRegistry().RootTagSet()),
	}

	vu := &modulestest.VU{
		RuntimeField: rt,
		InitEnvField: &common.InitEnvironment{},
		CtxField:     context.Background(),
		StateField:   state,
	}

	root := &RootModule{}
	m := root.NewModuleInstance(vu)
	require.NoError(t, rt.Set("Amqp", m.Exports().Named["Amqp"]))
	require.NoError(t, rt.Set("Queue", m.Exports().Named["Queue"]))
	require.NoError(t, rt.Set("Exchange", m.Exports().Named["Exchange"]))

	ev := eventloop.New(vu)
	vu.RegisterCallbackField = ev.RegisterCallback

	return rt
}

func TestAMQP(t *testing.T) {
	t.Parallel()

	rt := setupTestEnv(t)

	_, err := rt.RunString(`
Amqp.start({
	connection_url: "amqp://guest:guest@localhost:5672/"
})
  
const queueName = 'k6-intg-test'
  
Queue.declare({
	name: queueName,
	durable: true
})

Amqp.publish({
	queue_name: queueName,
	body: "Message from integration test",
	content_type: "text/plain"
})

const listener = function(data) { console.log('received data: ' + data) }
Amqp.listen({
    queue_name: queueName,
    listener: listener,
    auto_ack: true
})

`)
	require.NoError(t, err)
}
