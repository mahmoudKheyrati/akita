# AkitaSim Hands-On Exercises

This guide provides practical exercises to help you master AkitaSim concepts through building real components.

## Table of Contents

1. [Exercise 1: Event Basics](#exercise-1-event-basics)
2. [Exercise 2: Simple Counter Component](#exercise-2-simple-counter-component)
3. [Exercise 3: FIFO Queue Component](#exercise-3-fifo-queue-component)
4. [Exercise 4: Message Echo System](#exercise-4-message-echo-system)
5. [Exercise 5: Simple Cache](#exercise-5-simple-cache)
6. [Exercise 6: Adding Tracing](#exercise-6-adding-tracing)
7. [Exercise 7: Multi-Component System](#exercise-7-multi-component-system)

---

## Exercise 1: Event Basics

**Goal**: Understand events and the simulation engine

**Task**: Create a simulation where multiple events print messages at different times.

### Instructions

Create a file `exercises/ex1_events/main.go`:

```go
package main

import (
    "fmt"
    "github.com/sarchlab/akita/v4/sim"
)

// TODO: Define PrintEvent struct
type PrintEvent struct {
    // Your code here
}

// TODO: Define PrintHandler struct
type PrintHandler struct {
    // Your code here
}

// TODO: Implement Handle method
func (h *PrintHandler) Handle(e sim.Event) error {
    // Your code here
    return nil
}

func main() {
    // TODO: Create engine

    // TODO: Create handler

    // TODO: Schedule events at times: 1.0, 2.5, 5.0, 10.0
    // Each should print a different message

    // TODO: Run simulation

    fmt.Println("Simulation complete!")
}
```

### Solution

<details>
<summary>Click to see solution</summary>

```go
package main

import (
    "fmt"
    "github.com/sarchlab/akita/v4/sim"
)

type PrintEvent struct {
    *sim.EventBase
    message string
}

type PrintHandler struct {
    engine sim.Engine
}

func (h *PrintHandler) Handle(e sim.Event) error {
    evt := e.(*PrintEvent)
    fmt.Printf("[Time %.2f] %s\n", evt.Time(), evt.message)
    return nil
}

func main() {
    engine := sim.NewSerialEngine()
    handler := &PrintHandler{engine: engine}

    events := []struct {
        time    sim.VTimeInSec
        message string
    }{
        {1.0, "First event"},
        {2.5, "Second event"},
        {5.0, "Third event"},
        {10.0, "Final event"},
    }

    for _, evt := range events {
        event := &PrintEvent{
            EventBase: sim.NewEventBase(evt.time, handler),
            message:   evt.message,
        }
        engine.Schedule(event)
    }

    engine.Run()
    fmt.Println("Simulation complete!")
}
```

</details>

### Expected Output

```
[Time 1.00] First event
[Time 2.50] Second event
[Time 5.00] Third event
[Time 10.00] Final event
Simulation complete!
```

---

## Exercise 2: Simple Counter Component

**Goal**: Build a ticking component that counts clock cycles

**Task**: Create a counter component that increments every clock cycle and prints every 10 cycles.

### Instructions

Create `exercises/ex2_counter/counter.go`:

```go
package main

import (
    "fmt"
    "github.com/sarchlab/akita/v4/sim"
)

type Counter struct {
    *sim.TickingComponent

    count     int
    maxCount  int
}

// TODO: Implement Tick() method
func (c *Counter) Tick() bool {
    // Increment counter
    // Print every 10 cycles
    // Stop when reaching maxCount
    // Return true if still counting, false if done
    return false
}

// TODO: Create Builder
type Builder struct {
    engine   sim.Engine
    freq     sim.Freq
    maxCount int
}

func MakeBuilder() Builder {
    return Builder{
        freq:     1 * sim.GHz,
        maxCount: 100,
    }
}

func (b Builder) WithEngine(e sim.Engine) Builder {
    b.engine = e
    return b
}

func (b Builder) WithFreq(f sim.Freq) Builder {
    b.freq = f
    return b
}

func (b Builder) WithMaxCount(n int) Builder {
    b.maxCount = n
    return b
}

func (b Builder) Build(name string) *Counter {
    // TODO: Create and initialize counter
    return nil
}

func main() {
    engine := sim.NewSerialEngine()

    counter := MakeBuilder().
        WithEngine(engine).
        WithFreq(1 * sim.GHz).
        WithMaxCount(50).
        Build("Counter")

    engine.Run()

    fmt.Printf("Final count: %d\n", counter.count)
}
```

### Solution

<details>
<summary>Click to see solution</summary>

```go
package main

import (
    "fmt"
    "github.com/sarchlab/akita/v4/sim"
)

type Counter struct {
    *sim.TickingComponent
    count    int
    maxCount int
}

func (c *Counter) Tick() bool {
    c.count++

    if c.count%10 == 0 {
        fmt.Printf("[Cycle %d] Count: %d\n",
            c.count, c.count)
    }

    if c.count >= c.maxCount {
        return false // Done counting
    }

    return true // Keep counting
}

type Builder struct {
    engine   sim.Engine
    freq     sim.Freq
    maxCount int
}

func MakeBuilder() Builder {
    return Builder{
        freq:     1 * sim.GHz,
        maxCount: 100,
    }
}

func (b Builder) WithEngine(e sim.Engine) Builder {
    b.engine = e
    return b
}

func (b Builder) WithFreq(f sim.Freq) Builder {
    b.freq = f
    return b
}

func (b Builder) WithMaxCount(n int) Builder {
    b.maxCount = n
    return b
}

func (b Builder) Build(name string) *Counter {
    c := &Counter{
        maxCount: b.maxCount,
    }

    c.TickingComponent = sim.NewTickingComponent(
        name, b.engine, b.freq, c)

    return c
}

func main() {
    engine := sim.NewSerialEngine()

    counter := MakeBuilder().
        WithEngine(engine).
        WithFreq(1 * sim.GHz).
        WithMaxCount(50).
        Build("Counter")

    engine.Run()

    fmt.Printf("Final count: %d\n", counter.count)
}
```

</details>

---

## Exercise 3: FIFO Queue Component

**Goal**: Build a component with ports that stores and forwards messages

**Task**: Create a FIFO queue component that receives messages, stores them, and forwards them one per cycle.

### Instructions

Create `exercises/ex3_fifo/fifo.go`:

```go
package main

import (
    "fmt"
    "github.com/sarchlab/akita/v4/sim"
)

// DataMsg carries integer data
type DataMsg struct {
    sim.MsgMeta
    Data int
}

type FIFO struct {
    *sim.TickingComponent

    inPort  sim.Port
    outPort sim.Port

    queue    []sim.Msg
    capacity int
}

func (f *FIFO) Tick() bool {
    madeProgress := false

    // TODO: Try to send from queue to outPort
    // if queue not empty and outPort accepts

    // TODO: Try to receive from inPort to queue
    // if queue not full

    return madeProgress
}

type Builder struct {
    engine   sim.Engine
    freq     sim.Freq
    capacity int
}

func MakeBuilder() Builder {
    return Builder{
        freq:     1 * sim.GHz,
        capacity: 4,
    }
}

// TODO: Implement builder methods

func (b Builder) Build(name string) *FIFO {
    // TODO: Build FIFO with ports
    return nil
}

func main() {
    // TODO: Create FIFO
    // TODO: Create sender and receiver
    // TODO: Connect them
    // TODO: Send some messages
    // TODO: Run simulation
}
```

### Solution

<details>
<summary>Click to see solution</summary>

```go
package main

import (
    "fmt"
    "github.com/sarchlab/akita/v4/sim"
    "github.com/sarchlab/akita/v4/sim/directconnection"
)

type DataMsg struct {
    sim.MsgMeta
    Data int
}

type FIFO struct {
    *sim.TickingComponent
    inPort   sim.Port
    outPort  sim.Port
    queue    []sim.Msg
    capacity int
}

func (f *FIFO) Tick() bool {
    madeProgress := false

    // Try to send from queue
    if len(f.queue) > 0 {
        err := f.outPort.Send(f.queue[0])
        if err == nil {
            fmt.Printf("[%.0f] FIFO sent: %d\n",
                f.Engine().CurrentTime(),
                f.queue[0].(*DataMsg).Data)
            f.queue = f.queue[1:]
            madeProgress = true
        }
    }

    // Try to receive into queue
    if len(f.queue) < f.capacity {
        msg := f.inPort.Retrieve()
        if msg != nil {
            fmt.Printf("[%.0f] FIFO received: %d\n",
                f.Engine().CurrentTime(),
                msg.(*DataMsg).Data)
            f.queue = append(f.queue, msg)
            madeProgress = true
        }
    }

    return madeProgress
}

type Builder struct {
    engine   sim.Engine
    freq     sim.Freq
    capacity int
}

func MakeBuilder() Builder {
    return Builder{
        freq:     1 * sim.GHz,
        capacity: 4,
    }
}

func (b Builder) WithEngine(e sim.Engine) Builder {
    b.engine = e
    return b
}

func (b Builder) WithCapacity(c int) Builder {
    b.capacity = c
    return b
}

func (b Builder) Build(name string) *FIFO {
    f := &FIFO{
        capacity: b.capacity,
        queue:    make([]sim.Msg, 0, b.capacity),
    }

    f.TickingComponent = sim.NewTickingComponent(
        name, b.engine, b.freq, f)

    f.inPort = sim.NewLimitNumMsgPort(f, 1, name+".In")
    f.outPort = sim.NewLimitNumMsgPort(f, 1, name+".Out")

    return f
}

// Simple sender component
type Sender struct {
    *sim.TickingComponent
    port     sim.Port
    toSend   []int
    nextIdx  int
}

func (s *Sender) Tick() bool {
    if s.nextIdx >= len(s.toSend) {
        return false
    }

    msg := &DataMsg{
        MsgMeta: sim.MsgMeta{
            Src: s.port,
            Dst: s.port.Connection().(*directconnection.DirectConnection).PluggedPorts()[1],
        },
        Data: s.toSend[s.nextIdx],
    }

    err := s.port.Send(msg)
    if err == nil {
        s.nextIdx++
        return true
    }

    return false
}

// Simple receiver component
type Receiver struct {
    *sim.TickingComponent
    port     sim.Port
    received []int
}

func (r *Receiver) Tick() bool {
    msg := r.port.Retrieve()
    if msg != nil {
        data := msg.(*DataMsg).Data
        r.received = append(r.received, data)
        fmt.Printf("[%.0f] Receiver got: %d\n",
            r.Engine().CurrentTime(), data)
        return true
    }
    return false
}

func main() {
    engine := sim.NewSerialEngine()

    // Create components
    sender := &Sender{
        toSend: []int{1, 2, 3, 4, 5},
    }
    sender.TickingComponent = sim.NewTickingComponent(
        "Sender", engine, 1*sim.GHz, sender)
    sender.port = sim.NewLimitNumMsgPort(sender, 1, "Sender.Port")

    fifo := MakeBuilder().
        WithEngine(engine).
        WithCapacity(3).
        Build("FIFO")

    receiver := &Receiver{
        received: make([]int, 0),
    }
    receiver.TickingComponent = sim.NewTickingComponent(
        "Receiver", engine, 1*sim.GHz, receiver)
    receiver.port = sim.NewLimitNumMsgPort(receiver, 1, "Receiver.Port")

    // Connect: Sender -> FIFO -> Receiver
    conn1 := directconnection.NewDirectConnection(
        "SenderToFIFO", engine, 1*sim.GHz)
    conn1.PlugIn(sender.port, 8)
    conn1.PlugIn(fifo.inPort, 8)

    conn2 := directconnection.NewDirectConnection(
        "FIFOToReceiver", engine, 1*sim.GHz)
    conn2.PlugIn(fifo.outPort, 8)
    conn2.PlugIn(receiver.port, 8)

    // Run
    engine.Run()

    fmt.Printf("\nReceived %d messages: %v\n",
        len(receiver.received), receiver.received)
}
```

</details>

---

## Exercise 4: Message Echo System

**Goal**: Build a request-response system with two components

**Task**: Create a client that sends ping requests and a server that responds with pong.

### Instructions

```go
package main

import (
    "fmt"
    "github.com/sarchlab/akita/v4/sim"
)

// TODO: Define PingReqMsg
type PingReqMsg struct {
    sim.MsgMeta
    SeqNum int
}

// TODO: Define PingRspMsg
type PingRspMsg struct {
    sim.MsgMeta
    SeqNum int
}

// TODO: Implement Client component
type Client struct {
    *sim.TickingComponent
    port sim.Port

    numToSend    int
    numSent      int
    numReceived  int
    pendingReqs  map[int]bool
}

func (c *Client) Tick() bool {
    // TODO: Send ping requests (up to numToSend)
    // TODO: Receive ping responses
    // TODO: Return false when all responses received
    return false
}

// TODO: Implement Server component
type Server struct {
    *sim.TickingComponent
    port sim.Port
}

func (s *Server) Tick() bool {
    // TODO: Receive ping requests
    // TODO: Send ping responses
    return false
}

func main() {
    // TODO: Build client and server
    // TODO: Connect them
    // TODO: Run simulation
}
```

### Challenge

- Track round-trip latency for each request
- Limit client to only 2 outstanding requests at a time
- Add a processing delay in the server

---

## Exercise 5: Simple Cache

**Goal**: Build a basic cache with hit/miss logic

**Task**: Create a direct-mapped cache that sits between a CPU and memory.

### Specification

- **Capacity**: 4 cache lines
- **Each line**: 64 bytes
- **Direct-mapped**: Address maps to (addr / 64) % 4
- **Cache hit**: Return data immediately (1 cycle)
- **Cache miss**: Fetch from memory (10 cycle latency)

### Instructions

```go
package main

import "github.com/sarchlab/akita/v4/sim"

type CacheLine struct {
    Valid bool
    Tag   uint64
    Data  []byte
}

type Cache struct {
    *sim.TickingComponent

    topPort    sim.Port  // CPU side
    bottomPort sim.Port  // Memory side

    lines       [4]CacheLine
    pendingMiss map[uint64]*MemReadReq

    statsHits   int
    statsMisses int
}

// TODO: Implement Tick()
func (c *Cache) Tick() bool {
    // Handle memory responses
    // Handle CPU requests
    // Update statistics
    return false
}

// TODO: Implement lookup
func (c *Cache) lookup(addr uint64) (*CacheLine, bool) {
    // Return line and hit/miss
    return nil, false
}

// TODO: Implement install
func (c *Cache) install(addr uint64, data []byte) {
    // Install data in cache
}

func main() {
    // TODO: Build cache, CPU model, memory model
    // TODO: Send read requests
    // TODO: Print hit/miss statistics
}
```

---

## Exercise 6: Adding Tracing

**Goal**: Learn to use hooks for observing simulation

**Task**: Add tracing to any previous exercise to log all messages.

### Instructions

```go
package main

import (
    "fmt"
    "os"
    "github.com/sarchlab/akita/v4/sim"
)

// Message tracer hook
type MsgTracer struct {
    file *os.File
}

func NewMsgTracer(filename string) *MsgTracer {
    f, err := os.Create(filename)
    if err != nil {
        panic(err)
    }

    return &MsgTracer{file: f}
}

func (t *MsgTracer) Func(ctx sim.HookCtx) {
    // TODO: Log messages when they are sent/received
    // Format: [time] src -> dst: message_type
}

func (t *MsgTracer) Close() {
    t.file.Close()
}

func main() {
    // TODO: Create components
    // TODO: Create tracer
    // TODO: Install hooks on components
    // TODO: Run simulation
    // TODO: Close tracer
}
```

### Tasks

1. Log all messages with timestamp, src, dst, type
2. Count total messages sent
3. Calculate average message latency
4. Generate a CSV file for analysis

---

## Exercise 7: Multi-Component System

**Goal**: Build a complete system with multiple interacting components

**Task**: Create a simple memory hierarchy: CPU -> L1 Cache -> L2 Cache -> Memory

### System Specification

```
┌─────────┐     ┌──────────┐     ┌──────────┐     ┌────────┐
│   CPU   │────►│ L1 Cache │────►│ L2 Cache │────►│ Memory │
│         │◄────│  (4 KB)  │◄────│ (64 KB)  │◄────│        │
└─────────┘     └──────────┘     └──────────┘     └────────┘
   1 cycle         2 cycle          5 cycle        50 cycle
```

### Requirements

1. **CPU Component**
   - Generates read requests to sequential addresses
   - Tracks latency of each request
   - Prints statistics at end

2. **L1 Cache**
   - 4 KB capacity, direct-mapped
   - 2 cycle hit latency
   - Forwards misses to L2

3. **L2 Cache**
   - 64 KB capacity, 4-way set associative
   - 5 cycle hit latency
   - Forwards misses to Memory

4. **Memory**
   - Always hits
   - 50 cycle latency

### Tasks

1. Implement all four components
2. Connect them properly
3. Run simulation with 1000 memory accesses
4. Report:
   - L1 hit rate
   - L2 hit rate
   - Average memory access latency
   - Total simulation time

### Bonus Challenges

1. Implement LRU replacement in L2 cache
2. Add write support (write-through vs write-back)
3. Model port contention and queuing
4. Visualize access patterns
5. Add prefetching to L2

---

## Testing Your Solutions

For each exercise, add tests:

```go
func TestCounter_CountsCorrectly(t *testing.T) {
    engine := sim.NewSerialEngine()

    counter := MakeBuilder().
        WithEngine(engine).
        WithMaxCount(10).
        Build("TestCounter")

    engine.Run()

    if counter.count != 10 {
        t.Errorf("Expected 10, got %d", counter.count)
    }
}
```

---

## Learning Outcomes

After completing these exercises, you will be able to:

- ✅ Create and schedule events
- ✅ Build ticking components
- ✅ Design ports and messages
- ✅ Connect components with connections
- ✅ Implement request-response protocols
- ✅ Build stateful components (caches, queues)
- ✅ Use hooks for tracing
- ✅ Construct multi-component systems
- ✅ Measure and optimize performance
- ✅ Follow AkitaSim best practices

---

## Next Steps

After these exercises:

1. **Study real components** in `mem/` directory
2. **Explore GPU simulation** in MGPUSim project
3. **Build your own** architecture idea
4. **Contribute** improvements to AkitaSim

Happy coding! 🚀
