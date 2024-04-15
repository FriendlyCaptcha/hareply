# hareply
`hareply` is a small tool that implements a responder to HAProxy's `agent-check` that replies with the contents of a given file.

## CLI usage

```shell
# Print help
hareply -h

# Print help for serve
hareply serve -h

# Print version information
hareply version

# Serve from port 8020 (default=8442)
hareply serve -f /some/path/to/agentstate -p 8020
```

## As a library
`hareply` can also be used as a library, see the godoc for the `hareply` package.

## Error handling

The response "`agentstate`" file is read on startup, if that fails the program will exit. The file is read again 
on any TCP connection, if that fails the last known file contents are used.
