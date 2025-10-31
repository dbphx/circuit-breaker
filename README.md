# 🧠 Circuit Breaker for Go

This package provides a lightweight, flexible **Circuit Breaker** implementation in Go.  
It helps prevent cascading failures by stopping calls to failing dependencies and allowing recovery after a cooldown period.

---

## 🚀 Features

- Supports **failure and success thresholds**
- Configurable **open timeout duration**
- **OnStateChange** callback hook for monitoring transitions
- Thread-safe with `sync.Mutex`
- Simple and type-safe generic API: `Execute(func() (T, error)) (T, error)`
- Implements standard states:
  - `Closed`: normal operation
  - `Open`: rejecting all requests
  - `Half-Open`: testing recovery

---

## ⚙️ Installation

```bash
go get git.infiniband.vn/cloud/insky/eco/utils/circuit-breaker
