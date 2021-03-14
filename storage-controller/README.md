# Storage controller

A controller to manage nodes for key-value storage.

## Notes

A set of pods and a service is managed by the controller.

Each pod is aware of every other pod.

Any pod could take a request - get / put, through the service.

For simplicity, any put request is broadcast among all pods.
any read request is read from all nodes and majority wins.
if there is no majority it means no valid data.

