# Debugging Guide

## Debug Logging

ProjectVem includes comprehensive debug logging to help diagnose keybinding issues, mode transitions, and explorer behavior.

### Running with Debug Output

```bash
# Run with debug logging to console
go run .

# Or use the convenience script
./run_debug.sh
```

### Log Categories

#### [KEY] - Key Press Events

Logs every key press with detailed information:

```
[KEY] Key="t" Modifiers=none Mode=NORMAL ExplorerVisible=false ExplorerFocused=false
```

Fields:
- **Key**: The key name (e.g., "t", "⏎", "Ctrl")
- **Modifiers**: Active modifiers (none, Ctrl, Shift, Alt, or combinations)
- **Mode**: Current editor mode (NORMAL, INSERT, VISUAL, DELETE, COMMAND, EXPLORER)
- **ExplorerVisible**: Whether file tree is shown
- **ExplorerFocused**: Whether file tree has focus

#### [MOD_MATCH] - Modifier Matching

Shows how modifier keys are being matched:

```
[MOD_MATCH] Required: Ctrl=true Shift=false Alt=false | Actual: Ctrl=true Shift=false Alt=false
[MOD_MATCH] Exact match -> true
```

This helps debug issues where:
- Platform doesn't report modifiers correctly
- Extra modifiers prevent matches
- Required modifiers are missing

#### [MATCH] - Keybinding Matches

Shows which keybinding was matched:

```
[MATCH] Global keybinding matched: Action=ActionToggleExplorer
[MATCH] Mode-specific keybinding matched: Mode=NORMAL Action=ActionMoveDown
```

Types of matches:
- **Global keybinding**: Works in any mode (Ctrl+T, Ctrl+H, Ctrl+L, Shift+Enter)
- **Mode-specific keybinding**: Only works in specific mode (h/j/k/l in NORMAL)
- **Special handler**: Complex logic (counts, goto sequences)

#### [ACTION] - Action Execution

Logs when an action is executed:

```
[ACTION] Executing action=ActionToggleExplorer mode=NORMAL
```

Shows which action is running and in what mode context.

#### [TOGGLE_EXPLORER] - Explorer Toggle Events

Detailed logging for explorer visibility toggling:

```
[TOGGLE_EXPLORER] Before: visible=false focused=false mode=NORMAL
[TOGGLE_EXPLORER] After: visible=true focused=false mode=NORMAL
```

Helps diagnose:
- Why explorer isn't opening/closing
- Whether focus state is correct
- Mode transitions related to explorer

#### [MODE_CHANGE] - Mode Transitions

Tracks all mode changes:

```
[MODE_CHANGE] Entering EXPLORER mode from NORMAL
[MODE_CHANGE] Now in mode=EXPLORER explorerFocused=true
[MODE_CHANGE] Exiting mode=EXPLORER
[MODE_CHANGE] Exited EXPLORER -> now in NORMAL
```

#### [SPECIAL] - Special Handler Matches

Logs when special handlers (counts, goto) are invoked:

```
[SPECIAL] Normal mode special handler matched
```

#### [NO_MATCH] - Unmatched Keys

Logs keys that don't match any binding:

```
[NO_MATCH] No keybinding matched for key="x" modifiers=none
```

## Common Debugging Scenarios

### Issue: Ctrl+T doesn't toggle explorer

**What to check in logs:**

1. Is the key being received?
   ```
   [KEY] Key="t" Modifiers=Ctrl Mode=NORMAL ...
   ```

2. Is it matching the global keybinding?
   ```
   [MATCH] Global keybinding matched: Action=ActionToggleExplorer
   ```

3. Is the action executing?
   ```
   [ACTION] Executing action=ActionToggleExplorer mode=NORMAL
   ```

4. Is the explorer state changing?
   ```
   [TOGGLE_EXPLORER] Before: visible=false ...
   [TOGGLE_EXPLORER] After: visible=true ...
   ```

**Common problems:**
- Modifier not detected: Check if `Modifiers=none` instead of `Modifiers=Ctrl`
- Wrong key received: Platform might send different key code
- Match failed: Check `[MOD_MATCH]` logs for exact vs required modifiers

### Issue: Shift+Enter doesn't toggle fullscreen

**What to check:**

