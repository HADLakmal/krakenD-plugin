package main

import (
	"errors"
	"fmt"
	"io"
)

func main() {}

func init() {
	fmt.Println(string(ModifierRegisterer), "loaded!!!")
}

var extraHeader string = "extra-header"

var ModifierRegisterer = registerer("response-plugin")

type registerer string

// RegisterModifiers is the function the plugin loader will call to register the
// modifier(s) contained in the plugin using the function passed as argument.
// f will register the factoryFunc under the name and mark it as a request
// and/or response modifier.
func (r registerer) RegisterModifiers(f func(
	name string,
	factoryFunc func(map[string]interface{}) func(interface{}) (interface{}, error),
	appliesToRequest bool,
	appliesToResponse bool,
)) {
	f(string(r), r.responseDump, false, true)
	fmt.Println(string(r), " registered!")
}

// ResponseWrapper is an interface for passing proxy response between the krakend pipe
// and the loaded plugins
type ResponseWrapper interface {
	Data() map[string]interface{}
	Io() io.Reader
	IsComplete() bool
	StatusCode() int
	Headers() map[string][]string
}

type responseWrapper struct {
	data       map[string]interface{}
	isComplete bool
	metadata   metadataWrapper
	io         io.Reader
}

type metadataWrapper struct {
	headers    map[string][]string
	statusCode int
}

func (m metadataWrapper) Headers() map[string][]string { return m.headers }
func (m metadataWrapper) StatusCode() int              { return m.statusCode }

func (r responseWrapper) Data() map[string]interface{} { return r.data }
func (r responseWrapper) IsComplete() bool             { return r.isComplete }
func (r responseWrapper) Io() io.Reader                { return r.io }
func (r responseWrapper) Headers() map[string][]string { return r.metadata.headers }
func (r responseWrapper) StatusCode() int              { return r.metadata.statusCode }

var unkownTypeErr = errors.New("unknow request type")

func (r registerer) responseDump(
	cfg map[string]interface{},
) func(interface{}) (interface{}, error) {
	config := cfg[string(r)]
	var extraHeaderConfig string
	if config != nil {
		if cf, ok := config.(map[string]interface{}); ok {
			extraHeaderConfig = fmt.Sprintf(`%v`, cf[extraHeader])
		}
	}
	fmt.Println("response dumper injected!!!")
	return func(input interface{}) (interface{}, error) {
		fmt.Printf("input %v", input)
		fmt.Println("")
		resp, ok := input.(ResponseWrapper)
		if !ok {
			return nil, unkownTypeErr
		}

		resTemp := responseWrapper{}
		resTemp.data = resp.Data()
		resTemp.isComplete = resp.IsComplete()
		resTemp.io = resp.Io()
		resTemp.metadata.headers = make(map[string][]string)
		for k, v := range resp.Headers() {
			resTemp.metadata.headers[k] = v
		}
		if extraHeaderConfig != "" {
			resTemp.metadata.headers[extraHeader] = []string{extraHeaderConfig}
		}

		return resTemp, nil
	}
}
