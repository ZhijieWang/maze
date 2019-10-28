Personal practice to develop a distributed system of bunch of robots moving around, simulated by a bunch of threads making state mutations.

Sept 5 2018 -- 
Started diving in without much of design. Originally separates out to several components. Robot, Task/TaskManager and World. Using World as the point of entry.
Then run into a problem -- if the system is truly distributed, how much context/state does each robot maintains on its own.
Trying to solve the problem of state maintenance, I started passing pointers around and have wait locks on central world states. The system is somehow much slower. commit 574ab.
Each Robot execute in an GoRoutine concurrently.

Sept 21 2018 --
refactor the code base to have proper modules. Remove concurrency controls.

June 11 2019 --

Started implementing interfaces to abstract away all the confusing points of concurrent and distributed points. Trying to figure out a set of interfaces that work for both distributed message passing mode or contralized iterative  mode.


Sept 24 - Ocr 1 2019

Improve interface design. Trying to finish from previous work. Did some research. Other similar system are designed in stacks. Membership maintenance, state synchronization etc are in a different stack layer from main application. 

Started implementing a driving component called simulation that drives and maintains system starts and end. Let Robot be the one doing the bulk load, localize a copy of world state. Using driver to maintain the state synchronization in a centralized fashion.


Oct 14 2019

After first set of tests pass, reviewed the code. The intermingled pieces are due to trying to design interface and implementing basic objects fulfilling the interface -- making tests difficult. To ensure the interface is correct, test implementations should be implemented in test packages. This allows functions in test packages to introspect object states (Simulate), and some Basic Implementations could stay with the Interface package, like World, to avoid repetitive code (technically, for more robust testing, each test case should implement its own depended sub objects/interface) .

Oct 17 2019
Some basic simulation was run. Realized some limitation in existing Application stack interface. Start implement usecase further to identify limitation of the interface.


Oct 24 2019 
Implemented basic action graph generation and execution. With sufficient tests to pass as first implementation.
Fixed some minot code import cycles to make sure the build pass.
Improve the construct of robots. A robot should have a Init function which is its initialization, a Plan function which can read from senses and context to devise a plan; a Execute function which executes the action plan that derived from Plan function. The overall Plan then Execute can be wrapped into a Run function, as an entry point.

In the plan function, a robot can read its world or many other things. Technically the sensing can be done in another stack layer, constantly pulling or updating the sensor reading. Rather than the main thread triggers the data fetching and wait for cycle to complete.

A typical client and server model will require a sever, which server as many functionality:
1. Name discovery
2. Synchronization and message bus
3. Orchestration


Drawbacks -- not fast enough and not scalable to millions of complex requests.

Tiered approach like MVC (model view controller, and throw the persistence to database) allows data passing in sequential orders (best).
Drawbacks -- still not so scalable, when synchronization is needed.

Pure message passing based solution (distributed micro-services /Actor systems) with proper message design.
Drawbacks -- super pain in the neck to debug, stuck in message loops.

I guess I shouldn't be jumping in so early to achieve 3. But let's start with 1.

RPC based simulation. Bots announce their location to central. Bots observe their location from central. Central persist location, performs synchronization, performs message delivery.

First step is channels, let next be RPC, and finally just messages flying around. Let tracing be fun.



# maze
A simulator for robotic task planning over a network graph with task dependencies
