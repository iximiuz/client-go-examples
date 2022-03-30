# Deserialize kubeconfig YAML into kubeconfig Go struct and create rest.Config out of it

The trick is helpful when the Kubeconfig content is obtained from a non-disk location (e.g., read from a Kubernetes secret).