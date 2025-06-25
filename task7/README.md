# Select over Multiple Channels

Write a function that listens on three channels:

* jobs
* shutdown
* heartbeat

Use a select block to:

* Handle job messages
* Timeouts with `time.After`
* Clean shutdowns when signaled
