# Roadmap

## Current Goals

- Complete Wayland protocol implementation.
- Full X11 support with input method integration.
- Theme and layout customization.
- Multi-language support.

## Future Milestones

- Configuration GUI.
- Performance optimizations.

## Regenerating Wayland Protocol Stubs

If the Wayland protocol definitions change, you can regenerate the protocol stubs by following these steps:

1. Ensure the `wayland-scanner` tool is installed on your system.
2. Run the Wayland files through the `wayland-scanner`:
   ```bash
   wayland-scanner client-header path/to/protocol.xml > path/to/generated_protocol.h
   wayland-scanner private-code path/to/protocol.xml > path/to/generated_protocol.c
   ```

This will update the generated code to the latest protocol specifications.

---
