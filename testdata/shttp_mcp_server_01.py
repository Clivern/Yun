#!/usr/bin/env -S uv run --quiet --script
# /// script
# requires-python = ">=3.10"
# dependencies = [
#     "fastmcp",
# ]
# ///
# Copyright 2025 Clivern. All rights reserved.
# Use of this source code is governed by the MIT
# license that can be found in the LICENSE file.

from fastmcp import FastMCP

mcp = FastMCP(name="Mut")


# Tools
@mcp.tool()
def hello(name: str) -> str:
    """Greets a person by name.

    Args:
        name: The name of the person to greet

    Returns:
        A greeting message
    """
    return f"Hello, {name}!"


@mcp.tool()
def add(a: int, b: int) -> int:
    """Adds two numbers together.

    Args:
        a: First number
        b: Second number

    Returns:
        The sum of a and b
    """
    return a + b


# Resources
@mcp.resource("info://server")
async def server_info() -> str:
    """Server Information resource.

    Returns:
        Information about this MCP server
    """
    return """Mut v1.0.0

This is a Model Context Protocol (MCP) server that provides:

- Tool: 'hello' - Greets a person by name
- Tool: 'add' - Adds two numbers together
- Prompt: 'greeting' - A friendly greeting prompt
- Prompt: 'farewell' - A friendly farewell prompt
- Resource: 'info://server' - Server information (this document)
- Resource: 'help://guide' - User guide and documentation

The server implements the MCP protocol version 2024-11-05 and communicates
via STDIO using JSON-RPC 2.0.

Server built with FastMCP.
"""


@mcp.resource("help://guide")
async def user_guide() -> str:
    """User Guide resource.

    Returns:
        User guide documentation
    """
    return """Mut User Guide

USAGE EXAMPLES:

1. Using the hello tool:
   Call the 'hello' tool with a name parameter to get a greeting.

2. Using the add tool:
   Call the 'add' tool with two integers to get their sum.

3. Using the greeting prompt:
   Use the 'greeting' prompt with an optional name for a friendly greeting.

4. Using the farewell prompt:
   Use the 'farewell' prompt with a name for a friendly goodbye.

SUPPORT:
For more information, visit: https://gofastmcp.com
"""


# Prompts
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
                "text": f"Hello, {name}! Welcome to Mut. How can I help you today?"
            }
        }
    ]


@mcp.prompt()
def farewell(name: str = "friend") -> list[dict]:
    """A friendly farewell prompt.

    Args:
        name: The name to say goodbye to (default: "friend")

    Returns:
        A list of messages with the farewell
    """
    return [
        {
            "role": "user",
            "content": {
                "type": "text",
                "text": f"Goodbye, {name}! Thank you for using Mut. Have a great day!"
            }
        }
    ]


if __name__ == "__main__":
    mcp.run(
        transport="http",
        host="127.0.0.1",
        port=8000,
        show_banner=False,
    )
