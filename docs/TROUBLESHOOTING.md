# Troubleshooting Guide

## Common Issues and Solutions

### Installation Issues

#### Go is not installed
**Error:** `command not found: go`

**Solution:**
1. Download Go from https://golang.org/dl/
2. Follow installation instructions for your OS
3. Verify installation: `go version`
4. Ensure Go is in your PATH

#### Dependencies fail to download
**Error:** `go mod download` fails

**Solution:**
```bash
# Clear module cache
go clean -modcache

# Re-download dependencies
go mod download

# If still failing, try
go mod tidy
```

#### Permission denied when running start.sh
**Error:** `permission denied: ./start.sh`

**Solution:**
```bash
chmod +x start.sh
./start.sh
```

### Server Issues

#### Port 8080 already in use
**Error:** `listen tcp :8080: bind: address already in use`

**Solution:**
```bash
# Find process using port 8080
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Kill the process
kill -9 <PID>  # macOS/Linux
taskkill /PID <PID> /F  # Windows

# Or change port in main.go
```

#### Server crashes immediately
**Possible causes:**
1. Missing data files
2. Incorrect working directory

**Solution:**
```bash
# Ensure you're in the project root
cd /path/to/GoSim

# Check for required directories
ls -la data/
ls -la web/

# Run from project root
go run cmd/server/main.go
```

#### WebSocket connection fails
**Error:** WebSocket connection refused

**Solutions:**
1. Check firewall settings
2. Ensure server is running
3. Check browser console for errors
4. Try different browser
5. Check CORS settings

### Frontend Issues

#### Page loads but board doesn't appear
**Causes:**
1. JavaScript errors
2. Canvas not supported
3. CSS not loading

**Solutions:**
1. Open browser DevTools (F12)
2. Check Console for errors
3. Check Network tab for failed requests
4. Clear browser cache (Ctrl+Shift+R)
5. Try different browser

#### Can't place stones on the board
**Causes:**
1. Game not started
2. Not your turn
3. JavaScript error

**Solutions:**
1. Click a game mode first (vs AI, Local, etc.)
2. Check whose turn it is
3. Check browser console for errors
4. Refresh page and try again

#### AI doesn't make moves
**Causes:**
1. AI endpoint not working
2. Network request failing
3. Game state issue

**Solutions:**
1. Check Network tab in DevTools
2. Verify `/api/ai-move` returns data
3. Check server logs for errors
4. Try different AI difficulty

#### Multiplayer not working
**Causes:**
1. WebSocket connection failed
2. Room code issues
3. Network problems

**Solutions:**
1. Check WebSocket connection in Network tab
2. Both players must use same room code
3. Ensure both on same network (for local)
4. Check firewall settings

### Game Logic Issues

#### Illegal moves being allowed
**Check:**
1. Ko rule implementation
2. Suicide rule implementation
3. Board state synchronization

**Debug:**
```go
// Add logging in rules.go
log.Printf("Validating move at %v for %v", p, color)
log.Printf("Board state: %+v", g.Board)
```

#### Scoring seems wrong
**Check:**
1. Scoring method (Chinese vs Japanese)
2. Komi value
3. Territory calculation
4. Dead stone marking

**Debug:**
```go
// Add logging in scoring.go
log.Printf("Territory: %v", territory)
log.Printf("Captures: %v", captures)
log.Printf("Final scores: B:%f W:%f", score.Black, score.White)
```

#### Game freezes
**Causes:**
1. Infinite loop in AI
2. Browser hanging
3. WebSocket timeout

**Solutions:**
1. Check browser console
2. Restart server
3. Clear browser cache
4. Check for infinite loops in code

### Build Issues

#### Build fails with module errors
**Solution:**
```bash
# Update go.mod and go.sum
go mod tidy

# Upgrade dependencies
go get -u ./...

# Rebuild
go build -o gosim cmd/server/main.go
```

#### Tests failing
**Common issues:**
1. Changed API without updating tests
2. Timing issues in tests
3. Missing test data

