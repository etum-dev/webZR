# webZR

please dont use this lol


Websocket testing framework i built for ... reasons i forgot
It is made for identifying if a target site is using websockets. 
Gradually, it will also hold some quick n' easy vuln checks.
It was primarily built as a way to identify websockets at scale for research.


### Identification methodologies
1. Javascript parsing
2. CSP Headers
3. (Kind of useless, not run by default) Simple dir and sub fuzz
4. Shodan lookup



### Options     

| Flag | Type | Description |
|------|------|-------------|
| -file <path> | string | Path to file containing list of domains (one per line) |
| -fuzz <mode>| int | Fuzzing mode: `0` (none), `1` (basic), `2` (custom header), `3` (mutation) |
| `-test` | bool | Start local WebSocket test server |


WebZR supports three ways to provide domains:
      
  1. **File input**: Use `-d` flag
     ```bash
     webzr -d domains.txt
     ```
  
  2. **Command-line arguments**: Pass domains directly
     ```bash
     webzr example.com api.example.com
     ```
  
  3. **Stdin pipe**: Pipe domains from other tools
     ```bash
     cat domains.txt | webzr
     echo "example.com" | webzr
     ```

___
ROADMAP:
[ ] Add Shodan Enumeration
[ ] Add a config parser (see utils/IDEA.md)
[ ] Add caching for already scanned domains


