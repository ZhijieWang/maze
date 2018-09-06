A typical client and server model will require a sever, which server as many functionalities:
1. Name discovery
2. Syncrhonization and message bus
3. Orchestration


Drawbacks -- not fast enough and not scalable to millons of comples requests.

Tiered approach like MVC (model view controller, and throw the persistency to database) allows data passing in sequential orders (best).
Drawbacks -- still not so scallable, when syncrhonization is needed.

Pure message passing based solution (distributed micro-services /Actor systems) with proper message design.
Drawbacks -- super pain in the neck to debug, stuck in message loops.

I guess I shouldn't be jumping in so early to achieve 3. But let's start with 1.

RPC based simulation. Bots announce their location to central. Bots observe their location from central. Central persist location, performs syncrhonization, performs message delivery.

First step is channels, let next be RPC, and finally just messages flying around. Let tracing be fun.



# maze
A simulator for robotic task planning over a network graph with task dependencies
