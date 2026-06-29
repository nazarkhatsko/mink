# Drivers

Drivers are the building blocks of mink flows. Each driver encapsulates interaction with an external service or utility.

## Built-in drivers

| Driver | Description |
|---|---|
| [http](http.md) | HTTP requests |
| [generate](generate.md) | Fake data generation |
| [validate](validate.md) | JSON Schema validation |
| [sleep](sleep.md) | Execution delay |

## Custom drivers

You can implement your own driver using the public SDK in `pkg/driver`. See [Contributing](../../CONTRIBUTING.md).
