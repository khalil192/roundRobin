
I loved solving this problem:)

I understood that solving for scale was the real challenge here, ensuring that the system could handle a high concurrent load.

I love to build testable systems, and since round-robin can be implemented in several ways, I wanted to try different approaches and compare their throughput and fairness. 

Hence, I decided to build a testable LoadBalancer interface first. With a proper interface in place, I can plug in new implementations and even new resource allocation algorithms in a "plug-and-play" fashion.


The project is a mono repo with three directories and a shell file in the root:

- BackendServer: A dummy server, of which we will spawn multiple instances.

- LoadBalancer: The core of the project. A reverse proxy server that forwards requests from the client to one of the BackendServer instances.
  
  - The HTTP reverse proxy server is decoupled from the algorithm's logic.
  - The LoadBalancer algorithm is injected via dependency inversion into the HTTP handler. This enables the selection of any implementation when starting the program.
  - Unit tests are present to ensure the correctness of each implementation's data structure usage.


- StressTest: Client code that tests our LoadBalancer's performance under a high load.


- run.sh: A shell file for the final integration stress test. It outputs the server load distribution to server_stats.csv.

-------------------------------------

In the LoadBalancer/algorithms folder I've Implemented different versions of roundRobin


