#!/usr/bin/env python3
import os
import json
import sys

for line in sys.stdin:
    request = json.loads(line.strip())
    if request.get('method') == 'initialize':
        cwd = os.getcwd()
        response = {
            'jsonrpc': '2.0',
            'id': request['id'],
            'result': {
                'protocolVersion': '2024-11-05',
                'capabilities': {},
                'serverInfo': {
                    'name': 'cwd-test',
                    'version': '1.0.0'
                }
            }
        }
        # The test will verify the working directory via configuration
        sys.stdout.write(json.dumps(response) + '\n')
        sys.stdout.flush()

