# Mock Watch Server for Testing

This is a standalone test server that simulates the fulfillment service for testing the watch functionality of the CLI.

## Starting the Server

1. Build the server:
```bash
go build -o test-server ./cmd/test-server
```

2. Run the server with the default scenario:
```bash
./test-server
```

3. Or run with a custom scenario file:
```bash
./test-server -scenario path/to/your/scenario.yaml
```

The server will start listening on `127.0.0.1:8080` and display instructions for connecting with the CLI.

## Using the CLI with the Mock Server

1. First, set the fulfillment URL to point to the mock server:
```bash
export FULFILLMENT_URL=http://127.0.0.1:8080
```

2. Login with the fulfillment CLI (mock server requires no authentication):
```bash
./fulfillment-cli login --address http://127.0.0.1:8080
```

3. Watch for cluster events:
```bash
# Watch all clusters
./fulfillment-cli get clusters --watch

# Watch a specific cluster
./fulfillment-cli get clusters my-cluster --watch
```

4. Stop watching by pressing `Ctrl+C`

## Event Scenarios

The mock server uses event scenarios defined in YAML files to simulate cluster lifecycle events.

### Default Scenario

The default scenario is loaded from `internal/testing/testdata/default_scenario.yaml` and simulates a complete cluster creation lifecycle with 4 events:

1. Cluster created (PROGRESSING state)
2. Cluster updated (installing control plane)
3. Cluster updated (READY state)
4. Second cluster created

### Creating Custom Scenarios

You can create custom event scenarios by writing YAML files with the following structure:

```yaml
name: my-custom-scenario
description: Description of what this scenario tests
events:
  - id: event-1
    type: EVENT_TYPE_OBJECT_CREATED
    delaySeconds: 0
    cluster:
      id: my-cluster-id
      name: my-cluster-name
      state: CLUSTER_STATE_PROGRESSING
      conditions:
        - type: CLUSTER_CONDITION_TYPE_READY
          status: CONDITION_STATUS_FALSE
          message: Cluster is being created
  - id: event-2
    type: EVENT_TYPE_OBJECT_UPDATED
    delaySeconds: 3
    cluster:
      id: my-cluster-id
      name: my-cluster-name
      state: CLUSTER_STATE_READY
      conditions:
        - type: CLUSTER_CONDITION_TYPE_READY
          status: CONDITION_STATUS_TRUE
          message: Cluster is ready
```

### Event Types

Valid event types:
- `EVENT_TYPE_OBJECT_CREATED`
- `EVENT_TYPE_OBJECT_UPDATED`
- `EVENT_TYPE_OBJECT_DELETED`

### Cluster States

Valid cluster states:
- `CLUSTER_STATE_PROGRESSING`
- `CLUSTER_STATE_READY`
- `CLUSTER_STATE_ERROR`

### Condition Types

Valid condition types:
- `CLUSTER_CONDITION_TYPE_READY`

### Condition Status

Valid condition statuses:
- `CONDITION_STATUS_TRUE`
- `CONDITION_STATUS_FALSE`
- `CONDITION_STATUS_UNKNOWN`

### Delay Seconds

The `delaySeconds` field specifies how many seconds to wait before sending the event. This simulates the timing of real cluster provisioning.

## Server Behavior

- The server sends events from the scenario in sequence
- Events are filtered based on the watch request filter
- The server waits for the client to disconnect before cleaning up
- Each connection receives the full scenario from the beginning
- The server logs each event it sends for debugging

## Troubleshooting

### Server won't start
- Check if port 8080 is already in use
- Try changing the `serverPort` constant in `main.go`

### No events received
- Verify you're logged in with `./fulfillment-cli login`
- Check the FULFILLMENT_URL environment variable
- Look at server logs to see if events are being sent

### Wrong events received
- Check the filter being applied (server logs show the filter)
- Verify your scenario YAML syntax is correct
- Ensure event IDs and names match your filter criteria
