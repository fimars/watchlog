# watchlogs

WIP

## tldr

```bash
go run ./cmd/watchlogs/
```

## Issue

1.  When fs-notify emit a file-write event, we open the file and read the last line, if it's the same as the last line we read before, we ignore it.

- Solution A: count the file size and compare it with the last file size we read before.
- Solution B: only support one file, we read the file from the beginning every time.
