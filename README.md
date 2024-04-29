# 3COLOR-ZKP
## Team members: Timothy Fong, Edward Wibowo

## Description:
Contains:
- A server/verifier that verifies a given graph from the client/prover. 
- A graph directory which contains example graphs to play with. For example, the following texts will represent a graph with 3 vertices, 2 of which are red, 1 is green:
0 red
1 red
2 green

0 1
1 2
2 0
- The text files will then be parsed into a graph, which is given to the prover to be sent to the verifier. 
- Networking between the prover and the verifier.
- A 3Color-ZKP implementation between the prover and the verifier: the prover commits to a graph coloring, sends to verifier. Verifier chooses two random adjacent vertices, checks color of the two vertices. Prover then shuffle the colors so the verifier can't learn anything from the previous rounds. This shuffling happens every iteration. 

## Example usage:
- make
- Type "./verifier -port 1234 -repetitions 30" using a port number of your choice and repetition numbers of your choice. This starts the server on the specified port number. 
- In another terminal, type "./prover -address localhost -port 1234 -graph ./graphs/improper-c3.txt" using an example graph of your choosing. In this case, about a third of the iterations would fail, because there is about a 1/3 chance that the verifier chooses an edge with vertices of the same color.
- The terminal on the server/verifier side should print the number of repetitions passed. 

Estimated time to complete: >16.0 hours.