# AkitaSim Comprehensive Tutorial

Welcome to AkitaSim! This tutorial will guide you through learning the AkitaSim simulation framework from basics to advanced concepts.

## Table of Contents

1. [Introduction](#1-introduction)
2. [Core Concepts Overview](#2-core-concepts-overview)
3. [Getting Started: Your First Simulation](#3-getting-started-your-first-simulation)
4. [Components in Detail](#4-components-in-detail)
5. [Ports and Messages](#5-ports-and-messages)
6. [Connections](#6-connections)
7. [Hooks for Tracing and Monitoring](#7-hooks-for-tracing-and-monitoring)
8. [Domains](#8-domains)
9. [Building Real Components](#9-building-real-components)
10. [Best Practices and Patterns](#10-best-practices-and-patterns)
11. [Advanced Topics](#11-advanced-topics)

---

## 1. Introduction

### What is AkitaSim?

**AkitaSim** is a computer architecture simulation engine written in Go. Think of it as a "game engine" for hardware simulation - it's not a complete simulator but a **framework** for building simulators. It allows researchers and developers to model and experiment with computer architecture designs.

### Why AkitaSim?

- **Modular**: Build simulations from reusable components
- **Flexible**: Event-driven and cycle-accurate simulation modes
- **Traceable**: Built-in hooks for monitoring and debugging
- **Fast**: Written in Go for performance
- **Clean Architecture**: Clear separation of concerns

### Version Information

- **v4**: Current stable version (event-driven focus)
- **v5**: Next generation (middleware-based, better for snapshots)
- Official docs: https://akitasim.dev/docs/akita/

---

## 2. Core Concepts Overview

AkitaSim has several key abstractions that work together:

```
┌─────────────────────────────────────────────────┐
│              SIMULATION                          │
│  ┌────────────────────────────────────────┐    │
│  │         ENGINE                          │    │
│  │  (manages time & event queue)           │    │
│  └────────────────────────────────────────┘    │
│                                                  │
│  ┌──────────────┐        ┌──────────────┐      │
│  │  COMPONENT   │◄──────►│  COMPONENT   │      │
│  │              │        │              │      │
│  │  ┌────────┐  │        │  ┌────────┐  │      │
│  │  │ Port   │  │        │  │ Port   │  │      │
│  │  └────┬───┘  │        │  └───┬────┘  │      │
│  └───────┼──────┘        └──────┼───────┘      │
│          │                      │               │
│          │   ┌──────────────┐   │               │
│          └───┤ CONNECTION   ├───┘               │
│              │ (DirectLink) │                   │
│              └──────────────┘                   │
│                                                  │
│  ┌──────────────────────────────────────────┐  │
│  │        HOOKS (tracing/monitoring)         │  │
│  └──────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
```

### Key Terms

| Term | Description | File Location |
|------|-------------|---------------|
| **Component** | Hardware elements (CPU, cache, memory) | `sim/component.go` |
| **Port** | Communication endpoints for components | `sim/port.go` |
| **Connection** | Communication fabric between ports | `sim/connection.go` |
| **Message** | Data transferred between components | `sim/msg.go` |
| **Event** | Something that happens at a specific time | `sim/event.go` |
| **Engine** | Drives simulation time and events | `sim/engine.go` |
| **Hook** | Callback for tracing/monitoring | `sim/hook.go` |
| **Domain** | Logical grouping of components | `sim/domain.go` |
| **Middleware** | Tick-based logic (v5 pattern) | `sim/middleware.go` |

---

## 3. Getting Started: Your First Simulation

Let's start with the simplest possible simulation to understand the basics.

### Example 1: Print Event

**Location**: `examples/01_print_event/main.go`

This example shows the absolute basics: creating an engine, scheduling an event, and handling it.

```go
package main

import (
    "fmt"
    "github.com/sarchlab/akita/v4/sim"
)

// PrintEvent is a simple event that prints a message
type PrintEvent struct {
    *sim.EventBase
    message string
}

// PrintHandler handles PrintEvent
type PrintHandler struct{}

func (h PrintHandler) Handle(e sim.Event) error {
    evt := e.(*PrintEvent)
    fmt.Printf("[Time %.2f] %s\n", evt.Time(), evt.message)
    return nil
}

func main() {
    // 1. Create the simulation engine
    engine := sim.NewSerialEngine()

    // 2. Create a handler
    handler := PrintHandler{}

    // 3. Create and schedule an event
    event := &PrintEvent{
        EventBase: sim.NewEventBase(5, handler),
        message:   "Hello, AkitaSim!",
    }
    engine.Schedule(event)

    // 4. Run the simulation
    engine.Run()

    fmt.Println("Simulation complete!")
}
```

**What's happening:**
1. Create a `SerialEngine` to manage simulation time
2. Define an event type with a time and data
3. Define a handler that processes the event
4. Schedule the event at time 5
5. Run the simulation - the engine processes all events in time order

**Try it:**
```bash
cd examples/01_print_event
go run main.go
```

### Example 2: Cell Splitting

**Location**: `examples/02_cell_split/main.go`

This shows how events can schedule other events, creating a dynamic simulation.

```go
// A cell that splits into two cells after a delay
type CellSplitEvent struct {
    *sim.EventBase
    cellCount int
}

func (h *CellHandler) Handle(e sim.Event) error {
    evt := e.(*CellSplitEvent)
    newCount := evt.cellCount * 2

    fmt.Printf("[Time %.0f] Cell split: %d -> %d cells\n",
        evt.Time(), evt.cellCount, newCount)

    // Schedule next split in 1 time unit
    if newCount < 16 {
        nextEvent := &CellSplitEvent{
            EventBase: sim.NewEventBase(evt.Time() + 1, h),
            cellCount: newCount,
        }
        h.engine.Schedule(nextEvent)
    }

    return nil
}
```

**Key insight:** Events can schedule new events, creating chains of simulation behavior.

---

## 4. Components in Detail

Components are the heart of AkitaSim - they represent hardware elements.

### Component Interface

Every component must implement:

```go
type Component interface {
    Named           // Has a Name() string method
    Handler         // Can Handle(Event) error
    Hookable        // Can AcceptHook(Hook)
    PortOwner       // Owns Ports
}
```

### Two Patterns for Components

#### Pattern 1: Event-Driven Components

**When to use:** Discrete actions at specific times (network switches, controllers)

**Example:** `examples/ping/comp.go`

```go
type Comp struct {
    *sim.ComponentBase  // Provides basic component functionality

    topPort  sim.Port   // Communication port
    state    int        // Component state
}

// Handle processes events
func (c *Comp) Handle(e sim.Event) error {
    switch req := e.(type) {
    case *PingReqMsg:
        return c.handlePingReq(req)
    case *PingRspMsg:
        return c.handlePingRsp(req)
    }
    return nil
}

func (c *Comp) handlePingReq(req *PingReqMsg) error {
    // Process request
    rsp := &PingRspMsg{
        MsgMeta: sim.MsgMeta{
            Src: c.topPort,
            Dst: req.Src,
        },
    }

    // Send response
    c.topPort.Send(rsp)
    return nil
}
```

**Key points:**
- Embeds `ComponentBase` for common functionality
- Implements `Handle(Event)` to process events
- Events can be messages or custom types
- Component reacts only when events occur

#### Pattern 2: Ticking Components

**When to use:** Hardware with regular clock cycles (CPUs, caches, memories)

**Example:** `examples/tickingping/comp.go`

```go
type Comp struct {
    *sim.TickingComponent       // Provides ticking infrastructure
    sim.MiddlewareHolder        // Holds middleware chain (v5)

    topPort  sim.Port
    ctrlPort sim.Port

    state    CompState
}

// Tick runs every clock cycle
func (c *Comp) Tick() bool {
    madeProgress := false

    // Try to send pending messages
    if c.state.pendingRsp != nil {
        err := c.topPort.Send(c.state.pendingRsp)
        if err == nil {
            c.state.pendingRsp = nil
            madeProgress = true
        }
    }

    // Try to receive new messages
    msg := c.topPort.Retrieve()
    if msg != nil {
        c.processMessage(msg)
        madeProgress = true
    }

    return madeProgress
}
```

**Key points:**
- Embeds `TickingComponent`
- Implements `Tick() bool` - called every clock cycle
- Returns `true` if work was done (keeps simulation active)
- Returns `false` if idle (simulation can end if all components idle)
- More suitable for cycle-accurate hardware modeling

### Component Lifecycle

```go
// 1. Create using Builder pattern
builder := NewCompBuilder().
    WithFreq(1 * sim.GHz).
    WithEngine(engine)

comp := builder.Build("MyComponent")

// 2. Component is now ready
// - It's registered with the engine
// - It can send/receive messages via ports
// - It will receive Tick calls or Events

// 3. During simulation
// - Event-driven: Handle() called when events arrive
// - Ticking: Tick() called every cycle

// 4. Cleanup (if needed)
// - Usually automatic when simulation ends
```

---

## 5. Ports and Messages

Ports are how components communicate. Think of them as network interfaces.

### Port Basics

```go
type Port interface {
    Named
    Hookable

    // Sending
    Send(Msg) error

    // Receiving (component side)
    Retrieve() Msg
    Peek() Msg

    // Receiving (connection side)
    Deliver(Msg)

    // Connection management
    SetConnection(Connection)
    Connection() Connection
}
```

### Creating Ports

```go
// Create a port for a component
port := sim.NewLimitNumMsgPort(
    comp,                  // Owner component
    4,                     // Buffer size (max 4 messages)
    "ComponentPort",       // Port name
)

// Associate with component
comp.topPort = port
```

### Message Structure

All messages implement:

```go
type Msg interface {
    Meta() *MsgMeta
    Clone() Msg
}

type MsgMeta struct {
    ID            string
    Src           Port      // Source port
    Dst           Port      // Destination port
    SendTime      VTimeInSec
    RecvTime      VTimeInSec
    TrafficClass  int
    TrafficBytes  int
}
```

### Creating Custom Messages

```go
// Define your message type
type PingReqMsg struct {
    sim.MsgMeta
    SeqNum int
    Data   string
}

// Send it
msg := &PingReqMsg{
    MsgMeta: sim.MsgMeta{
        Src: myComponent.port,
        Dst: otherComponent.port,
    },
    SeqNum: 1,
    Data:   "ping",
}

err := myComponent.port.Send(msg)
if err != nil {
    // Port buffer full, try again later
}
```

### Receiving Messages

```go
// In component's Tick() or Handle()
msg := c.topPort.Retrieve()
if msg != nil {
    // Got a message!
    switch m := msg.(type) {
    case *PingReqMsg:
        // Handle ping request
    case *PingRspMsg:
        // Handle ping response
    }
}
```

### Port Buffer Management

Ports have limited buffers:

```go
// Limited buffer
port := sim.NewLimitNumMsgPort(comp, 4, "port")  // Max 4 messages

// Unlimited buffer (use carefully!)
port := sim.NewUnlimitedNumMsgPort(comp, "port")

// Send can fail if buffer full
err := port.Send(msg)
if err != nil {
    // Buffer full - need to retry later
    // Component should save msg and try again next tick
}
```

---

## 6. Connections

Connections model the communication fabric between components (wires, buses, networks).

### Connection Types

#### DirectConnection (Zero Latency)

**Location:** `sim/directconnection/directconnection.go`

Simplest connection - messages arrive immediately.

```go
// Create connection
conn := directconnection.NewDirectConnection(
    "MyConnection",
    engine,
    1*sim.GHz,  // Frequency
)

// Plug in ports
conn.PlugIn(component1.port, 8)  // 8 bytes/cycle bandwidth
conn.PlugIn(component2.port, 8)
```

**When to use:**
- Components on same chip
- When latency doesn't matter
- Testing and prototyping

#### Network-on-Chip (With Latency)

**Location:** `noc/`

Models realistic networks with switches, routers, latency, and bandwidth.

```go
// More complex - see noc/acceptance/one_to_one/main.go
// Involves:
// - Switches/routers
// - Links with bandwidth/latency
// - Flow control
// - Routing algorithms
```

### Setting Up Connections

```go
// 1. Create components
comp1 := builder1.Build("Comp1")
comp2 := builder2.Build("Comp2")

// 2. Create their ports
port1 := sim.NewLimitNumMsgPort(comp1, 4, "Comp1.Port")
port2 := sim.NewLimitNumMsgPort(comp2, 4, "Comp2.Port")

// 3. Create connection
conn := directconnection.NewDirectConnection("Link", engine, 1*sim.GHz)

// 4. Plug ports into connection
conn.PlugIn(port1, 64)  // 64 bytes/cycle bandwidth
conn.PlugIn(port2, 64)

// Now comp1 and comp2 can communicate!
```

### Message Flow Through Connection

```
Component A                  Connection                  Component B
┌─────────┐                 ┌──────────┐                ┌─────────┐
│         │                 │          │                │         │
│  Port   ├────── Send ────►│          ├──── Deliver ──►│  Port   │
│ OutBuf  │                 │  Fabric  │                │  InBuf  │
│         │◄─── Retrieve ───┤          │◄─── Send ──────┤         │
│         │                 │          │                │         │
└─────────┘                 └──────────┘                └─────────┘
```

1. Component A calls `port.Send(msg)`
2. Message goes to port's outgoing buffer
3. Connection retrieves from outgoing buffer
4. Connection delivers to destination port's incoming buffer
5. Component B calls `port.Retrieve()` to get message

---

## 7. Hooks for Tracing and Monitoring

Hooks let you observe simulation internals without modifying component code.

### What Are Hooks?

Hooks are callbacks triggered at specific points:

```go
type Hook interface {
    Func(HookCtx)
}

type HookCtx struct {
    Domain string       // Where it happened
    Pos    HookPos     // When it happened
    Item   interface{} // What happened
    Detail interface{} // Additional context
}
```

### Hook Positions

```go
const (
    HookPosBeforeEvent  // Before handling an event
    HookPosAfterEvent   // After handling an event
    HookPosPortMsgSend  // When message sent to port
    HookPosPortMsgRecvd // When message received from port
    HookPosBufPush      // When data pushed to buffer
    HookPosBufPop       // When data popped from buffer
    // ... and more
)
```

### Creating a Simple Hook

```go
// Message tracer hook
type MsgTracer struct {
    file *os.File
}

func (t *MsgTracer) Func(ctx sim.HookCtx) {
    if ctx.Pos == sim.HookPosPortMsgSend {
        msg := ctx.Item.(sim.Msg)
        fmt.Fprintf(t.file, "[%.2f] %s -> %s: %T\n",
            ctx.Domain,
            msg.Meta().Src.Name(),
            msg.Meta().Dst.Name(),
            msg,
        )
    }
}

// Install hook
tracer := &MsgTracer{file: os.Stdout}
component.AcceptHook(tracer)
```

### Using Built-in Tracing

AkitaSim has a comprehensive tracing system in `tracing/`:

```go
import "github.com/sarchlab/akita/v4/tracing"

// Create trace file
tracer := tracing.NewMySQLTracer()
tracer.Init()
defer tracer.Fini()

// Add hooks to components
tracing.CollectTrace(component, tracer)

// Run simulation - traces are automatically collected

// Analyze traces
// - SQL queries on trace database
// - Visualization tools
// - Performance analysis
```

### Hook Best Practices

1. **Use for observation only** - Don't modify simulation state in hooks
2. **Keep hooks fast** - They're called frequently
3. **Filter appropriately** - Only hook what you need
4. **Use built-in tracing** - Leverage existing infrastructure

---

## 8. Domains

Domains organize components hierarchically.

### What Are Domains?

Domains are named containers for related components:

```go
type Domain struct {
    *sim.PortOwnerBase
    name string
}
```

### Why Use Domains?

- **Organization**: Group related components (e.g., "GPU.ComputeUnit0")
- **Tracing**: Easier to identify component locations
- **Hierarchy**: Model chip/die/unit relationships

### Creating Domains

```go
// Create a domain for a GPU
gpu := sim.NewDomain("GPU")

// Create components within domain
for i := 0; i < 64; i++ {
    name := fmt.Sprintf("GPU.CU%d", i)
    cu := builder.Build(name)
    // Component name includes domain
}
```

### Domain Naming Convention

```
ChipName.SubUnit.Component
└─────┬──────┘ └──┬───┘ └───┬──┘
   Domain      Domain    Component

Example: "GPU.ComputeUnit0.VectorALU"
```

---

## 9. Building Real Components

Let's build a realistic component: a simple cache.

### Simple Cache Example

```go
package simplecache

import "github.com/sarchlab/akita/v4/sim"

// Cache stores data with tags
type CacheLine struct {
    Valid bool
    Tag   uint64
    Data  []byte
}

// Cache component
type Cache struct {
    *sim.TickingComponent

    topPort    sim.Port  // CPU-facing port
    bottomPort sim.Port  // Memory-facing port

    lines      []CacheLine
    numSets    int
    numWays    int
    blockSize  int

    pendingReqs map[uint64]*MemReadReq
}

// Tick handles one clock cycle
func (c *Cache) Tick() bool {
    madeProgress := false

    // Handle responses from memory
    madeProgress = c.handleMemoryResponse() || madeProgress

    // Handle requests from CPU
    madeProgress = c.handleCPURequest() || madeProgress

    return madeProgress
}

func (c *Cache) handleCPURequest() bool {
    req := c.topPort.Retrieve()
    if req == nil {
        return false
    }

    readReq := req.(*MemReadReq)

    // Check cache
    line := c.lookup(readReq.Address)

    if line != nil && line.Valid {
        // Cache hit!
        rsp := &MemReadRsp{
            MsgMeta: sim.MsgMeta{
                Src: c.topPort,
                Dst: readReq.Src,
            },
            Data: line.Data,
        }
        c.topPort.Send(rsp)
        return true
    }

    // Cache miss - forward to memory
    memReq := &MemReadReq{
        MsgMeta: sim.MsgMeta{
            Src: c.bottomPort,
            Dst: c.bottomPort.Connection().PluggedPorts()[1],
        },
        Address: readReq.Address,
    }

    c.bottomPort.Send(memReq)
    c.pendingReqs[readReq.Address] = readReq

    return true
}

func (c *Cache) handleMemoryResponse() bool {
    rsp := c.bottomPort.Retrieve()
    if rsp == nil {
        return false
    }

    memRsp := rsp.(*MemReadRsp)

    // Install in cache
    c.install(memRsp.Address, memRsp.Data)

    // Forward to CPU
    origReq := c.pendingReqs[memRsp.Address]
    cpuRsp := &MemReadRsp{
        MsgMeta: sim.MsgMeta{
            Src: c.topPort,
            Dst: origReq.Src,
        },
        Data: memRsp.Data,
    }

    c.topPort.Send(cpuRsp)
    delete(c.pendingReqs, memRsp.Address)

    return true
}

func (c *Cache) lookup(addr uint64) *CacheLine {
    // Simple direct-mapped lookup
    set := (addr / uint64(c.blockSize)) % uint64(c.numSets)
    line := &c.lines[set]

    tag := addr / uint64(c.blockSize*c.numSets)
    if line.Tag == tag && line.Valid {
        return line
    }

    return nil
}

func (c *Cache) install(addr uint64, data []byte) {
    set := (addr / uint64(c.blockSize)) % uint64(c.numSets)
    tag := addr / uint64(c.blockSize*c.numSets)

    c.lines[set].Valid = true
    c.lines[set].Tag = tag
    c.lines[set].Data = data
}
```

### Builder Pattern for Cache

```go
type Builder struct {
    engine    sim.Engine
    freq      sim.Freq
    numSets   int
    numWays   int
    blockSize int
}

func MakeBuilder() Builder {
    return Builder{
        freq:      1 * sim.GHz,
        numSets:   64,
        numWays:   4,
        blockSize: 64,
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

func (b Builder) WithNumSets(n int) Builder {
    b.numSets = n
    return b
}

func (b Builder) Build(name string) *Cache {
    c := &Cache{
        numSets:     b.numSets,
        numWays:     b.numWays,
        blockSize:   b.blockSize,
        pendingReqs: make(map[uint64]*MemReadReq),
    }

    c.TickingComponent = sim.NewTickingComponent(
        name, b.engine, b.freq, c)

    c.topPort = sim.NewLimitNumMsgPort(c, 4, name+".TopPort")
    c.bottomPort = sim.NewLimitNumMsgPort(c, 4, name+".BottomPort")

    c.lines = make([]CacheLine, b.numSets*b.numWays)

    return c
}
```

### Using the Cache

```go
func main() {
    engine := sim.NewSerialEngine()

    // Build cache
    cache := simplecache.MakeBuilder().
        WithEngine(engine).
        WithFreq(2 * sim.GHz).
        WithNumSets(128).
        Build("L1Cache")

    // Build memory
    memory := idealmemcontroller.MakeBuilder().
        WithEngine(engine).
        WithLatency(100).  // 100ns
        Build("MainMemory")

    // Connect them
    conn := directconnection.NewDirectConnection("CacheMem", engine, 1*sim.GHz)
    conn.PlugIn(cache.bottomPort, 64)
    conn.PlugIn(memory.topPort, 64)

    // Now cache can talk to memory!
}
```

---

## 10. Best Practices and Patterns

### Component Design

1. **Use Builder Pattern**
   - Fluent configuration API
   - Clear default values
   - Type-safe construction

2. **Separate Spec and State** (v5)
   - Spec: Immutable configuration (primitives only)
   - State: Mutable runtime data
   - Makes testing and checkpointing easier

3. **Keep Tick() Simple**
   - One clear task per tick
   - Return early if no work
   - Use middlewares to organize logic

### Message Design

1. **Embed MsgMeta**
   ```go
   type MyMsg struct {
       sim.MsgMeta  // Embedding, not pointer
       // ... fields
   }
   ```

2. **Use Request/Response Pairs**
   ```go
   type ReadReq struct { ... }
   type ReadRsp struct {
       RspTo string  // ID of request
       // ...
   }
   ```

3. **Set Traffic Bytes**
   ```go
   msg.TrafficBytes = len(data)  // For bandwidth modeling
   ```

### Performance Tips

1. **Use ParallelEngine for Large Simulations**
   ```go
   engine := sim.NewParallelEngine()
   ```

2. **Limit Port Buffers**
   - Prevents memory bloat
   - Forces backpressure modeling

3. **Profile First**
   - Use Go pprof
   - Trace only what you need

### Testing

1. **Unit Test Components**
   ```go
   func TestCache_Hit(t *testing.T) {
       engine := sim.NewSerialEngine()
       cache := MakeBuilder().WithEngine(engine).Build("cache")

       // Send request
       // Assert response
   }
   ```

2. **Use Mock Components**
   ```go
   type MockMemory struct {
       responses map[uint64][]byte
   }
   ```

3. **Test Builders**
   ```go
   func TestBuilder_Defaults(t *testing.T) {
       cache := MakeBuilder().Build("test")
       assert.Equal(t, 64, cache.numSets)
   }
   ```

---

## 11. Advanced Topics

### Middleware Pattern (v5)

Middlewares organize component logic as a pipeline:

```go
type Middleware interface {
    Tick() bool
}

// Component with middlewares
type Cache struct {
    *sim.TickingComponent
    sim.MiddlewareHolder
    // ...
}

func (c *Cache) Tick() bool {
    madeProgress := false

    // Run each middleware in order
    for _, mw := range c.Middlewares() {
        madeProgress = mw.Tick() || madeProgress
    }

    return madeProgress
}

// Example middlewares
type WritebackMiddleware struct {
    cache *Cache
}

func (m *WritebackMiddleware) Tick() bool {
    // Handle writeback logic
}

type EvictionMiddleware struct {
    cache *Cache
}

func (m *EvictionMiddleware) Tick() bool {
    // Handle eviction logic
}
```

### Frequency Domains

Components can run at different frequencies:

```go
cpu := builder.WithFreq(3 * sim.GHz).Build("CPU")
memory := builder.WithFreq(800 * sim.MHz).Build("Memory")

// Engine handles frequency conversion automatically
```

### Snapshots and Checkpointing (v5)

v5 architecture enables deterministic snapshots:

```go
// Save component state
state := component.SaveState()

// Restore later
component.RestoreState(state)
```

### Tracing Analysis

Analyze traces with SQL:

```go
// Example queries on trace database
SELECT AVG(latency) FROM messages WHERE type='MemRead';
SELECT component, COUNT(*) FROM events GROUP BY component;
```

### Network-on-Chip

Building complex networks:

```go
// Create switches
switches := make([]*switch.Switch, 16)
for i := range switches {
    switches[i] = switch.MakeBuilder().Build(fmt.Sprintf("Switch%d", i))
}

// Create links
for i := 0; i < 16; i++ {
    link := link.MakeBuilder().
        WithLatency(10).
        WithBandwidth(64).
        Build(fmt.Sprintf("Link%d", i))

    link.Connect(switches[i], switches[(i+1)%16])
}
```

---

## Learning Path

### Beginner (Week 1)

1. ✅ Read "Introduction" and "Core Concepts"
2. ✅ Run `examples/01_print_event`
3. ✅ Run `examples/02_cell_split`
4. ✅ Modify examples to schedule different events
5. ✅ Read "Components" and "Ports and Messages"

### Intermediate (Week 2)

1. ✅ Study `examples/ping/` (event-driven)
2. ✅ Study `examples/tickingping/` (ticking)
3. ✅ Build your own simple component (e.g., counter, FIFO)
4. ✅ Connect two custom components
5. ✅ Add tracing hooks

### Advanced (Week 3-4)

1. ✅ Study `mem/idealmemcontroller/`
2. ✅ Study `mem/cache/` for realistic cache
3. ✅ Build a mini memory hierarchy
4. ✅ Implement custom middleware
5. ✅ Analyze traces with SQL

### Expert

1. ✅ Explore NoC (Network-on-Chip) in `noc/`
2. ✅ Study GPU simulation in MGPUSim repository
3. ✅ Contribute to AkitaSim
4. ✅ Build complete system simulator

---

## Key Files Reference

### Essential Reading

| Purpose | File Path |
|---------|-----------|
| **Examples (Start Here!)** | |
| Print event | `examples/01_print_event/main.go` |
| Cell splitting | `examples/02_cell_split/main.go` |
| Event-driven ping | `examples/ping/example_test.go` |
| Ticking ping | `examples/tickingping/example_test.go` |
| **Core Framework** | |
| Component interface | `sim/component.go` |
| Ticking component | `sim/ticker.go` |
| Port interface | `sim/port.go` |
| Message interface | `sim/msg.go` |
| Connection interface | `sim/connection.go` |
| Engine | `sim/engine.go` |
| Hooks | `sim/hook.go` |
| **Real Components** | |
| Ideal memory | `mem/idealmemcontroller/comp.go` |
| Cache | `mem/cache/` |
| **Infrastructure** | |
| Tracing | `tracing/` |
| NoC | `noc/` |
| **Documentation** | |
| README | `README.md` |
| Migration guide | `migration_guide.md` |

---

## Common Patterns Cheat Sheet

### Creating a Component

```go
type MyComp struct {
    *sim.TickingComponent
    port sim.Port
    // state fields
}

func (c *MyComp) Tick() bool {
    // 1. Handle incoming messages
    msg := c.port.Retrieve()
    if msg != nil {
        // process msg
        return true
    }

    // 2. Send outgoing messages
    if c.hasPendingWork() {
        err := c.port.Send(response)
        if err == nil {
            return true
        }
    }

    return false
}
```

### Builder Pattern

```go
type Builder struct {
    engine sim.Engine
    freq   sim.Freq
}

func MakeBuilder() Builder {
    return Builder{freq: 1 * sim.GHz}
}

func (b Builder) WithEngine(e sim.Engine) Builder {
    b.engine = e
    return b
}

func (b Builder) Build(name string) *MyComp {
    c := &MyComp{}
    c.TickingComponent = sim.NewTickingComponent(name, b.engine, b.freq, c)
    c.port = sim.NewLimitNumMsgPort(c, 4, name+".Port")
    return c
}
```

### Connecting Components

```go
conn := directconnection.NewDirectConnection("link", engine, freq)
conn.PlugIn(comp1.port, 64)
conn.PlugIn(comp2.port, 64)
```

### Adding Hooks

```go
type MyHook struct{}

func (h *MyHook) Func(ctx sim.HookCtx) {
    if ctx.Pos == sim.HookPosPortMsgSend {
        msg := ctx.Item.(sim.Msg)
        fmt.Printf("Message sent: %T\n", msg)
    }
}

component.AcceptHook(&MyHook{})
```

---

## Next Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/sarchlab/akita
   cd akita
   ```

2. **Run examples**
   ```bash
   cd examples/01_print_event
   go run main.go
   ```

3. **Read the code**
   - Start with `sim/component.go`
   - Follow example components

4. **Build something**
   - Start small (counter, FIFO)
   - Add complexity gradually

5. **Join the community**
   - Official docs: https://akitasim.dev
   - GitHub issues for questions

---

## Glossary

- **Component**: A simulated hardware element
- **Port**: Communication endpoint for a component
- **Message**: Data sent between components
- **Connection**: Communication fabric linking ports
- **Event**: Something that happens at a specific time
- **Engine**: Drives simulation time and event processing
- **Hook**: Callback for observing simulation
- **Domain**: Logical grouping of components
- **Middleware**: Reusable tick logic (v5)
- **Builder**: Pattern for constructing components
- **Tick**: One clock cycle of execution
- **Freq**: Frequency (clock rate) of a component

---

## Troubleshooting

### "Port buffer full"
- Increase buffer size: `NewLimitNumMsgPort(comp, 8, name)`
- Add backpressure handling in component

### "Simulation not terminating"
- Check all Tick() methods return false when idle
- Ensure no infinite event loops
- Add timeout in engine.Run()

### "Messages not arriving"
- Verify ports are connected: `port.Connection() != nil`
- Check Src and Dst are set correctly
- Ensure both ports plugged into same connection

### "Performance issues"
- Use ParallelEngine for large sims
- Reduce tracing (only trace what you need)
- Profile with pprof

---

**Happy Simulating! 🚀**

For questions and support:
- Docs: https://akitasim.dev
- GitHub: https://github.com/sarchlab/akita
