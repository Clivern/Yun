#!/usr/bin/env python3
import json
import sys

for line in sys.stdin:
    request = json.loads(line.strip())
    if 'id' in request:
        response = {
            'jsonrpc': '2.0',
            'id': request['id'],
            'error': {
                'code': -32600,
                'message': 'Invalid Request',
                'data': 'Test error'
            }
        }
        sys.stdout.write(json.dumps(response) + '\n')
        sys.stdout.flush()

