#!/bin/bash
# Manual testing script for Gabel

echo "=== Gabel Manual Test Suite ==="
echo

echo "1. Testing help flag..."
./gabel -h
echo

echo "2. Testing with invalid arguments..."
./gabel 2>&1 | head -5
echo

echo "3. Testing with invalid repo format..."
./gabel invalid another-invalid 2>&1
echo

echo "4. Testing debug output..."
./gabel -d golang/go golang/go 2>&1 | grep DEBUG | head -5
echo

echo "5. Testing GitHub CLI check..."
PATH="" ./gabel foo/bar foo/baz 2>&1
echo

echo "=== Tests complete ==="
echo "Note: Full interactive testing requires manual interaction"