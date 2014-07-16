go-hsm
======

```go-hsm``` is a Golang library that implements the Hierarchical State Machine(HSM).

Hierarchical State Machine(HSM) is a method to express state machine. It's introduced in the book [```Practical Statecharts in C/C++: Quantum Programming for Embedded System```][qp-book-homepage] by Dr. Miro M. Samek. Comparing to the traditional methods to implement state machine(e.g. nested if-else/switch, state table, state design pattern in OOP), HSM provides some major advantage as the followings:

1. It supports nested states and hehavior inheritance
2. It provides entry and exit action for state
3. It uses class hierarchy to express state hierarchy. Easy to write and understand.

The book mentioned above gives an elaboration on state machine and HSM. For more details, please refer to it.

Whenever there arises the requirement to write codes which has high complexity but need to be definitely correct, state machine is a promising way to get there. HSM is a great pattern to express state machine. And HSM is the most powerful pattern I know so far. At the moment I ran into some Golang project like [```Raft```][raftconsensus-homepage], HSM occurs to me. So I port HSM to Golang and create this project.

## Port HSM to Golang

Along with the theory of HSM, there are two implementations of it in the book, one in C and the other in C++. Some language specified features are used in these implementations. However, Golang has some significant differencies in the programming language features. So porting HSM to Golang would not be doing a simple copycat of the original implementations.

In both implementations, methods(member functions of class) are used to represent states. However no member function pointer in Golang. And In the C/C++ implementation the class inheritance hierarchy are used to represent the structure of states for the whole state machine. The OOP inheritance in Golang is quite different by using interface. In Golang, there is no parent pointer for parent class object in an object's memory layout.

My first thought is to use function variable/lambda function  to represent state. But it's not easy to maintain the state hierarchy between functions. After some experiments on trying to get HSM up and running in Golang, it turns out to be the way it's as now: class represents State, and state hierarchy is done munually(and verbosely) in initialization of the whole state machine. 

IMO, implementing HSM in this way isn't perfect. There are some pitfalls:

1. To define a new state, you always need to write super.AddChild() in the New() constructor, which I consider as a abstraction leak. Anyway it's boring.
2. Defining states could be lot of chunks of codes to write. All of states have the similar Entry/Exit/Handle code structure. To provide some kinds of code template for code reusage, the lack of meta programming ability in Golang leads me to a DSL and code generator, which is too heavy and too complicated.
3. The type casts when handling event. Casting is inevitable in such a generic framework. But I don't think the library user should write these codes by hand again and again. It's unwelcomed especially while assert is taken out from language libraries in Golang, although casting needs it to ensure type correctness. (IMO, assert is a dark corner in Golang. The language just take it out. But when someone revives it as a library(see [testify.assert][testify-github]), people find it quite usefull and wants it.


## The Missings

This project contains only the HSM, briefly a method to construct state machine and dispatch events. It's not a full Quantum Framework. So there are a lot missing:

1. other state machine methods, e.g. FSM or optimized FSM
2. concurrency mechanism, e.g. go-routine
3. composite states
4. local and external transitions, internal transitions
5. event queuing
6. orthogonal regions
7. validation of state machine structure, with checks for:
    * machine having single top state
    * unreachable states
    * multiple occurrences of same state object instance
    * multiple states with same name
    * transitions that start from or point to nonexistent states

## Usage

Refer to the project [```go-hsm-examples```][go-hsm-examples-github].


[go-hsm-examples-github]: https://github.com/hhkbp2/go-hsm-examples
[qp-book-homepage]: http://www.state-machine.com/psicc/
[raftconsensus-homepage]: http://raftconsensus.github.io/
[testify-github]: https://github.com/stretchr/testify
