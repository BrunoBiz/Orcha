# Request logs and remote command diagnostics

Orcha writes the same concise, structured text logs to its console and the
existing `Log/YYYY-MM-DD-orcha.log` file. Each HTTP request produces **one**
completion record containing the request ID, method, path, HTTP status and
duration. HTTP 5xx failures are ERROR, 4xx are WARN, and successful requests are
INFO. Request/response headers, bodies, query strings and routine execution
steps are not dumped at the default level.

Every response includes `X-Request-ID`. Server-operation failures also include
this ID in the HTTP error message. Orcha Web preserves the header and displays
the ID with the error, so it identifies the exact failed request even when
automatic status polling is running concurrently.

## Failed manager commands

For `start`, `stop`, `restart` and `details`, the same request record includes:

- `vmid` and `action`
- `command` (the `pct exec ... ./Narwhal ...` command)
- `error` and `exit_code` when available
- `output`: captured stdout and stderr, with ANSI colors removed
- `error_code`: `SSH_COMMAND_FAILED` or `UNRECOGNIZED_MANAGER_OUTPUT`
- `output_truncated`: whether a long output was shortened

Orcha captures output on both non-zero exits and exit-zero responses it cannot
interpret. This distinguishes a missing user or executable from a manager
output-format mismatch. A remote exit code is not available for every SSH
failure (for example, failure to open a session).

Command output is bounded to an 8 KiB excerpt plus the truncation marker, keeping
the beginning and end. Raw command output is retained in server logs rather
than copied into browser responses. The response contains the operation, VMID,
error category and request ID. If truncation occurs, a deeper investigation may
still need the manager's own log.

An illustrative failure record (time/duration omitted here) looks like:

```text
level=ERROR msg="request completed" request_id=example method=POST path=/containers/server/101/details status=500 vmid=101 action=details exit_code=1 error_code=SSH_COMMAND_FAILED output="su: user gameserver does not exist"
```

Read the running Orcha process's console to see errors live, or follow the log
file it opened on startup:

```bash
tail -F Log/YYYY-MM-DD-orcha.log
```

Replace the date with the filename created when that Orcha process started.
No manual re-execution of the manager command is needed to capture its error.

## Log levels

Set the **process environment variable** `LOG_LEVEL` to `debug`, `info`, `warn`
or `error` before launching Orcha. Default: `info`. An invalid value prevents
startup with a configuration error. This setting is read by logger startup;
putting it only in `pm.env` does not set the process environment variable.

```bash
LOG_LEVEL=debug ./Orcha
```

Debug mode adds source locations in the file and enables existing internal
diagnostic steps. It still does not enable HTTP header/body dumps. Use `warn`
to keep warnings/errors while suppressing successful polling summaries.

The existing log-file retention behavior is unchanged: cleanup runs on startup
and removes files older than 15 days. Logging changes require rebuilding and
restarting the Orcha binary; frontend changes alone cannot apply them.
