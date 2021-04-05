# gebus
GeBus is simple but powerful event bus implementation for golang.

It supports:
1. Synchronous event bus
2. Goroutine pool event bus (acutally allocating new ones rather than pooling, so that closeing is not needed)
3. Distributed AMQP based event bus(using external library with AMQP client)
4. Event bus using reflection to make typed event handlers, which converts function into functional event handler.

For examples take a look at *_test files.