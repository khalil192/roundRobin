
I loved solving this problem:)

-----------------------------------
Problem Statement: Implement a reverse-proxy server which uses round-robin algorithm to route the requests to one of the BE server.

Core problem to solve: Handling heavy concurrent load and ensuring fairness of allocation.

Tech stack used: `Golang`, `shell-script`.

--------------------------

Since round-robin can be implemented in several ways, I wanted to try different approaches and compare their throughput and fairness. 

The project is a mono repo with three directories and a shell file in the root:

- BackendServer: A dummy server, we will spawn multiple instances.

- LoadBalancer: The core of the project. A reverse proxy server that forwards requests from the client to one of the BackendServer instances.
  - The HTTP reverse proxy server is decoupled from the algorithm's logic.
  - The LoadBalancer algorithm is injected via dependency inversion into the HTTP handler. This enables the selection of any implementation when starting the program.
  - Unit tests are present to ensure the correctness of each implementation's data structure usage.

- StressTest: Client code that tests our Load Balancers performance under a high load.

- run.sh: A shell file for the final integration stress test. It outputs the server load distribution to server_stats.csv.

-------------------------------------

``Implementaion Details``

In the LoadBalancer/algorithms folder I've Implemented different versions of roundRobin

Our Main function is : GetNextHealthyServer()
    -> all the implementation's should implement this function which our reverse-proxy server uses before assigning server to req

Requirements: All implementations should be thread safe.

- At first, I started with simple global lock version, in which we lock the load-balancer data and find a healthy server and unlock.
- `Global Lock version` is taking too much time under heavy load. Since its efficient to lock the whole loadbalancer's data(serverList)
- Improved upon it by using `atomic operations` instead of mutex lock for fetching the next Healthy Server.
- Atomic is performing good and in case of huge concurrent requests, atomic implementation is 2X speed than global_lock.
- I identified Atomic implementation lagging when the healthy servers are skewed in the list, since we traverse the whole list till we find the next healthy server.
- To make it better I want to have an implementation where we store all the healthy servers in one place, then we don't need to traverse anything.
- `Separate slice` Implementation is the solution for this, maintaining a separate array/slice for healthy and unhealthy servers.
- `Separate Slice` Implementation uses `Read` lock during GetNextHealthyServer() ensuring Massive parallel read throughput and it only uses write lock when updating the data of server.

`Bonus` : I added a periodic health check worker which checks the health of servers in background (for every x secs) and updates the server health.

The interface is built in a plug-n-play fashion which makes adding newer algorithms/implementations very easy:
Sample MR for reference: https://github.com/khalil192/roundRobin/pull/3.

-------------------------------

Testing and Results: 

- Correctness of the algorithm is tested by unittests in LoadBalancer folder.
- Additionally, I want to test the robustness under heavy load -> stress-testing.
- I did stress testing using postman, but i want to have a configurable script in one place so that i demo it effectively.
- `Stresstest` folder contains the go script which makes concurrent requests to reverse proxy server.
- the final results are posted to server_stats.csv

--------------------

Running the code: 

Install all dependencies:
- `go mod tidy` 

Read the script file and configure as per requirements:
- sh run.sh


---------------- 