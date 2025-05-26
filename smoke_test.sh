#!/bin/bash
# Manual testing script for Gable

echo "=== Gable Manual Test Suite ==="
echo

echo "1. Testing help flag..."
./gable -h
echo

echo "2. Testing with invalid arguments..."
./gable 2>&1 | head -5
echo

echo "3. Testing with invalid repo format..."
./gable invalid another-invalid 2>&1
echo

echo "4. Testing debug output..."
./gable -d golang/go golang/go 2>&1 | grep DEBUG | head -5
echo

echo "5. Testing GitHub CLI check..."
PATH="" ./gable foo/bar foo/baz 2>&1
echo

echo "=== Tests complete ==="
echo "Note: Full interactive testing requires manual interaction"