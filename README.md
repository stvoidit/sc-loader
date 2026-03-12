# sc-loader

Lightweight Go utility to record streaming video (download playlists / TS segments), write to a temporary file and perform a final remux using ffmpeg.

## Installation

Clone repository:
```
git clone https://github.com/stvoidit/sc-loader.git
cd sc-loader
```

Build using provided Makefile:
```
make build
# produces executable sc-loader in repo root
```

Install to ~/.local/bin:
```
make install
```

Uninstall:
```
make uninstall
```

Or build manually:
```
go build -C src/ -ldflags "-s -w" -o sc-loader
```

Make sure ffmpeg is installed and available in PATH:
```
ffmpeg -version
```

## Usage

Run with a username or a user URL:
```
./sc-loader <username_or_user_url>
```

Examples:
```
./sc-loader some_user
./sc-loader https://example.com/user/some_user
```

If a URL is provided (it starts with `http`), the program uses the basename of the path as the username.

Behavior:
- Loads configuration (see below).
- Creates a temporary recording file named `<username>_YYYYMMDD-HHMMSS_tmp.mp4` in the configured folder.
- Polls the playlist and downloads new segments into the temporary file.
- On finish, runs `ffmpeg` to remux the temporary file into a final file without the `_tmp` suffix using stream copy and `-movflags +faststart`.
- Supports graceful termination (Ctrl+C).

## Configuration

sc-loader looks for `config.json` in the following order:
1. $XDG_CONFIG_HOME/sc-loader/config.json (usually ~/.config/sc-loader/config.json)
2. `./config.json` (current working directory)
3. `./config.json` (literal path)

Format: JSON. Required/used fields:

- `debug` (boolean) — enable debug logging (slog.LevelDebug)
- `host` (string) — host used to generate user URL
- `folder` (string) — directory where recordings are saved (will be created)
- `keys` (object) — map of DRM keys (string -> string)

Example `config.json`:
```json
{
  "debug": false,
  "host": "example.com",
  "folder": "~/videos",
  "keys": {
    "drm-key-1": "string",
    "drm-key-2": "string"
  }
}
```

Notes:
- `folder` is converted to an absolute path and created with 0700 permissions if necessary.
- `keys` are cached in memory as SHA-256 digests and not serialized back to disk.
- If config is missing or unreadable, the program will fail at startup.

## Output files

Temporary files use the pattern:
```
<username>_YYYYMMDD-HHMMSS_tmp.mp4
```
After successful processing, `ffmpeg` creates a final file without `_tmp` (stream copy + `-movflags +faststart`). The temporary file is removed.

## Requirements

- Go (see go.mod for required Go version)
- ffmpeg in PATH
- Write permission for the `folder` configured in `config.json`

## Notes & Recommendations

- The program exits if `ffmpeg` is not found.
- Ctrl+C is handled to allow graceful shutdown and final remux.
- Tune poll interval and retry behavior in code/config if needed.
