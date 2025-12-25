# Getting Started with AkitaSim

Welcome to AkitaSim! This guide will help you get started with learning and using the AkitaSim simulation framework.

## 📚 Learning Resources

We've prepared comprehensive learning materials to help you master AkitaSim:

### 1. **[Comprehensive Tutorial](./AKITASIM_TUTORIAL.md)** 📖
   - **Start here!** Complete guide covering all core concepts
   - Progressive learning path from basics to advanced topics
   - Detailed explanations of components, ports, messages, connections, hooks, and domains
   - Real-world examples and code samples
   - Best practices and design patterns
   - **Time**: 3-4 weeks for full mastery

### 2. **[Hands-On Exercises](./HANDS_ON_EXERCISES.md)** 💻
   - 7 practical exercises to build real components
   - Progressive difficulty: from simple events to complete systems
   - Solutions provided for self-checking
   - Build: counters, FIFOs, caches, memory hierarchies
   - **Time**: 1-2 weeks of practice

### 3. **[Quick Reference Guide](./QUICK_REFERENCE.md)** ⚡
   - Condensed API reference for quick lookup
   - Code snippets and patterns
   - Common imports and usage examples
   - Debugging tips and troubleshooting
   - **Use**: Keep open while coding!

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/sarchlab/akita
cd akita

# Ensure Go is installed (1.18+)
go version

# Run your first example
cd examples/01_print_event
go run main.go
```

### Your First Simulation (5 minutes)

Create `hello.go`:

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

type PrintHandler struct{}

func (h PrintHandler) Handle(e sim.Event) error {
    evt := e.(*PrintEvent)
    fmt.Printf("[Time %.2f] %s\n", evt.Time(), evt.message)
    return nil
}

func main() {
    engine := sim.NewSerialEngine()
    handler := PrintHandler{}

    event := &PrintEvent{
        EventBase: sim.NewEventBase(5.0, handler),
        message:   "Hello, AkitaSim!",
    }
    engine.Schedule(event)

    engine.Run()
    fmt.Println("Simulation complete!")
}
```

Run it:

```bash
go run hello.go
```

Output:
```
[Time 5.00] Hello, AkitaSim!
Simulation complete!
```

Congratulations! You just ran your first AkitaSim simulation! 🎉

## 📖 Recommended Learning Path

### Week 1: Fundamentals
- [ ] Read [Tutorial](./AKITASIM_TUTORIAL.md) sections 1-6
- [ ] Complete [Exercises](./HANDS_ON_EXERCISES.md) 1-2
- [ ] Run all examples in `examples/` directory
- [ ] Understand: Events, Components, Ports, Messages

### Week 2: Building Components
- [ ] Read [Tutorial](./AKITASIM_TUTORIAL.md) sections 7-9
- [ ] Complete [Exercises](./HANDS_ON_EXERCISES.md) 3-5
- [ ] Study `mem/idealmemcontroller/`
- [ ] Build your own simple component
- [ ] Understand: Connections, Ticking, Builders

### Week 3: Advanced Topics
- [ ] Read [Tutorial](./AKITASIM_TUTORIAL.md) sections 10-11
- [ ] Complete [Exercises](./HANDS_ON_EXERCISES.md) 6-7
- [ ] Study `mem/cache/` implementation
- [ ] Add tracing to your components
- [ ] Build a multi-component system

### Week 4: Real-World Application
- [ ] Study complete systems in `examples/`
- [ ] Explore Network-on-Chip in `noc/`
- [ ] Design your own architecture experiment
- [ ] Use [Quick Reference](./QUICK_REFERENCE.md) as needed
- [ ] Read research papers using AkitaSim

## 🎯 Learning by Example

AkitaSim includes several examples of increasing complexity:

| Example | Concepts | Location |
|---------|----------|----------|
| **Print Event** | Basic events, engine | `examples/01_print_event/` |
| **Cell Split** | Event chains, scheduling | `examples/02_cell_split/` |
| **Ping (Event)** | Components, messages, ports | `examples/ping/` |
| **Ping (Ticking)** | Ticking components, middlewares | `examples/tickingping/` |
| **Ideal Memory** | Real component, builder pattern | `mem/idealmemcontroller/` |
| **Cache** | Complex state, hit/miss logic | `mem/cache/` |
| **Network** | NoC, switches, routing | `noc/acceptance/` |

**Recommended order**: Follow the table from top to bottom.

## 🔑 Key Concepts Overview

