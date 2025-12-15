#!/bin/bash

mkdir -p logs
OUTPUT_FILE="logs/commit-logs.txt"

echo "Generating commit logs to $OUTPUT_FILE..."

# git log format:
# %H: Commit Hash
# %ad: Author Date (ISO 8601)
# %B: Raw Body (Message)
# --name-status: Show changed files with status (A/M/D)

git log --date=iso --pretty=format:"------------------------------------------------------------------------%nCommit: %H%nDate:   %ad%n%nMessage:%n%B%nFiles:" --name-status > "$OUTPUT_FILE"

echo "Done."
