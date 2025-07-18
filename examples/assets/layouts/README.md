# Keyboard Layout Files

This directory contains JSON files defining different keyboard layouts for the on-screen keyboard.

## File Structure

Each layout file follows this structure:
- `name`: Unique identifier for the layout
- `description`: Human-readable description of the layout
- `width` & `height`: Overall dimensions of the keyboard
- `keys`: Array of key objects with the following properties:
  - `id`: Unique identifier for the key
  - `label`: Display text on the key
  - `code`: Linux evdev input code
  - `x`, `y`: Position coordinates
  - `width`, `height`: Key dimensions
  - `modifier`: Boolean flag for modifier keys (optional)

## Custom Negative Placeholder Codes

For keys that don't have obvious Linux evdev codes, the following custom negative codes are used:

| Code | Description | Usage |
|------|-------------|--------|
| -1   | Numbers toggle | Switches to numeric/symbol input mode |
| -2   | Symbols toggle | Switches to symbol input mode |
| -3   | Microphone | Voice input activation |
| -4   | Language switch | Changes input language |
| -5   | Emoji selector | Opens emoji picker |
| -6   | Swipe/Gesture | Gesture input mode |

## Layout Files

- `style_one.json`: Standard mobile QWERTY layout
- `style_two.json`: Mobile QWERTY with number row and additional functions
- `style_three.json`: Compact layout with emoji support
- `style_four.json`: Minimal design with gesture support

## Notes

- The Go code safely ignores unknown negative codes
- Modifier keys should have the `modifier: true` flag
- Coordinates are relative to the keyboard's top-left corner
- All dimensions are in pixels
