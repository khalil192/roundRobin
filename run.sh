#in case of any script termination flow, kill all the child process
trap cleanup INT TERM EXIT

#killing all the child processes of this script
cleanup() {
    echo "stopping backend servers"
    pkill -P $$
}

#------------Set up multiple BE servers------------

echo "Setting up multiple BE servers"

go build -o ./BackendServer/BackendServer ./BackendServer

./BackendServer/BackendServer --port=8080 &
./BackendServer/BackendServer --port=8081 &
./BackendServer/BackendServer --port=8082 &
./BackendServer/BackendServer --port=8083 &

#-----------------------------------------------------

#---------------Set up the Load Balancer------------

echo "Setting up Load Balancer"

go build -o ./LoadBalancer/LoadBalancer ./LoadBalancer

#  select one of the algo implementation as a command line argument : "lock", "atomic" , "queue",  "separate_slice":

./LoadBalancer/LoadBalancer --port=9000 --algoToUse=separate_slice &

#Pauce enough time so that our LB can run  first health.
sleep 10

#---------------Running Stress Test---------------------------------

echo "Running stress-tests "

go build -o ./stresstest/stresstest ./stresstest

#maxCons => no of concurrent client of connections
#numReqs => no of request each concurrent client will make
./stresstest/stresstest --numReqs=5000 --maxCons=50


echo "View the results in server_stats.csv"

