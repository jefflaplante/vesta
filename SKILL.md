---
name: vestaboard
description: Control Vestaboard displays - send messages, read state, create dashboards
---

# Vestaboard Skill

Control Vestaboard Note and Flagship displays via the `vesta` CLI.

## Commands

```sh
vesta send "<message>"           # Send a message
vesta send -c "<message>"        # Send centered
vesta send --dry-run "<message>" # Preview without sending
vesta read                       # Read current board state
vesta config show                # Show current config
echo "msg" | vesta send -        # Read from stdin (scripting)
```

## Board Dimensions

| Device   | Rows | Columns | Total chars |
|----------|------|---------|-------------|
| Note     | 3    | 15      | 45          |
| Flagship | 6    | 22      | 132         |

Default is Note. Use `--device flagship` to override or `vesta config set device flagship` to change default.

## Escape Codes

### Colors (all devices)
`{red}` `{orange}` `{yellow}` `{green}` `{blue}` `{violet}` `{white}` `{black}` `{filled}`

### Symbols
- `{deg}` - degree symbol (Flagship only)
- `<3` or `{<3}` - heart (Note only)
- `{0}`-`{71}` - raw character codes

## Formatting

- Text auto-wraps at word boundaries
- Use `\n` for explicit line breaks
- Use `-c` flag for centered text
- Unsupported characters are skipped with warnings

## Examples

### Simple message
```sh
vesta send "Hello World"
```

### Centered with colors
```sh
vesta send -c "{red}{orange}{yellow}{green}{blue}{violet}"
```

### Multi-line
```sh
vesta send "Line 1\nLine 2\nLine 3"
```

### Weather dashboard (Flagship)
```sh
vesta send --device flagship "{blue}WEATHER{white}        72{deg}F\nSunny\n\n{yellow}HIGH 78{deg}  LOW 65{deg}\nHumidity 45%\nWind 8 mph NW"
```

### Status indicator
```sh
vesta send "{green} All systems operational"
```

### Preview before sending
```sh
vesta send --dry-run "Test message"
```

### Scripting with stdin
```sh
# Pipe output from another command
date +"%H:%M" | vesta send -

# Read from file
cat message.txt | vesta send -c -
```

## Best Practices

1. **Always preview first** - Use `--dry-run` before sending to verify formatting
2. **Mind the dimensions** - Note has only 45 chars; keep messages brief
3. **Use colors sparingly** - One or two accent colors work better than rainbows
4. **Test escape codes** - Not all codes work on all devices
5. **Handle errors** - Use `-v` flag for verbose error details if sends fail
