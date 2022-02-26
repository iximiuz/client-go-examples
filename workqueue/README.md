# workqueue example - controllers' fundamentals

Based on client-go <a href="https://github.com/kubernetes/client-go/tree/cc43a708a08eb9ff6a436f0cb00c5ee05121d2cd/examples/workqueue">`examples/workqueue`</a>.
The example shows how to implement a primitive controller watching ADD/UPDATE/DELETE events for a particul kind of object.
The events are queued to allow safe parallel processing.
