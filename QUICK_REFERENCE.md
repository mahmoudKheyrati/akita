# AkitaSim Quick Reference Guide

A condensed reference for AkitaSim APIs and patterns.

## Table of Contents

- [Component Patterns](#component-patterns)
- [Port Operations](#port-operations)
- [Message Types](#message-types)
- [Connection Setup](#connection-setup)
- [Engine Operations](#engine-operations)
- [Hook Usage](#hook-usage)
- [Builder Pattern](#builder-pattern)
- [Common Imports](#common-imports)
- [Middleware (v5)](#middleware-v5)
- [Debugging Tips](#debugging-tips)

---

## Component Patterns

### Event-Driven Component

```go
type MyComp struct {
    *sim.ComponentBase
    port sim.Port
}

func (c *MyComp) Handle(e sim.Event) error {
    switch evt := e.(type) {
    case *MyMsg:
        return c.handleMyMsg(evt)
    }
    return nil
}
```

**Use when**: Discrete events, no regular clock

### Ticking Component

```go
type MyComp struct {
    *sim.TickingComponent
    port sim.Port
}

func (c *MyComp) Tick() bool {
    madeProgress := false

    // Try to send
    if c.hasPendingWork() {
        err := c.port.Send(msg)
        if err == nil {
            madeProgress = true
        }
    }

    // Try to receive
    msg := c.port.Retrieve()
    if msg != nil {
        c.process(msg)
        madeProgress = true
    }

    return madeProgress
}
```

**Use when**: Clock-based hardware (CPU, cache, memory)

---

## Port Operations

### Creating Ports

```go
// Limited buffer (recommended)
port := sim.NewLimitNumMsgPort(component, bufferSize, "PortName")

// Unlimited buffer (use carefully)
port := sim.NewUnlimitedNumMsgPort(component, "PortName")
```

### Sending Messages

```go
msg := &MyMsg{
    MsgMeta: sim.MsgMeta{
        Src: myPort,
        Dst: destPort,
    },
    Data: value,
}

err := myPort.Send(msg)
if err != nil {
    // Buffer full - retry later
}
```

### Receiving Messages

```go
// Retrieve and remove from buffer
msg := port.Retrieve()
if msg != nil {
    // Process message
}

// Peek without removing
msg := port.Peek()
```

### Port Status

```go
// Check connection
if port.Connection() == nil {
    // Not connected
}
```

---

## Message Types

### Basic Message

```go
type MyMsg struct {
    sim.MsgMeta
    Data interface{}
}

// Usage
msg := &MyMsg{
    MsgMeta: sim.MsgMeta{
        Src:          srcPort,
        Dst:          dstPort,
        TrafficBytes: dataSize,
    },
    Data: myData,
}
```

### Request-Response Pattern

```go
type ReadReq struct {
    sim.MsgMeta
    Address uint64
}

type ReadRsp struct {
    sim.MsgMeta
    RspTo string  // Request ID
    Data  []byte
}

// In handler
func (c *Comp) handleReadReq(req *ReadReq) error {
    rsp := &ReadRsp{
        MsgMeta: sim.MsgMeta{
            Src: c.port,
            Dst: req.Src,
        },
        RspTo: req.ID,
        Data:  c.read(req.Address),
    }
    return c.port.Send(rsp)
}
```

---

## Connection Setup

### DirectConnection (Zero Latency)

```go
import "github.com/sarchlab/akita/v4/sim/directconnection"

conn := directconnection.NewDirectConnection(
    "ConnectionName",
    engine,
    1 * sim.GHz,  // Frequency
)

// Connect ports
conn.PlugIn(port1, 64)  // 64 bytes/cycle bandwidth
conn.PlugIn(port2, 64)
```

### Connecting Multiple Components

```go
// A <-> B
connAB := directconnection.NewDirectConnection("AB", engine, freq)
connAB.PlugIn(compA.port, bandwidth)
connAB.PlugIn(compB.port, bandwidth)

// B <-> C
connBC := directconnection.NewDirectConnection("BC", engine, freq)
connBC.PlugIn(compB.port2, bandwidth)
connBC.PlugIn(compC.port, bandwidth)
```

---

## Engine Operations

### Creating Engines

```go
// Serial (single-threaded)
engine := sim.NewSerialEngine()

// Parallel (multi-threaded)
engine := sim.NewParallelEngine()
```

### Scheduling Events

```go
// Create event
event := &MyEvent{
    EventBase: sim.NewEventBase(
        targetTime,  // When to execute
        handler,     // Who handles it
    ),
    data: myData,
}

// Schedule
engine.Schedule(event)
```

### Running Simulation

```go
// Run until no more events/ticks
engine.Run()

// Run until specific time
engine.RunTil(100 * sim.NS)

// Check current time
currentTime := engine.CurrentTime()
```

### Pausing Components

```go
// Pause a component
engine.Pause(component)

// Resume a component
engine.Continue(component)
```

---

## Hook Usage

### Creating a Hook

```go
type MyHook struct {
    // Optional state
}

func (h *MyHook) Func(ctx sim.HookCtx) {
    switch ctx.Pos {
    case sim.HookPosPortMsgSend:
        msg := ctx.Item.(sim.Msg)
        fmt.Printf("Sent: %T from %s to %s\n",
            msg,
            msg.Meta().Src.Name(),
            msg.Meta().Dst.Name())

    case sim.HookPosPortMsgRecvd:
        msg := ctx.Item.(sim.Msg)
        fmt.Printf("Received: %T\n", msg)
    }
}
```

### Installing Hooks

```go
// Create hook
hook := &MyHook{}

// Install on component
component.AcceptHook(hook)

// Install on port
port.AcceptHook(hook)

// Install on connection
connection.AcceptHook(hook)
```

### Common Hook Positions

```go
sim.HookPosBeforeEvent   // Before event handling
sim.HookPosAfterEvent    // After event handling
sim.HookPosPortMsgSend   // When message sent
sim.HookPosPortMsgRecvd  // When message received
sim.HookPosBufPush       // Buffer push
sim.HookPosBufPop        // Buffer pop
```

---

## Builder Pattern

### Standard Builder Structure

```go
type Builder struct {
    engine    sim.Engine
    freq      sim.Freq
    // Configuration fields
    capacity  int
    latency   sim.VTimeInSec
}

func MakeBuilder() Builder {
    return Builder{
        freq:     1 * sim.GHz,
        capacity: 4,
        latency:  10,
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

func (b Builder) WithCapacity(c int) Builder {
    b.capacity = c
    return b
}

func (b Builder) Build(name string) *MyComp {
    c := &MyComp{
        capacity: b.capacity,
    }

    c.TickingComponent = sim.NewTickingComponent(
        name, b.engine, b.freq, c)

    c.port = sim.NewLimitNumMsgPort(c, 4, name+".Port")

    return c
}
```

### Using Builder

```go
comp := MakeBuilder().
    WithEngine(engine).
    WithFreq(2 * sim.GHz).
    WithCapacity(8).
    Build("MyComponent")
```

---

## Common Imports

```go
import (
    "github.com/sarchlab/akita/v4/sim"
    "github.com/sarchlab/akita/v4/sim/directconnection"
    "github.com/sarchlab/akita/v4/tracing"
)

// For memory components
import (
    "github.com/sarchlab/akita/v4/mem"
    "github.com/sarchlab/akita/v4/mem/cache"
    "github.com/sarchlab/akita/v4/mem/idealmemcontroller"
)

// For network-on-chip
import (
    "github.com/sarchlab/akita/v4/noc"
)
```

---

## Middleware (v5)

### Defining Middleware

```go
type MyMiddleware struct {
    component *MyComp
}

func (m *MyMiddleware) Tick() bool {
    // Perform one unit of work
    // Return true if made progress
    return false
}
```

### Component with Middlewares

```go
type MyComp struct {
    *sim.TickingComponent
    sim.MiddlewareHolder

    port sim.Port
}

func (c *MyComp) Tick() bool {
    madeProgress := false

    for _, mw := range c.Middlewares() {
        madeProgress = mw.Tick() || madeProgress
    }

    return madeProgress
}

// In builder
func (b Builder) Build(name string) *MyComp {
    c := &MyComp{}

    // ... initialize ...

    // Add middlewares
    c.AppendMiddleware(&FirstMiddleware{c})
    c.AppendMiddleware(&SecondMiddleware{c})

    return c
}
```

---

## Debugging Tips

### Print Current Time

```go
fmt.Printf("[Time %.2f ns] Event occurred\n",
    engine.CurrentTime() * 1e9)
```

### Debug Message Flow

```go
func (c *Comp) Tick() bool {
    msg := c.port.Retrieve()
    if msg != nil {
        fmt.Printf("[%s] Received: %T from %s\n",
            c.Name(),
            msg,
            msg.Meta().Src.Name())
    }
    return false
}
```

### Check Port Buffer Status

```go
// Custom debugging
type DebugPort struct {
    sim.Port
}

func (p *DebugPort) Send(msg sim.Msg) error {
    fmt.Printf("Sending %T to %s\n", msg, p.Name())
    return p.Port.Send(msg)
}
```

### Verify Connections

```go
func verifyConnections(ports []sim.Port) {
    for _, port := range ports {
        if port.Connection() == nil {
            fmt.Printf("WARNING: %s not connected!\n",
                port.Name())
        }
    }
}
```

### Simulation Stuck?

```go
// Add timeout to engine
engine := sim.NewSerialEngine()

// In your main loop
done := make(chan bool)
go func() {
    engine.Run()
    done <- true
}()

select {
case <-done:
    fmt.Println("Simulation completed")
case <-time.After(10 * time.Second):
    fmt.Println("Simulation timeout - stuck!")
    // Debug: print component states
}
```

### Check Component Activity

```go
type DebugComp struct {
    *sim.TickingComponent
    tickCount int
}

func (c *DebugComp) Tick() bool {
    c.tickCount++
    if c.tickCount % 1000 == 0 {
        fmt.Printf("[%s] Ticked %d times\n",
            c.Name(), c.tickCount)
    }
    return false
}
```

---

## Time Units

```go
// Frequencies
1 * sim.Hz
1 * sim.KHz
1 * sim.MHz
1 * sim.GHz

// Time
1 * sim.PS   // Picosecond
1 * sim.NS   // Nanosecond
1 * sim.US   // Microsecond
1 * sim.MS   // Millisecond
1 * sim.S    // Second

// Conversions
freq := 2 * sim.GHz
period := 1.0 / float64(freq)  // In seconds

time := 100 * sim.NS
cycles := int(time * float64(freq))
```

---

## Common Patterns Checklist

### Starting a New Component

- [ ] Define struct with `*sim.TickingComponent`
- [ ] Add ports as fields
- [ ] Implement `Tick() bool`
- [ ] Create Builder struct
- [ ] Implement `Build(name string)`
- [ ] Initialize TickingComponent in Build
- [ ] Create ports in Build
- [ ] Return pointer to component

### Connecting Components

- [ ] Create all components with builders
- [ ] Create ports for each component
- [ ] Create connection(s)
- [ ] PlugIn all ports to appropriate connections
- [ ] Verify connections are not nil

### Sending Messages

- [ ] Define message struct with `sim.MsgMeta`
- [ ] Set Src and Dst ports
- [ ] Set TrafficBytes if modeling bandwidth
- [ ] Handle `Send()` error (buffer full)
- [ ] Store pending message if send fails

### Testing

- [ ] Create test engine
- [ ] Build components
- [ ] Connect components
- [ ] Schedule initial events or send messages
- [ ] Run engine
- [ ] Assert expected state/outputs

---

## Performance Checklist

- [ ] Use `ParallelEngine` for large simulations
- [ ] Limit port buffer sizes (avoid unlimited)
- [ ] Minimize hook overhead (only hook what you need)
- [ ] Profile with `go tool pprof` if slow
- [ ] Use middleware to organize tick logic
- [ ] Batch operations where possible
- [ ] Return early from Tick() if no work

---

## File Locations Cheat Sheet

| What | Where |
|------|-------|
| Component interface | `sim/component.go` |
| Port interface | `sim/port.go` |
| Message interface | `sim/msg.go` |
| Engine | `sim/engine.go` |
| DirectConnection | `sim/directconnection/directconnection.go` |
| Hooks | `sim/hook.go` |
| Examples | `examples/` |
| Ideal memory | `mem/idealmemcontroller/` |
| Cache | `mem/cache/` |
| Tracing | `tracing/` |

---

## Quick Start Template

```go
package main

import (
    "fmt"
    "github.com/sarchlab/akita/v4/sim"
    "github.com/sarchlab/akita/v4/sim/directconnection"
)

type MyComp struct {
    *sim.TickingComponent
    port sim.Port
}

func (c *MyComp) Tick() bool {
    // Your logic here
    return false
}

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
    c.TickingComponent = sim.NewTickingComponent(
        name, b.engine, b.freq, c)
    c.port = sim.NewLimitNumMsgPort(c, 4, name+".Port")
    return c
}

func main() {
    engine := sim.NewSerialEngine()

    comp1 := MakeBuilder().WithEngine(engine).Build("Comp1")
    comp2 := MakeBuilder().WithEngine(engine).Build("Comp2")

    conn := directconnection.NewDirectConnection(
        "Link", engine, 1*sim.GHz)
    conn.PlugIn(comp1.port, 64)
    conn.PlugIn(comp2.port, 64)

    // Send initial message or schedule event

    engine.Run()

    fmt.Println("Done!")
}
```

---

## Common Errors and Solutions

| Error | Cause | Solution |
|-------|-------|----------|
| "Port buffer full" | Sending too fast | Check `Send()` error, retry later |
| "Simulation won't end" | Tick() always returns true | Return false when idle |
| "Nil pointer on port" | Port not initialized | Create port in Build() |
| "No connection" | Forgot to PlugIn | Call `conn.PlugIn()` for all ports |
| "Messages not arriving" | Wrong Src/Dst | Verify Src/Dst ports are correct |
| "Panic: nil engine" | Builder missing engine | Call `WithEngine()` on builder |

---

**This reference assumes AkitaSim v4. For v5 features, see migration_guide.md**

Quick links:
- Tutorial: `AKITASIM_TUTORIAL.md`
- Exercises: `HANDS_ON_EXERCISES.md`
- Examples: `examples/`
- Docs: https://akitasim.dev
