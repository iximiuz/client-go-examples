# Retry on Conflict

`retry.RetryOnConflict(backoff wait.Backoff, fn func() error) error` is a
useful function when creating custom controllers or operators.
`RetryOnConflict` can help reduce the number of transient errors within
your program. For more in depth information [see the official docs][].

## Usage

Typical usage for `RetryOnConflict` occurs when there's a possibility
of many clients interacting with a resource at the same time or multiple controllers
are interacting with the same resource, or a mixture of both.

A typical implementation would be:

```golang
err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
    // Always fetch the new version of the resource
    pod, err := c.Pods("mynamespace").Get(name, metav1.GetOptions{})
    if err != nil {
        return err
    }

    // *************
    // Make some form of long running change, hitting external apis,
    // spinning up and validating external resources, etc
    // *************

    // Try to update
    _, err = c.Pods("mynamespace").UpdateStatus(pod)
    // You have to return err itself here (not wrapped inside another error)
    // so that RetryOnConflict can identify it correctly.
    return err
})
if err != nil {
    return err
}
```

## This Example

This example is intended to fail to show `retry.RetryOnConflict` attempting to
retry the operation. The expected outcome is:

```log
$ go run main.go 
Successfully updated ConfigMap
Operation cannot be fulfilled on configmaps "foobar": the object has been modified; please apply your changes to the latest version and try again
Successfully updated ConfigMap
Operation cannot be fulfilled on configmaps "foobar": the object has been modified; please apply your changes to the latest version and try again
Successfully updated ConfigMap
Operation cannot be fulfilled on configmaps "foobar": the object has been modified; please apply your changes to the latest version and try again
Successfully updated ConfigMap
Operation cannot be fulfilled on configmaps "foobar": the object has been modified; please apply your changes to the latest version and try again
Successfully updated ConfigMap
Operation cannot be fulfilled on configmaps "foobar": the object has been modified; please apply your changes to the latest version and try again
```

[see the official docs]: https://pkg.go.dev/k8s.io/client-go/util/retry#RetryOnConflict
