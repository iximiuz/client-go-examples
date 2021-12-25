# Create config from kubeconfig file using clientcmd.BuildConfigFromFlags()

The example code looks simple, but the internals aren't so simple, actually.
`clientcmd.BuildConfigFromFlags()` under the hood uses the `DeferredLoadingClientConfig` struct
(via `NewNonInteractiveDeferredLoadingClientConfig`) that implements `ClientConfigLoader` interface.
Such kind of loader is needed to defer the actual config creation util all the possible tweaks (via
extra flags and/or env vars) are done. However, since the `ExplicitPath` is used, neither the deferred
loading nor actual merging happens.
