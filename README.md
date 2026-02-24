
Installation:
-------------

Copy mcp-server-snapper to /usr/bin.

When using mcphost add it to the configuration:

```json
{
    "mcpServers": {
        "my-server-snapper": {
            "command": "/usr/bin/mcp-server-snapper"
        }
    }
}
```

If you do not run mcp servers as root you have to allow
mcp-server-snapper to use the snapper config nu adding the user
(e.g. "mcp-test") to ALLOW_USERS using (as root):

snapper set-config ALLOW_USERS=mcp-test


Examples:
---------

"Please list all snapshots in a table."

"Please create a snapshot before I update glibc."

"I have now updated glibc. Please create the post snapshot."

"How old is the newest snapshot?"

For the last example the LLM additional needs access to the current
time and date.

