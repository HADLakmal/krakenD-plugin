# KrakenD Response Plugin
KrakenD is a high performance API Gateway. that can communicate between the client and all the source servers, adding a new layer that removes all the complexity to the clients, providing them only the information that the API response needs.
### Install KrakenD
Clone git-hub [repo](https://github.com/krakendio/krakend-ce) then follow the steps in the readme to build the krakend binary file. That is easy to match underline requirement of the custom plugins. Also you can install the binary file from [here](https://www.krakend.io/docs/overview/installing/).
### Build Plugin
Response plugin build in this project,

```go build -buildmode=plugin -o res-plugin.so .```

