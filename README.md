DocShare
========

An application to share documents via blockchain. Hyperledger fabric is used as a framework here.


To start the network go to the network folder and execute:
- `$ ./network.sh up`

To compile the code run the following commands in the base directory:
- `$ go get -v`
- `$ go build`

To run the code execute the following command:
- `$ ./docshare`

To shut down the network go to the network folder and execute:
- `$ ./network.sh down`

The base code is mostly inspired from [here](https://github.com/chainHero/heroes-service/)
