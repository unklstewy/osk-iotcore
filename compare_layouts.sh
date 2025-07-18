#!/bin/bash

# Visual regression test for keyboard layouts
# Compare generated layouts against reference images

echo "=== Visual Regression Test ==="
echo "Comparing generated layouts against reference images..."
echo

OUT_DIR="assets/reference/out"
DIFF_DIR="$OUT_DIR/diffs"
mkdir -p "$DIFF_DIR"

layouts=("style_one" "style_two" "style_three" "style_four")

for layout in "${layouts[@]}"; do
    echo "Testing layout: $layout"
    
    generated="$OUT_DIR/${layout}.png"
    reference="$OUT_DIR/onscreen-${layout//_/-}-ref.png"
    diff_output="$DIFF_DIR/${layout}_diff.png"
    
    if [[ -f "$generated" ]] && [[ -f "$reference" ]]; then
        # Get dimensions of both images
        gen_size=$(identify -format "%wx%h" "$generated")
        ref_size=$(identify -format "%wx%h" "$reference")
        
        echo "  Generated: $gen_size"
        echo "  Reference: $ref_size"
        
        # Resize images to compare (take the smaller dimensions)
        gen_width=$(echo $gen_size | cut -d'x' -f1)
        gen_height=$(echo $gen_size | cut -d'x' -f2)
        ref_width=$(echo $ref_size | cut -d'x' -f1)
        ref_height=$(echo $ref_size | cut -d'x' -f2)
        
        # Use the smaller dimensions for comparison
        compare_width=$((gen_width < ref_width ? gen_width : ref_width))
        compare_height=$((gen_height < ref_height ? gen_height : ref_height))
        
        echo "  Comparing at: ${compare_width}x${compare_height}"
        
        # Resize both images to the same size for comparison
        convert "$generated" -resize "${compare_width}x${compare_height}!" "$DIFF_DIR/${layout}_gen_resized.png"
        convert "$reference" -resize "${compare_width}x${compare_height}!" "$DIFF_DIR/${layout}_ref_resized.png"
        
        # Compare the images and get difference percentage
        # Using the mean error metric to calculate percentage difference
        diff_result=$(compare -metric AE "$DIFF_DIR/${layout}_gen_resized.png" "$DIFF_DIR/${layout}_ref_resized.png" "$diff_output" 2>&1)
        
        if [[ $? -eq 0 ]]; then
            echo "  ✓ Images are identical"
            deviation="0"
        else
            # Calculate percentage difference
            total_pixels=$((compare_width * compare_height))
            if [[ "$diff_result" =~ ^[0-9]+$ ]]; then
                deviation=$(echo "scale=2; $diff_result * 100 / $total_pixels" | bc -l)
                echo "  Different pixels: $diff_result out of $total_pixels"
                echo "  Deviation: ${deviation}%"
            else
                echo "  ⚠ Could not calculate difference: $diff_result"
                deviation="unknown"
            fi
        fi
        
        # Check if deviation is within acceptable range (< 5%)
        if [[ "$deviation" != "unknown" ]] && (( $(echo "$deviation < 5" | bc -l) )); then
            echo "  ✓ PASS - Deviation ${deviation}% is within acceptable range (<5%)"
        elif [[ "$deviation" == "0" ]]; then
            echo "  ✓ PASS - Images are identical"
        else
            echo "  ⚠ ATTENTION - Deviation ${deviation}% exceeds threshold (5%)"
            echo "    Difference visualization saved to: $diff_output"
        fi
    else
        echo "  ⚠ Missing files - Generated: $generated, Reference: $reference"
    fi
    
    echo
done

echo "=== Visual Regression Test Complete ==="
echo "Generated images are stored in: $OUT_DIR"
echo "Difference visualizations are in: $DIFF_DIR"
