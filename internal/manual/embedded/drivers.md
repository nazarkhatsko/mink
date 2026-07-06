# Drivers

Drivers are the building blocks of mink flows. Each driver encapsulates interaction with an external service or utility.

## Built-in drivers

| Driver | Description |
|---|---|
| http | HTTP requests |
| generate | Fake data generation |
| validate | JSON Schema validation |
| sleep | Execution delay |
| shell | Shell command execution |
| python | Python code/script execution |
| claude | Claude API messages |
| firestore | Fetch a Firestore document by path |

Run `mink manual drivers <name>` for details on a specific driver.
