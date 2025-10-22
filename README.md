# 🎭 Anonymous Real-Time Chat (Omegle-style)

Pure real-time anonymous chat service built with Go and WebSockets. No database, no registration - just instant anonymous conversations!

## ✨ Features

- 🚀 **Pure Real-Time** - All communication in memory, no database
- 🎲 **Random Matching** - Get paired with random strangers
- 💬 **Instant Messaging** - WebSocket-based real-time chat
- 🔄 **Skip Feature** - Don't like your partner? Skip to the next one!
- 🔒 **Anonymous** - No registration, no data storage
- ⚡ **Lightweight** - Minimal dependencies, fast startup

## 🏗️ Architecture

```
Client (WebSocket) 
    ↓
WebSocket Handler
    ↓
Router (Message Type Based)
    ↓
Handlers (FindMatch, SendMessage, NextStranger, StopChat)
    ↓
Hub + MatchingService (In-Memory)
    ↓
Paired Stranger (WebSocket)
```

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- No database required!

### Installation

```bash
# Clone or download the project
cd real-time-service

# Install dependencies
go mod download

# Run the service
go run main.go
```

Server will start on `http://localhost:8080`

## 📡 WebSocket API

Connect to: `ws://localhost:8080/ws`

### Message Types

#### 1. Find Match (Start Chatting)
```json
{
  "type": "findMatch"
}
```

**Responses:**
- If waiting: `{"type": "searching"}`
- If matched: `{"type": "strangerJoined", "pairId": "uuid"}`

#### 2. Send Message
```json
{
  "type": "sendMessage",
  "text": "Hello stranger!"
}
```

Partner receives:
```json
{
  "type": "message",
  "text": "Hello stranger!",
  "userId": "sender-uuid",
  "pairId": "pair-uuid",
  "timestamp": "2025-10-22T10:30:00Z"
}
```

#### 3. Next Stranger (Skip)
```json
{
  "type": "nextStranger"
}
```

Current partner gets: `{"type": "strangerLeft"}`
You get: `{"type": "searching"}` or `{"type": "strangerJoined"}`

#### 4. Stop Chat
```json
{
  "type": "stopChat"
}
```

Disconnects from current chat and removes you from matching queue.

## 🧪 Testing with JavaScript

```html
<!DOCTYPE html>
<html>
<head>
    <title>Anonymous Chat Test</title>
</head>
<body>
    <div id="status">Disconnected</div>
    <div id="messages"></div>
    <input type="text" id="input" placeholder="Type a message...">
    <button onclick="sendMessage()">Send</button>
    <button onclick="findMatch()">Start Chat</button>
    <button onclick="nextStranger()">Next</button>

    <script>
        let ws = new WebSocket('ws://localhost:8080/ws');
        
        ws.onopen = () => {
            document.getElementById('status').textContent = 'Connected';
        };
        
        ws.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            const div = document.getElementById('messages');
            
            if (msg.type === 'strangerJoined') {
                div.innerHTML += '<p><b>Stranger connected! Start chatting.</b></p>';
            } else if (msg.type === 'message') {
                div.innerHTML += '<p><b>Stranger:</b> ' + msg.text + '</p>';
            } else if (msg.type === 'strangerLeft') {
                div.innerHTML += '<p><b>Stranger disconnected.</b></p>';
            } else if (msg.type === 'searching') {
                div.innerHTML += '<p><i>Searching for stranger...</i></p>';
            }
        };
        
        function findMatch() {
            ws.send(JSON.stringify({type: 'findMatch'}));
        }
        
        function sendMessage() {
            const input = document.getElementById('input');
            ws.send(JSON.stringify({
                type: 'sendMessage',
                text: input.value
            }));
            document.getElementById('messages').innerHTML += 
                '<p><b>You:</b> ' + input.value + '</p>';
            input.value = '';
        }
        
        function nextStranger() {
            ws.send(JSON.stringify({type: 'nextStranger'}));
        }
    </script>
</body>
</html>
```

## 📦 Project Structure

```
real-time-service/
├── main.go                          # Entry point
├── config.json                      # Configuration
├── configuration/
│   └── configuration.go             # Config loader
├── handlers/
│   ├── ws.go                        # WebSocket handler
│   └── wsrouter/
│       ├── router.go                # Message router
│       └── handlers/
│           ├── find_match_handler.go
│           ├── send_handler.go
│           ├── next_stranger_handler.go
│           └── stop_chat_handler.go
├── hubs/
│   └── main_hub.go                  # Connection hub
├── services/
│   └── matching_service.go          # Pair matching logic (in-memory)
├── models/
│   ├── client.go
│   ├── chat_pair.go                 # Pair model
│   ├── message.go
│   └── incoming_message.go
├── middlewares/
│   └── auth_middleware.go           # Simple session middleware
└── interfaces/
    └── container_interface.go       # DI interface
```

## 🎯 Key Components

### MatchingService
- Maintains waiting queue of users looking for chat
- Creates pairs when two users are available
- Manages active pairs in memory
- Thread-safe with mutex locks

### Hub
- Manages all connected WebSocket clients
- Routes messages between paired users
- Handles disconnections and notifications
- Integrates with MatchingService

### WebSocket Handlers
- **FindMatchHandler**: Matches users with strangers
- **SendHandler**: Forwards messages to partner
- **NextStrangerHandler**: Ends current chat and finds new partner
- **StopChatHandler**: Gracefully ends chat session

## ⚙️ Configuration

Edit `config.json`:

```json
{
  "httpPort": ":8080"
}
```

## 🔧 Development

### Run with auto-reload (using air)
```bash
go install github.com/cosmtrek/air@latest
air
```

### Build
```bash
go build -o chat-service
```

### Run
```bash
./chat-service
```

## 📝 TODO / Future Features

- [ ] Add HTML/CSS/JS web interface
- [ ] Typing indicators
- [ ] Connection statistics (users online, pairs active)
- [ ] Interest tags for better matching
- [ ] Rate limiting
- [ ] Profanity filter
- [ ] Message history (in-memory, per session)
- [ ] Video/audio chat support

## 🤝 Contributing

Feel free to submit issues and pull requests!

## 📄 License

MIT License - feel free to use this for learning or building your own chat service!

---

**Made with ❤️ using Go, Gin, and WebSockets**
