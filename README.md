Personal practice to develop a distributed system of bunch of robots moving around, simulated by a bunch of threads making state mutations.

Sept 5 2018 -- 
Started diving in without much of design. Originally separates out to several components. Robot, Task/TaskManager and World. Using World as the point of entry.
Then run into a problem -- if the system is trully distributed, how much context/state does each robot maintains on its own.
Trying to solve the problem of state maintainance, I started passng pointers around and have wait locks on central world states. The system is somehow much slower. commit 574ab.
Each Robot execute in an GoRoutine concurrently.

Sept 21 2018 --
refactor the code base to have proper modules. Remove concurrency controls.

June 11 2019 --

Started implementing interfaces to abstract away all the confusing points of concurrent and distributed points. Trying to figure out a set of interfaces that work for both distributed message passing mode or contralized iterative  mode.


Sept 24 - Ocr 1 2019

Improve interface design. Trying to finish from previous work. Did some research. Other similar system are designed in stacks. Membership maintenance, state synchronization etc are in a different stack layer from main application. 

Started implementing a driving component called simuation that drives and maintains system starts and end. Let Robot be the one doing the bulk load, localize a copy of world state. Using driver to maintain the state synchronization in a centralized fashion.


Oct 14 2019

After first set of tests pass, reviewed the code. The interganbled pieces are due to trying to design interface and implementing basic objects fulfilling the interface -- making tests difficult. To ensure the interface is correct, test implementations should be implemented in test packages. This allows functions in test packages to introspect object states (Simulate), and some Basic Implementations could stay with the Interface package, like World, to avoid repetitive code (technically, for more robust testing, each test case should implement its own depended sub objects/interface) .


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
