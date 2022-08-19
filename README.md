# goose
A Golang implementation of Server Sent Events

### Overview

Server Sent Events (SSE) is a simple one to provide near real-time data to a HTTP client. It is unidirectional - the serer pushes new information to the client. This library provides a simple wrapper to setting up an SSE connection and pushing event data to clients simple and easy. 

### Setup

The `EventStream` struct conatains all the logic for handling pushing data. Create and initialize your `EventStream` inside a standard HTTP handler func:

```
func myHttpHandler(w http.ResponseWriter, r *http.Request) {

    eventStream := goose.NewEventStream(w)
}
```

Nothing is sent to the client until the stream's `Begin` function is called. First, you must create a channel to pass event messages to. The `Begin` function blocks until the channel is closed, so be sure to call it at the end of your handler function.

```
func myHttpHandler(w http.ResponseWriter, r *http.Request) {

    eventStream := goose.NewEventStream(w)

    defer eventStream.Close()

    eventChan := make(chan string)

    // asynchronous method to push data to the channel
    go func() {
        for i:=0;i<10;i++ {
            time.Sleep(5 * time.Second)
            eventChan <- fmt.Sprintf("%d", i+1)
        }

        close(eventChan)
    }()

    // Blocks until eventChan is closed
    if err := eventStream.Begin(eventChan); err != nil {
        // Handle Error here
    }
}
```

### Examples

Serialize JSON Object and push to event stream with the builtin `json` package

```
func myHttpHandler(w http.ResponseWriter, r *http.Request) {

    var objectChan chan interface{} 

    eventStream := goose.NewEventStream(w)

    defer eventStream.Close()

    eventChan := make(chan string)

    // asynchronous method to push data to the channel
    go func() {
        for obj := range objectChan {
            data, _ := json.Marshal(&obj)
            eventChan <- string(data)
        }
        close(eventChan)
    }()

    // Blocks until eventChan is closed
    if err := eventStream.Begin(eventChan); err != nil {
        // Handle Error here
    }
}
```