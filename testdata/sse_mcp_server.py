#!/usr/bin/env -S uv run --quiet --script
# /// script
# requires-python = ">=3.10"
# dependencies = [
#     "fastmcp",
#     "uvicorn[standard]",
#     "fastapi",
# ]
# ///
"""
SSE MCP Server with one tool, one resource, and one prompt using FastMCP.
Run with: uv run cache/sse_mcp_server.py
"""

from fastmcp import FastMCP

PORT = 8080

# Create FastMCP server
mcp = FastMCP("SSE MCP Server")


# Define the calculate tool
@mcp.tool()
def calculate(operation: str, a: float, b: float) -> str:
    """Performs basic arithmetic calculations.

    Args:
        operation: The arithmetic operation to perform (add, subtract, multiply, divide)
        a: First number
        b: Second number

    Returns:
        The result of the calculation as a string
    """
    try:
        if operation == "add":
            result = a + b
        elif operation == "subtract":
            result = a - b
        elif operation == "multiply":
            result = a * b
        elif operation == "divide":
            if b == 0:
                return "Error: Division by zero"
            result = a / b
        else:
            return f"Unknown operation: {operation}"

        return f"Result: {a} {operation} {b} = {result}"
    except Exception as e:
        return f"Error: {str(e)}"


# Define the greeting prompt
@mcp.prompt()
def greeting(name: str = "friend") -> list[dict]:
    """A friendly greeting prompt.

    Args:
        name: The name to greet (default: "friend")

    Returns:
        A list of messages with the greeting
    """
    return [
        {
            "role": "user",
            "content": {
                "type": "text",
                "text": f"Hello, {name}! How can I help you today?"
            }
        }
    ]


# Add a resource
@mcp.resource("info://server")
async def server_info() -> str:
    """Server Information resource.

    Returns:
        Information about this SSE MCP server
    """
    return """SSE MCP Server v1.0.0

This is a Server-Sent Events (SSE) Model Context Protocol (MCP) server that provides:

- Tool: 'calculate' - Performs basic arithmetic operations (add, subtract, multiply, divide)
- Prompt: 'greeting' - A friendly greeting prompt
- Resource: 'info://server' - This information document

The server implements the MCP protocol version 2024-11-05 and communicates
via HTTP POST using JSON-RPC 2.0.

Server runs on port 8080 by default.
Server built with FastMCP.
"""


# Create the FastAPI/ASGI app from FastMCP
app = mcp.http_app()


if __name__ == '__main__':
    import uvicorn
    print(f"SSE MCP Server starting on http://localhost:{PORT}")
    print(f"Endpoint: http://localhost:{PORT}/mcp")
    uvicorn.run(app, host='0.0.0.0', port=PORT)
