# HTTP API with Context Support
Create a basic HTTP API (using net/http or Gin) with an endpoint that simulates a long-running task. Use context.WithTimeout to:

Cancel the operation if it exceeds a certain duration,

Detect client disconnection and stop work accordingly.
