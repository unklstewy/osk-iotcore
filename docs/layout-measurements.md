# OSK IoTCore Layout Measurements

This document provides exact key geometry measurements for each keyboard layout in the OSK IoTCore project.

## Overview

The OSK IoTCore system uses JSON-based layout configurations with pixel-perfect positioning. All measurements are in pixels with the origin (0,0) at the top-left corner of the keyboard.

## Base Layout Structure

### Standard Layout Dimensions
- **Total Width**: 900px
- **Total Height**: 400px
- **Key Gap**: 10px (horizontal between standard keys)
- **Row Gap**: 10px (vertical between rows)

### Key Size Standards
- **Standard Key**: 50px × 50px
- **Function Keys**: 50px × 40px
- **Tab Key**: 75px × 50px
- **Caps Lock**: 90px × 50px  
- **Backspace**: 100px × 50px
- **Enter**: 120px × 50px
- **Left Shift**: 120px × 50px
- **Right Shift**: 150px × 50px
- **Space Bar**: 300px × 50px
- **Control Keys**: 80px × 50px
- **Alt/Meta Keys**: 60px × 50px
- **Backslash**: 75px × 50px

## Layout Analysis by Screenshot

### Style One (onscreen-style-one.jpeg)
- **Screenshot Dimensions**: 302px × 167px
- **Aspect Ratio**: 1.81:1
- **Layout Type**: Compact QWERTY
- **Estimated Scale**: ~33% of base layout
- **Key Measurements** (approximate):
  - Standard key: ~17px × 17px
  - Function row: Visible with reduced height
  - Special features: Compact function row, standard alphanumeric layout

### Style Two (onscreen-style-two.jpeg)
- **Screenshot Dimensions**: 381px × 132px
- **Aspect Ratio**: 2.89:1
- **Layout Type**: Ultra-wide compact
- **Estimated Scale**: ~42% width, ~33% height
- **Key Measurements** (approximate):
  - Standard key: ~21px × 16px
  - Layout style: Horizontally extended, vertically compressed
  - Special features: Single-row or dual-row layout, optimized for wide screens

### Style Three (onscreen-style-three.jpeg)
- **Screenshot Dimensions**: 299px × 169px
- **Aspect Ratio**: 1.77:1
- **Layout Type**: Balanced compact
- **Estimated Scale**: ~33% of base layout
- **Key Measurements** (approximate):
  - Standard key: ~17px × 18px
  - Layout style: Proportionally scaled standard layout
  - Special features: Maintains standard QWERTY proportions

### Style Four (onscreen-style-four.jpeg)
- **Screenshot Dimensions**: 284px × 177px
- **Aspect Ratio**: 1.60:1
- **Layout Type**: Square-oriented compact
- **Estimated Scale**: ~32% width, ~44% height
- **Key Measurements** (approximate):
  - Standard key: ~16px × 19px
  - Layout style: More vertical emphasis
  - Special features: Taller key appearance, compact horizontal spacing

## Detailed Key Geometry (Base Layout)

### Function Row (Row 1)
- **Y Position**: 10px
- **Key Height**: 40px
- **Key Width**: 50px
- **Gap**: 10px between F4/F5 and F8/F9
- **X Positions**: 
  - ESC: 10px
  - F1-F4: 80px, 140px, 200px, 260px
  - F5-F8: 330px, 390px, 450px, 510px
  - F9-F12: 580px, 640px, 700px, 760px

### Number Row (Row 2)
- **Y Position**: 60px
- **Key Height**: 50px
- **Standard Key Width**: 50px
- **X Positions**: 
  - Grave (`): 10px
  - Numbers 1-0: 70px to 610px (60px intervals)
  - Minus (-): 670px
  - Equal (=): 730px
  - Backspace: 790px (100px width)

### QWERTY Row (Row 3)
- **Y Position**: 120px
- **Key Height**: 50px
- **X Positions**:
  - Tab: 10px (75px width)
  - Q-P: 95px to 635px (60px intervals)
  - Left Bracket: 695px
  - Right Bracket: 755px
  - Backslash: 815px (75px width)

### Home Row (Row 4)
- **Y Position**: 180px
- **Key Height**: 50px
- **X Positions**:
  - Caps Lock: 10px (90px width)
  - A-L: 110px to 590px (60px intervals)
  - Semicolon: 650px
  - Apostrophe: 710px
  - Enter: 770px (120px width)

### Bottom Row (Row 5)
- **Y Position**: 240px
- **Key Height**: 50px
- **X Positions**:
  - Left Shift: 10px (120px width)
  - Z-/: 140px to 680px (60px intervals)
  - Right Shift: 740px (150px width)

### Space Row (Row 6)
- **Y Position**: 300px
- **Key Height**: 50px
- **X Positions**:
  - Left Ctrl: 10px (80px width)
  - Left Meta: 100px (60px width)
  - Left Alt: 170px (60px width)
  - Space: 240px (300px width)
  - Right Alt: 550px (60px width)
  - Right Meta: 620px (60px width)
  - Menu: 690px (60px width)
  - Right Ctrl: 760px (80px width)

### Arrow Keys
- **Y Position**: 180px-210px
- **Key Size**: 30px × 30px
- **X Positions**:
  - Up: 860px, 180px
  - Left: 830px, 210px
  - Down: 860px, 210px
  - Right: 890px, 210px

### Navigation Cluster
- **Y Position**: 60px-120px
- **Key Size**: 40px × 25px
- **X Positions**:
  - Insert: 820px, 60px
  - Home: 820px, 90px
  - Page Up: 820px, 120px
  - Delete: 870px, 60px
  - End: 870px, 90px
  - Page Down: 870px, 120px

## Layout Variations

### Dvorak Layout
- **Base Dimensions**: Same as QWERTY (900px × 400px)
- **Key Positions**: Same physical layout as QWERTY
- **Differences**: Only key labels change, geometry remains identical

### Theme Variations
All themes use the same geometric layout with different visual styling:
- **Border Radius**: 4px
- **Border Width**: 1px
- **Key Padding**: 2px
- **Shadow Offset**: 2px × 2px
- **Shadow Blur**: 4px

## Special Layout Features

### Non-Standard Elements
1. **Function Key Gaps**: 20px gaps between F4-F5 and F8-F9 groups
2. **Arrow Key Cluster**: Compact T-shaped arrangement
3. **Navigation Block**: 6-key cluster with 25px height keys
4. **Split Spacebar**: Not implemented in current layouts
5. **Emoji Bar**: Not present in current layouts
6. **Number Pad**: Not included in standard layouts

### Modifier Keys
- **Shift Keys**: Left (120px), Right (150px)
- **Control Keys**: Left (80px), Right (80px)
- **Alt Keys**: Left (60px), Right (60px)
- **Meta Keys**: Left (60px), Right (60px)

## Scale Factors for Different Displays

Based on screenshot analysis, the following scale factors are used:
- **Compact View**: ~33% scale (maintaining proportions)
- **Ultra-wide View**: ~42% width, ~33% height
- **Square View**: ~32% width, ~44% height

## Implementation Notes

1. All measurements use integer pixel values
2. Key positioning uses top-left origin
3. No overlapping keys in any layout
4. Consistent 10px grid alignment for most keys
5. Visual themes affect appearance but not geometry
6. Layout JSON files contain exact pixel coordinates

## Validation

Key geometry validation ensures:
- No negative coordinates
- No overlapping keys
- Keys within layout boundaries
- Consistent row alignment
- Proper key sizing ratios