1. Key press with Shift modifier:
   ```
   [KEY] Key="⏎" Modifiers=Shift Mode=INSERT ...
   ```

2. Modifier matching (should match exactly):
   ```
   [MOD_MATCH] Required: Ctrl=false Shift=true Alt=false | Actual: Ctrl=false Shift=true Alt=false
   [MOD_MATCH] Exact match -> true
   ```

3. Global keybinding match:
   ```
   [MATCH] Global keybinding matched: Action=ActionToggleFullscreen
   ```

**Common problems:**
- Extra modifiers: `[MOD_MATCH] Extra Ctrl present -> false`
- Mode-specific binding intercepting: Return without Shift matches INSERT mode binding first
- Platform sends different key for Shift+Enter

### Issue: Mode gets stuck or changes unexpectedly

**What to check:**

1. Mode transitions:
   ```
   [MODE_CHANGE] Entering EXPLORER mode from NORMAL
   [MODE_CHANGE] Now in mode=EXPLORER ...
   ```

2. Actions that change modes:
   - `ActionEnterInsert` → INSERT mode
   - `ActionEnterExplorer` → EXPLORER mode
   - `ActionExitMode` → Returns to NORMAL (usually)

3. Unexpected mode changes:
   - Check what action triggered it
   - Verify keybinding isn't firing unexpectedly

### Issue: Keys not working in certain modes

**What to check:**

1. Current mode when key is pressed:
   ```
   [KEY] ... Mode=EXPLORER ...
   ```

2. Which phase matched:
   - Global: Works in any mode
   - Mode-specific: Only that mode
   - Special: Custom logic

3. If no match:
   ```
   [NO_MATCH] No keybinding matched for key="h" modifiers=none
   ```
   
   Check if keybinding is defined for that mode in `keybindings.go`

## Debug Log Examples

### Example: Opening Explorer with Ctrl+T

```
[KEY] Key="t" Modifiers=Ctrl Mode=NORMAL ExplorerVisible=false ExplorerFocused=false
[MOD_MATCH] Required: Ctrl=true Shift=false Alt=false | Actual: Ctrl=true Shift=false Alt=false
[MOD_MATCH] Exact match -> true
[MATCH] Global keybinding matched: Action=ActionToggleExplorer
[ACTION] Executing action=ActionToggleExplorer mode=NORMAL
[TOGGLE_EXPLORER] Before: visible=false focused=false mode=NORMAL
[TOGGLE_EXPLORER] After: visible=true focused=false mode=NORMAL
```

### Example: Entering INSERT mode

```
[KEY] Key="i" Modifiers=none Mode=NORMAL ExplorerVisible=false ExplorerFocused=false
[MOD_MATCH] Required none, got 0 -> true (implied)
[MATCH] Mode-specific keybinding matched: Mode=NORMAL Action=ActionEnterInsert
[ACTION] Executing action=ActionEnterInsert mode=NORMAL
```

### Example: Failed modifier match (Shift+Enter with extra Ctrl)

```
[KEY] Key="⏎" Modifiers=Ctrl+Shift Mode=NORMAL ...
[MOD_MATCH] Required: Ctrl=false Shift=true Alt=false | Actual: Ctrl=true Shift=true Alt=false
[MOD_MATCH] Extra Ctrl present -> false
[MOD_MATCH] Required: Ctrl=false Shift=false Alt=false | Actual: Ctrl=true Shift=true Alt=false
[MOD_MATCH] Extra Ctrl present -> false
[NO_MATCH] No keybinding matched for key="⏎" modifiers=Ctrl+Shift
```

## Adding Custom Debug Logging

To add more debug logging to track specific behavior:

```go
// In any function where you want logging
log.Printf("[YOUR_TAG] Your message: var=%v", someVariable)
```

Common patterns:
- Use UPPERCASE tags in square brackets: `[MY_TAG]`
- Include relevant state: mode, flags, variables
- Log before and after state changes
- Use descriptive messages

## Performance Considerations

Debug logging adds overhead. For production builds, consider:

1. Conditional compilation with build tags
2. Log level filtering (only errors/warnings in production)
3. Disabling verbose logs for performance-critical paths

Currently, all debug logs are always active. This is intentional during development to help diagnose issues.
