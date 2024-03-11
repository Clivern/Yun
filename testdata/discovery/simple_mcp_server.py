#!/usr/bin/env -S uv run --quiet --script
# /// script
# requires-python = ">=3.10"
# dependencies = []
# ///
"""
Simple MCP Server with one tool, one resource, and one prompt.
Run with: uv run simple_mcp_server.py
"""

import json
import sys
from typing import Any, Dict


def send_response(response: Dict[str, Any]) -> None:
    """Send a JSON-RPC response to stdout."""
    sys.stdout.write(json.dumps(response) + '\n')
    sys.stdout.flush()


def handle_initialize(request_id: int) -> Dict[str, Any]:
    """Handle initialize request."""
    return {
        'jsonrpc': '2.0',
        'id': request_id,
        'result': {
            'protocolVersion': '2024-11-05',
            'capabilities': {
                'tools': {},
                'resources': {},
                'prompts': {}
            },
            'serverInfo': {
                'name': 'simple-mcp-server',
                'version': '1.0.0'
            }
        }
    }


def handle_tools_list(request_id: int) -> Dict[str, Any]:
    """Handle tools/list request."""
    return {
        'jsonrpc': '2.0',
        'id': request_id,
        'result': {
            'tools': [{
                'name': 'calculate',
                'description': 'Performs basic arithmetic calculations',
                'inputSchema': {
                    'type': 'object',
                    'properties': {
                        'operation': {
                            'type': 'string',
                            'enum': ['add', 'subtract', 'multiply', 'divide'],
                            'description': 'The arithmetic operation to perform'
                        },
                        'a': {
                            'type': 'number',
                            'description': 'First number'
                        },
                        'b': {
                            'type': 'number',
                            'description': 'Second number'
                        }
                    },
                    'required': ['operation', 'a', 'b']
                }
            }]
        }
    }


def handle_tools_call(request_id: int, params: Dict[str, Any]) -> Dict[str, Any]:
    """Handle tools/call request."""
    tool_name = params.get('name')
    arguments = params.get('arguments', {})

    if tool_name == 'calculate':
        operation = arguments.get('operation')
        a = arguments.get('a', 0)
        b = arguments.get('b', 0)

        try:
            if operation == 'add':
                result = a + b
            elif operation == 'subtract':
                result = a - b
            elif operation == 'multiply':
                result = a * b
            elif operation == 'divide':
                if b == 0:
                    return {
                        'jsonrpc': '2.0',
                        'id': request_id,
                        'result': {
                            'content': [{
                                'type': 'text',
                                'text': 'Error: Division by zero'
                            }],
                            'isError': True
                        }
                    }
                result = a / b
            else:
                result = 'Unknown operation'

            return {
                'jsonrpc': '2.0',
                'id': request_id,
                'result': {
                    'content': [{
                        'type': 'text',
                        'text': f'Result: {a} {operation} {b} = {result}'
                    }],
                    'isError': False
                }
            }
        except Exception as e:
            return {
                'jsonrpc': '2.0',
                'id': request_id,
                'result': {
                    'content': [{
                        'type': 'text',
                        'text': f'Error: {str(e)}'
                    }],
                    'isError': True
                }
            }

    return {
        'jsonrpc': '2.0',
        'id': request_id,
        'error': {
            'code': -32601,
            'message': f'Unknown tool: {tool_name}'
        }
    }


def handle_prompts_list(request_id: int) -> Dict[str, Any]:
    """Handle prompts/list request."""
    return {
        'jsonrpc': '2.0',
        'id': request_id,
        'result': {
            'prompts': [{
                'name': 'greeting',
                'description': 'A friendly greeting prompt',
                'arguments': [{
                    'name': 'name',
                    'description': 'The name to greet',
                    'required': False
                }]
            }]
        }
    }


def handle_prompts_get(request_id: int, params: Dict[str, Any]) -> Dict[str, Any]:
    """Handle prompts/get request."""
    prompt_name = params.get('name')
    arguments = params.get('arguments', {})

    if prompt_name == 'greeting':
        name = arguments.get('name', 'friend')
        return {
            'jsonrpc': '2.0',
            'id': request_id,
            'result': {
                'description': 'A friendly greeting',
                'messages': [{
                    'role': 'user',
                    'content': {
                        'type': 'text',
                        'text': f'Hello, {name}! How can I help you today?'
                    }
                }]
            }
        }

    return {
        'jsonrpc': '2.0',
        'id': request_id,
        'error': {
            'code': -32601,
            'message': f'Unknown prompt: {prompt_name}'
        }
    }


def handle_resources_list(request_id: int) -> Dict[str, Any]:
    """Handle resources/list request."""
    return {
        'jsonrpc': '2.0',
        'id': request_id,
        'result': {
            'resources': [{
                'uri': 'info://server',
                'name': 'Server Information',
                'description': 'Information about this MCP server',
                'mimeType': 'text/plain'
            }]
        }
    }


def handle_resources_read(request_id: int, params: Dict[str, Any]) -> Dict[str, Any]:
    """Handle resources/read request."""
    uri = params.get('uri')

    if uri == 'info://server':
        info = """Simple MCP Server v1.0.0

This is a simple Model Context Protocol (MCP) server that provides:

- Tool: 'calculate' - Performs basic arithmetic operations (add, subtract, multiply, divide)
- Prompt: 'greeting' - A friendly greeting prompt
- Resource: 'info://server' - This information document

The server implements the MCP protocol version 2024-11-05 and communicates
via stdio using JSON-RPC 2.0.
"""
        return {
            'jsonrpc': '2.0',
            'id': request_id,
            'result': {
                'contents': [{
                    'uri': uri,
                    'mimeType': 'text/plain',
                    'text': info
                }]
            }
        }

    return {
        'jsonrpc': '2.0',
        'id': request_id,
        'error': {
            'code': -32602,
            'message': f'Unknown resource: {uri}'
        }
    }


def handle_request(request: Dict[str, Any]) -> Dict[str, Any]:
    """Route request to appropriate handler."""
    method = request.get('method')
    request_id = request.get('id')
    params = request.get('params', {})

    handlers = {
        'initialize': lambda: handle_initialize(request_id),
        'tools/list': lambda: handle_tools_list(request_id),
        'tools/call': lambda: handle_tools_call(request_id, params),
        'prompts/list': lambda: handle_prompts_list(request_id),
        'prompts/get': lambda: handle_prompts_get(request_id, params),
        'resources/list': lambda: handle_resources_list(request_id),
        'resources/read': lambda: handle_resources_read(request_id, params),
    }

    handler = handlers.get(method)
    if handler:
        return handler()

    return {
        'jsonrpc': '2.0',
        'id': request_id,
        'error': {
            'code': -32601,
            'message': f'Method not found: {method}'
        }
    }


def main():
    """Main server loop."""
    try:
        for line in sys.stdin:
            line = line.strip()
            if not line:
                continue

            try:
                request = json.loads(line)

                # Skip notifications (requests without id)
                if 'id' not in request or request['id'] is None:
                    continue

                response = handle_request(request)
                send_response(response)

            except json.JSONDecodeError as e:
                sys.stderr.write(f'JSON decode error: {e}\n')
                sys.stderr.flush()
            except Exception as e:
                sys.stderr.write(f'Error processing request: {e}\n')
                sys.stderr.flush()

    except KeyboardInterrupt:
        pass
    except Exception as e:
        sys.stderr.write(f'Fatal error: {e}\n')
        sys.stderr.flush()


if __name__ == '__main__':
    main()

