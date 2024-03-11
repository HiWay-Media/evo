---
layout: default title: Keycloak Example parent: Examples
---

# Keycloak

Example usage of keycloak:

```go
package main

import (
	"github.com/hiway-media/evo"
	"github.com/hiway-media/evo/apps/keycloak"
)

func main() {

	evo.Setup()
	keycloak.Register("https://auth.example.com", "realm", "client")
	evo.Run()
}
```

[Keycloak Example Project](https://github.com/hiway-media/examples/tree/master/keycloak_example)
