# webZR

please dont use this lol


Websocket testing framework i built for ... reasons i forgot
It is made for identifying if a target site is using websockets. 
Gradually, it will also hold some quick n' easy vuln checks.
It was primarily built as a way to identify websockets at scale for research.


### Identification methodologies
1. Javascript parsing
2. CSP Headers
3. Enhanced subdomain scanning with worker pools
4. Directory and endpoint fuzzing
5. Shodan lookup (planned)



### Options     

| Flag | Type | Description |
|------|------|-------------|
| -d <path> | string | Single domain or path to file containing list of domains (one per line) |
| -fuzz <mode>| int | Fuzzing mode: `0` (none), `1` (basic), `2` (custom header), `3` (mutation) |
| -subdomain <mode> | string | Subdomain scanning: `off` (default), `basic`, `aggressive` |
| -subdomain-max <num> | int | Maximum number of subdomains to test (default: 50) |
| -subdomain-workers <num> | int | Number of concurrent subdomain workers (default: 8) |
| -m <mode> | string | General scan mode |
| -of <path> | string | Output file path (default: scan_results.json) |
| -debug | bool | Enable debug output |
| -v | bool | Enable verbose output |
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

### Subdomain Scanning Examples

**Basic subdomain scanning** (fast, stops on first WebSocket found):
```bash
webzr -subdomain basic example.com
```

**Aggressive subdomain scanning** (comprehensive, finds all WebSocket subdomains):
```bash
webzr -subdomain aggressive -subdomain-max 100 -subdomain-workers 15 example.com
```

**Combined with file input**:
```bash
webzr -d domains.txt -subdomain basic -of results.json
```

**Piped from subdomain enumeration tools**:
```bash
subfinder -d example.com | webzr -subdomain aggressive
```

___
ROADMAP:
[✅] Enhanced subdomain scanning with worker pools
[✅] Configurable subdomain scanning modes (basic/aggressive)
[✅] Concurrent subdomain testing with prioritization
[ ] Add Shodan Enumeration
[ ] Add a config parser (see utils/IDEA.md)  
[ ] Add caching for already scanned domains
[ ] Integration with external subdomain tools (subfinder, amass)


