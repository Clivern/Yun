#!/usr/bin/env python3
import sys
import time

# Just read but never respond
try:
    for line in sys.stdin:
        time.sleep(10)  # Simulate slow response
except KeyboardInterrupt:
    pass

