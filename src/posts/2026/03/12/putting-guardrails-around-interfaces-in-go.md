---
date: 2026-03-12
tags:
- golang
- Tools
title: Putting guardrails around interfaces in Go
type: post
permalink: /2026/03/12/putting-guardrails-around-interfaces-in-go/
summary: Sometimes you don't want anyone else touching that.
---

Recently at work we had a new system that needed to communicate with an older system, who's data transfer format was even older than that and consisted of amalgamated JSON built up over years.  There was not time to rework or replace the communication format to something modern, and the older system still had to communicate with contemporary systems, so a complete replacement wasn't really possible anyway.

This second system was, however, written in Go, as was the new project.  That meant we could take the types from the existing application and use them in the new one, which reduces duplicated code and also ensure we have the weird and tangled types correct in both places.

However, we did not want to infect our new project with these types, as they matched the JSON structure and were not well formed or easy to work with.  We decided to add a package which did all the marshalling from external types to internal types, and a strict boundary line to keep them from bleeding over.

The problem then was, how do we enforce it?

We can certainly socialize the issue, and encourage PR's to check for it, but that's very dependent on keeping knowledge alive in the team, and vigilance in peer reviews. The new system won't change often, and this context will be lost over time.

We needed a way to enforce, or at the least loudly complain, about this in the codebase.  Something that we can add documentation around, and gives us a single point of contact to document.

The solution, for us, was [depguard](https://github.com/OpenPeeDeeP/depguard).  Depguard is a linter which lets you create allow or deny lists of packages for your codebase.  In our case, we used depguard through out meta-linter [golangci-lint](https://golangci-lint.run/)

We created a deny list based rule for the external packages (`pkg: github.com/chromaui/capture-exchange`), then carved out an exception for our internal package (`!**internal/exchange/*.go`) that handled the boundary.

```yaml
version: "2"
linters:
  enable:
    - depguard
  settings:
    depguard:
      rules:
        no-exchange:
          files:
            - "$all"
            - "!**/internal/exchange/*.go"
          allow: []
          deny:
            - pkg: "github.com/chromaui/capture-exchange"
              desc: "Use internal/exchange/ instead of importing capture-exchange directly."
```

If you attempt to use the external package elsewhere in the project, you now get a lint error, which you'll need to resolve and will (hopefully) lead you to the information you need to know why we have quarantined that package the way we did.

```
$ golangci-lint run
main.go:10:2: import 'github.com/chromaui/capture-exchange/pkg/types' is not allowed from list 'no-exchange': Use internal/exchange/ instead of importing capture-exchange directly. (depguard)
        "github.com/chromaui/capture-exchange/pkg/types"
        ^
```

It's something that can run on developer machines and CI, we can customize the error message, and we can then add lots of documentation to the internal package explaining why we did it this way.  Problem solved!