```
┌─────────────────────────────────────────────────┐
│              SIMULATION                          │
│                                                  │
│  ┌────────────────────────────────────────┐    │
│  │         ENGINE                          │    │
│  │  (manages time & events)                │    │
│  └────────────────────────────────────────┘    │
│                                                  │
│  ┌──────────────┐        ┌──────────────┐      │
│  │  COMPONENT   │◄──────►│  COMPONENT   │      │
│  │   (CPU,      │  CONN  │  (Memory,    │      │
│  │   Cache...)  │        │   Cache...)  │      │
│  └──────────────┘        └──────────────┘      │
│                                                  │
│  Communication via MESSAGES through PORTS       │
│  Observation via HOOKS                          │
│  Organization via DOMAINS                       │
└─────────────────────────────────────────────────┘
```

### Core Abstractions

- **Engine**: Drives simulation time
- **Component**: Hardware element (CPU, cache, memory)
- **Port**: Communication endpoint
- **Message**: Data exchanged between components
- **Connection**: Links ports together
- **Event**: Something that happens at a specific time
- **Hook**: Callback for monitoring/tracing
- **Domain**: Logical grouping of components

## 🛠️ Development Workflow

1. **Design** your component's interface (ports, messages)
2. **Implement** the component logic (Tick or Handle)
3. **Create** a Builder for easy construction
4. **Test** the component in isolation
5. **Connect** to other components
6. **Trace** with hooks to verify behavior
7. **Analyze** simulation results

## 📊 What Can You Build?

AkitaSim is used for:

- ✅ CPU microarchitecture simulation
- ✅ GPU architecture research (see MGPUSim)
- ✅ Memory hierarchy experiments
- ✅ Network-on-Chip design
- ✅ Cache coherence protocols
- ✅ Interconnect modeling
- ✅ Performance analysis
- ✅ Power modeling
- ✅ Hardware/software co-design

## 🔗 Additional Resources

### Official Documentation
- **Website**: https://akitasim.dev
- **API Docs**: https://pkg.go.dev/github.com/sarchlab/akita/v4
- **GitHub**: https://github.com/sarchlab/akita

### Related Projects
- **MGPUSim**: Multi-GPU simulator built with AkitaSim
- **Daisen**: Web-based trace visualization tool

### Community
- **Issues**: Report bugs or ask questions on GitHub Issues
- **Discussions**: Join conversations on GitHub Discussions

## 🐛 Troubleshooting

### Common Issues

**Simulation won't end**
- Check that all `Tick()` methods return `false` when idle
- Verify no infinite event loops

**Messages not arriving**
- Verify ports are connected: `port.Connection() != nil`
- Check message Src/Dst are set correctly
- Ensure both ports plugged into same connection

**Port buffer full errors**
- Increase buffer size in `NewLimitNumMsgPort()`
- Implement backpressure in sender
- Store pending messages and retry

**Need help?**
- Check [Quick Reference](./QUICK_REFERENCE.md) for common patterns
- Review [Tutorial](./AKITASIM_TUTORIAL.md) for detailed explanations
- Look at working examples in `examples/`
- Search GitHub Issues

## 🎓 Advanced Topics

Once you're comfortable with basics:

- **v5 Migration**: Read `migration_guide.md` for v5 features
- **Middleware Pattern**: Organize complex component logic
- **Parallel Simulation**: Use `ParallelEngine` for speed
- **Tracing & Analysis**: Deep dive into `tracing/` package
- **Network Modeling**: Explore `noc/` for realistic networks
- **GPU Simulation**: Study MGPUSim for large-scale examples

## 🚦 You're Ready When...

You understand AkitaSim when you can:

- ✅ Explain the difference between event-driven and ticking components
- ✅ Create a custom component with ports
- ✅ Send messages between components
- ✅ Use the Builder pattern
- ✅ Add hooks for tracing
- ✅ Connect multiple components
- ✅ Debug common simulation issues

## 📝 Next Actions

1. **Read**: [Comprehensive Tutorial](./AKITASIM_TUTORIAL.md)
2. **Practice**: [Hands-On Exercises](./HANDS_ON_EXERCISES.md)
3. **Code**: Build something simple
4. **Reference**: Use [Quick Reference](./QUICK_REFERENCE.md)
5. **Explore**: Study real components in `mem/`
6. **Create**: Design your own architecture experiment

---

**Ready to start? Open [AKITASIM_TUTORIAL.md](./AKITASIM_TUTORIAL.md) and begin your journey!** 🚀

Good luck, and happy simulating!
