package jsonrpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/shiftdevices/godbb/util/errp"
	"github.com/shiftdevices/godbb/util/jsonp"
	"github.com/sirupsen/logrus"
)

const responseTimeout = 100 * time.Second

type callbacks struct {
	success func([]byte) error
	cleanup func()
}

// SocketError indicates an error when reading from or writing to a network socket.
type SocketError error

// ResponseError indicates an error when parsing the response from the server.
type ResponseError error

// RPCClient is a generic json rpc client, which is able to invoke remote methods and subscribe to
// remote notifications.
type RPCClient struct {
	conn                  io.ReadWriteCloser
	msgID                 int
	msgIDLock             sync.Mutex
	responseCallbacks     map[int]callbacks
	responseCallbacksLock sync.RWMutex
	close                 bool

	notificationsCallbacks     map[string]func([]byte)
	notificationsCallbacksLock sync.RWMutex
	failureCallback            func(error)
	log                        *logrus.Entry
}

// NewRPCClient creates a new RPCClient. conn is used for transport (e.g. a tcp/tls connection).
func NewRPCClient(conn io.ReadWriteCloser, failureCallback func(error), log *logrus.Entry) *RPCClient {
	client := &RPCClient{
		conn:                   conn,
		msgID:                  0,
		responseCallbacks:      map[int]callbacks{},
		notificationsCallbacks: map[string]func([]byte){},
		failureCallback:        failureCallback,
		log:                    log,
	}
	go client.read(client.handleResponse, client.failureCallback)
	return client
}

func (client *RPCClient) read(success func([]byte), failure func(error)) {
	defer func() {
		_ = client.conn.Close()
		if r := recover(); r != nil {
			failure(r.(error))
		}
	}()
	reader := bufio.NewReader(client.conn)
	for !client.close {
		line, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			panic(errp.Wrap(err, "Failed to read from socket"))
		}
		success(line)
	}
}

func (client *RPCClient) handleResponse(responseBytes []byte) {
	// fmt.Println("got response ", string(responseBytes))

	// Catch all response.
	// A notification contains:
	// - jsonrpc
	// - method
	// - params
	// A method call response contains:
	// - jsonrpc
	// - id
	// - result
	// - (error)
	response := &struct {
		JSONRPC string           `json:"jsonrpc"`
		ID      *int             `json:"id"`
		Error   *json.RawMessage `json:"error"`
		Result  json.RawMessage  `json:"result"`
		Method  *string          `json:"method"`
		Params  json.RawMessage  `json:"params"`
	}{}
	if err := json.Unmarshal(responseBytes, response); err != nil {
		client.failureCallback(ResponseError(errp.Wrap(err, "Failed to unmarshal response")))
		return
	}
	if response.JSONRPC != "2.0" {
		client.failureCallback(ResponseError(errp.Newf("Unexpected json rpc version: %s", response.JSONRPC)))
		return
	}
	if response.Error != nil {
		client.failureCallback(ResponseError(errp.New(string(*response.Error))))
		return
	}

	// Handle method response.
	if response.ID != nil {
		func() {
			client.responseCallbacksLock.RLock()
			responseCallbacks, ok := client.responseCallbacks[*response.ID]
			client.responseCallbacksLock.RUnlock()
			if ok {
				if len(response.Result) == 0 {
					client.failureCallback(ResponseError(errp.New("unexpected reply")))
				}
				if err := responseCallbacks.success([]byte(response.Result)); err != nil {
					client.failureCallback(ResponseError(errp.WithMessage(err, "success callback error")))
				}
				responseCallbacks.cleanup()
				client.responseCallbacksLock.Lock()
				delete(client.responseCallbacks, *response.ID)
				client.responseCallbacksLock.Unlock()
			}
		}()
	}

	// Handle notification.
	if response.Method != nil {
		if len(response.Params) == 0 {
			client.failureCallback(ResponseError(errp.New("unexpected reply")))
			return
		}
		func() {
			client.notificationsCallbacksLock.RLock()
			responseCallback, ok := client.notificationsCallbacks[*response.Method]
			client.notificationsCallbacksLock.RUnlock()
			if ok {
				responseCallback([]byte(response.Params))
			}
		}()
	}
}

// Method sends invokes the remote method with the provided parameters. The success callback is
// called with the response. cleanup is called afterwards, regardless of whether an error occurred
// anywhere.
func (client *RPCClient) Method(
	success func([]byte) error,
	cleanup func(),
	method string,
	params ...interface{},
) error {
	client.msgIDLock.Lock()
	msgID := client.msgID
	client.msgID++
	client.msgIDLock.Unlock()
	if params == nil {
		params = []interface{}{}
	}
	jsonText := append(jsonp.MustMarshal(map[string]interface{}{
		"id":     msgID,
		"method": method,
		"params": params,
	}), byte('\n'))

	client.responseCallbacksLock.Lock()
	client.responseCallbacks[msgID] = callbacks{
		success: success,
		cleanup: cleanup,
	}
	client.responseCallbacksLock.Unlock()

	return client.send(jsonText)
}

// MethodSync is the same as method, but blocks until the response is available. The result is
// json-deserialized into response.
func (client *RPCClient) MethodSync(response interface{}, method string, params ...interface{}) error {
	responseChan := make(chan []byte)
	if err := client.Method(
		func(responseBytes []byte) error {
			responseChan <- responseBytes
			return nil
		},
		func() {},
		method, params...); err != nil {
		return err
	}
	select {
	case responseBytes := <-responseChan:
		if response != nil {
			if err := json.Unmarshal(responseBytes, response); err != nil {
				return ResponseError(errp.Wrap(err, "Failed to unmarshal response"))
			}
		}
	case <-time.After(responseTimeout):
		return SocketError(errp.New("response timeout"))
	}
	return nil
}

// SubscribeNotifications installs a callback for a method which is called with notifications from
// the server.
func (client *RPCClient) SubscribeNotifications(method string, callback func([]byte)) {
	client.notificationsCallbacksLock.Lock()
	defer client.notificationsCallbacksLock.Unlock()
	if _, ok := client.notificationsCallbacks[method]; ok {
		panic(fmt.Sprintf("already subscribed to notifications of %s", method))
	}
	client.notificationsCallbacks[method] = callback
}

func (client *RPCClient) send(msg []byte) error {
	_, err := client.conn.Write(msg)
	if err != nil {
		return SocketError(errp.WithStack(err))
	}
	return nil
}

// Close shuts down the connection.
func (client *RPCClient) Close() {
	client.close = true
}
