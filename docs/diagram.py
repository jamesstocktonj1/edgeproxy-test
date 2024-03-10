from diagrams import Diagram, Cluster
from diagrams.onprem.network import Envoy
from diagrams.programming.language import Go

with Diagram("edgeproxy", show=False):

    with Cluster("Container 1"):
        server1 = Go("api")
        with Cluster("edgeproxy"):
            proxy1 = Envoy()
        proxy1 >> server1
        
    with Cluster("Container 2"):
        server2 = Go("api")
        with Cluster("edgeproxy"):
            proxy2 = Envoy()
        proxy2 >> server2

    with Cluster("Container 3"):
        server3 = Go("api")
        with Cluster("edgeproxy"):
            proxy3 = Envoy()
        proxy3 >> server3

    proxy1 >> proxy2
    proxy2 >> proxy3
    # proxy3 >> proxy1

    proxy3 >> proxy2
    proxy2 >> proxy1
    # proxy1 >> proxy3

    proxy = Envoy("main proxy")
    proxy >> proxy1
    proxy >> proxy2
    proxy >> proxy3

