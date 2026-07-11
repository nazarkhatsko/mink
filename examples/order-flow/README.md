# 📦 order-flow

A local Go HTTP server implementing an order lifecycle (`/orders`, add items, checkout) with `X-Api-Key` auth. Demonstrates active use of `mutate_on` for state accumulation across many actions.

🚀 **Start the server**
```bash
go run examples/order-flow/main.go
```

▶️ **Run the flow**
```bash
mink run examples/order-flow/mink.yaml
```

📋 **Flow**
| Flow | What it does |
|---|---|
| `order_lifecycle` | Creates an order, adds three items while accumulating their IDs/names into `state`, validates the accumulated count, fetches the order to capture the server-computed total into `state`, checks out, and validates the final state against the server response |

⬅️ [Back to examples](../README.md)
