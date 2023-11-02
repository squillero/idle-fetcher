# The Idle Fetcher

A small program that shows

```
HOSTNAME [LOCAL IP/PUBLIC IP]
```

That is, a simple excuse for writing few lines of Go.

* Reasonably smart output (when directly connected to the Internet or not connected at all)
* Reasonably portable (tested on different Linux flavors, OSX, and Windows)
* Reasonably fast (exploits multiple methods to get IPs in parallel and caches the result for some time)

**Copyright © 2023 by Giovanni Squillero**  
Distributed under a [Zero-Clause BSD](https://tldrlegal.com/license/bsd-0-clause-license) License (SPDX: [0BSD](https://spdx.org/licenses/0BSD.html)), which allows unlimited freedom without the requirement to include legal notices.
