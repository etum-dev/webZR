# webZR

please dont use this lol


Websocket testing framework i built for ... reasons i forgot


### Identification methodologies
Simple dir and sub fuzz
TODO: Javascript parsing
TODO: CSP Headers

### Options     

| Flag | Type | Description |
|------|------|-------------|
| -file <path> | string | Path to file containing list of domains (one per line) |
| -fuzz <mode>| int | Fuzzing mode: `0` (none), `1` (basic), `2` (custom header), `3` (mutation) |
| `-test` | bool | Start local WebSocket test server |


WebZR supports three ways to provide domains:
      
  1. **File input**: Use `-file` flag
     ```bash
     webzr -file domains.txt
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
TODOs:
- [ ] More vuln scans
- [ ] Optimize speed

