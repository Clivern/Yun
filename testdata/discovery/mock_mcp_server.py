#!/usr/bin/env python3
import json
import sys

def send_response(response):
    sys.stdout.write(json.dumps(response) + '\n')
    sys.stdout.flush()

def handle_request(request):
    method = request.get('method')
    req_id = request.get('id')

    if method == 'initialize':
        return {
            'jsonrpc': '2.0',
            'id': req_id,
            'result': {
                'protocolVersion': '2024-11-05',
                'capabilities': {},
                'serverInfo': {
                    'name': 'mock-server',
                    'version': '1.0.0',
                    'protocolVersion': '2024-11-05'
                }
            }
        }
    elif method == 'tools/list':
        return {
            'jsonrpc': '2.0',
            'id': req_id,
            'result': {
                'tools': [{
                    'name': 'test_tool',
                    'description': 'A test tool',
                    'inputSchema': {'type': 'object'}
                }]
            }
        }
    elif method == 'tools/call':
        return {
            'jsonrpc': '2.0',
            'id': req_id,
            'result': {
                'content': [{
                    'type': 'text',
                    'text': 'Tool result for test'
                }],
                'isError': False
            }
        }
    elif method == 'prompts/list':
        return {
            'jsonrpc': '2.0',
            'id': req_id,
            'result': {
                'prompts': [{
                    'name': 'test_prompt',
                    'description': 'A test prompt'
                }]
            }
        }
    elif method == 'prompts/get':
        return {
            'jsonrpc': '2.0',
            'id': req_id,
            'result': {
                'description': 'Test prompt',
                'messages': [{
                    'role': 'user',
                    'content': {'type': 'text', 'text': 'Test message'}
                }]
            }
        }
    elif method == 'resources/list':
        return {
            'jsonrpc': '2.0',
            'id': req_id,
            'result': {
                'resources': [{
                    'uri': 'test://resource',
                    'name': 'test_resource',
                    'description': 'A test resource',
                    'mimeType': 'text/plain'
                }]
            }
        }
    elif method == 'resources/read':
        return {
            'jsonrpc': '2.0',
            'id': req_id,
            'result': {
                'contents': [{
                    'uri': 'test://resource',
                    'mimeType': 'text/plain',
                    'text': 'Resource content'
                }]
            }
        }
    else:
        return {
            'jsonrpc': '2.0',
            'id': req_id,
            'error': {
                'code': -32601,
                'message': 'Method not found'
            }
        }

# Main loop
try:
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue

        request = json.loads(line)

        # Handle notifications (no response needed)
        if 'id' not in request or request['id'] is None:
            continue

        response = handle_request(request)
        send_response(response)

except KeyboardInterrupt:
    pass
except Exception as e:
    sys.stderr.write(f'Error: {e}\n')

