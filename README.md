# Edgeproxy Example

This proof of concept looks at local meshing in microservices. This is done by deploying a light-weight Envoy proxy in the same container as the code. All the meshing is handled by this proxy which simplifies the server code as it doesn't need to know which hosts are available. By sending any requests straight to the proxy, Envoy then handles which services are available to process the requst.

![System Diagram](https://github.com/jamesstocktonj1/edgeproxy-test/blob/main/docs/edgeproxy.png)