**Solution:**
```bash
# Run tests verbosely
go test -v ./...

# Run specific test
go test -v -run TestBoardCreation ./pkg/game

# Check coverage
go test -cover ./...
```

### Browser Compatibility

#### Works in Chrome but not Firefox
**Check:**
1. Console errors specific to Firefox
2. CSS compatibility
3. JavaScript syntax

**Solutions:**
1. Use standard JavaScript features
2. Test in multiple browsers
3. Add polyfills if needed

#### Mobile browser issues
**Common problems:**
1. Touch events not working
2. Board too small
3. CSS layout broken

**Solutions:**
1. Add touch event handlers
2. Make responsive design
3. Test on actual devices

### Performance Issues

#### Game runs slowly
**Causes:**
1. AI taking too long
2. Too many DOM updates
3. Memory leaks

**Solutions:**
1. Profile with browser DevTools
2. Optimize AI algorithms
3. Reduce unnecessary renders
4. Check for memory leaks

#### High CPU usage
**Check:**
1. Infinite loops
2. Inefficient algorithms
3. Too frequent updates

**Solutions:**
```javascript
// Throttle updates
let updateScheduled = false;
function scheduleUpdate() {
    if (!updateScheduled) {
        updateScheduled = true;
        requestAnimationFrame(() => {
            update();
            updateScheduled = false;
        });
    }
}
```

### Data Issues

#### Puzzles not loading
**Check:**
1. File path correct
2. JSON format valid
3. Server endpoint working

**Debug:**
```bash
# Check puzzle files exist
ls -la data/puzzles/

# Validate JSON
python -m json.tool data/puzzles/beginner.json

# Test endpoint
curl http://localhost:8080/api/puzzles
```

#### Lessons not displaying
**Similar checks as puzzles:**
1. Verify file existence
2. Check JSON validity
3. Test API endpoint
4. Check browser console

### Development Environment

#### Hot reload not working
**For Go (using Air):**
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with air
air

# Check .air.toml configuration
```

**For Frontend:**
- Manually refresh browser
- Use browser auto-refresh extensions
- Set up webpack if needed

#### IDE issues
**VS Code:**
1. Install Go extension
2. Install gopls: `go install golang.org/x/tools/gopls@latest`
3. Reload VS Code

**GoLand/IntelliJ:**
1. Configure GOPATH
2. Enable Go modules
3. Index project

### Debugging Tips

#### Enable verbose logging
```go
// In main.go
import "log"

func init() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}
```

#### Browser debugging
```javascript
// Add breakpoints
debugger;

// Conditional logging
if (DEBUG) console.log('State:', state);

// Preserve logs
// Check "Preserve log" in DevTools Network/Console
```

#### Network debugging
```bash
# Monitor network traffic
tcpdump -i any -n port 8080  # Linux/macOS

# Check WebSocket frames
# Use browser DevTools Network tab
# Filter by WS to see WebSocket traffic
```

### Getting Help

If these solutions don't work:

1. **Search existing issues:** https://github.com/Prawal-Sharma/GoSim/issues
2. **Create new issue with:**
   - Error messages (full text)
   - Steps to reproduce
   - System information (OS, Go version, browser)
   - Screenshots if UI issue
   - Relevant code snippets
3. **Provide logs:**
   - Server logs
   - Browser console logs
   - Network traces if relevant

### Quick Fixes Checklist

- [ ] Go installed and in PATH?
- [ ] In correct directory?
- [ ] Dependencies downloaded?
- [ ] Server running?
- [ ] Using correct URL (http://localhost:8080)?
- [ ] Browser cache cleared?
- [ ] Firewall/antivirus not blocking?
- [ ] Port 8080 available?
- [ ] JavaScript enabled in browser?
- [ ] Using modern browser?

### Reset Everything

If all else fails, clean slate:
```bash
# Clean build
go clean -cache
go clean -modcache
rm -rf gosim gosim.exe

# Re-clone
cd ..
rm -rf GoSim
git clone https://github.com/Prawal-Sharma/GoSim.git
cd GoSim

# Fresh start
go mod download
go run cmd/server/main.go
```