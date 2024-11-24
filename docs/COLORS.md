# Colors
Colors for some of the output has been added, if the terminal supports it.

It can be disabled with `TERM=` or `NO_COLOR=`. `NO_COLOR` can be set to any
value other than empty to disable colors. This `NO_COLOR` any value behavior
seems to come from the shell, and isn't written into this code.

Examples:
```
TERM= ./att-fiber-gateway-info --help
NO_COLOR=1 ./att-fiber-gateway-info --help
NO_COLOR=false ./att-fiber-gateway-info --help # When set to anything, no color
```
