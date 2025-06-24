# Panic Recovery in Goroutines
Create a `SafeGo(func())` function that:

Runs the given function in a new goroutine,

Recovers from any panic,

Logs the error without crashing the app.

Make sure other goroutines continue running after a panic.
