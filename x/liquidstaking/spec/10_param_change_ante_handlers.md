<!-- order: 10 -->

# Param Change Ante Handlers

## Slashing Param Change Limit Decorator

The liquidstaking module works closely with the slashing params. (e.g. MinimumCollateral constant is calculated based on the slashing params)
To reduce unexpected risks, it is important to reduce the maximum slashing penalty that can theoretically occur.
