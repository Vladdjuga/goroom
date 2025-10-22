# 🎭 Anonymous Real-Time Chat (Omegle-style)

Pure real-time anonymous chat service built with **Go, Gin, WebSockets, and MVC architecture**. No database, no registration - just instant anonymous conversations with a beautiful modern UI!

## ✨ Features

- 🚀 **Pure Real-Time** - All communication in memory, no database
- 🎲 **Random Matching** - Get paired with random strangers
- 💬 **Instant Messaging** - WebSocket-based real-time chat
- 🔄 **Skip Feature** - Don't like your partner? Skip to the next one!
- 🔒 **Anonymous** - No registration, no data storage
- ⚡ **Lightweight** - Minimal dependencies, fast startup
- 🎨 **Modern UI** - Beautiful gradient design with smooth animations
- 📱 **Responsive** - Works perfectly on desktop and mobile
- 🏗️ **MVC Architecture** - Clean separation of concerns

## 🎨 Screenshots

**Home Page**: Beautiful landing page with features and call-to-action
**Chat Room**: Modern chat interface with real-time messaging

## 🏗️ Architecture

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ HTTP/WebSocket
┌──────▼──────────────────────┐
│    Gin Web Server (MVC)     │
├─────────────────────────────┤
│ Controllers (HTTP Routes)   │
│  ├─ HomeController          │
│  └─ ChatController          │
├─────────────────────────────┤
│ WebSocket Handler           │
│  └─ Router (Message Types)  │
├─────────────────────────────┤
│ Handlers                    │
│  ├─ FindMatchHandler        │
│  ├─ SendHandler             │
│  ├─ NextStrangerHandler     │
│  └─ StopChatHandler         │
├─────────────────────────────┤
│ Hub + MatchingService       │
│  └─ In-Memory Pairs         │
└─────────────────────────────┘
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

### Open in Browser

1. **Home Page**: Open `http://localhost:8080` 
2. **Chat Page**: Click "Start Chatting Now" or go to `http://localhost:8080/chat`
3. **Open Multiple Tabs**: Open the chat page in 2+ browser tabs to test matching

## 🌐 Pages

### Home Page (`/`)
- Beautiful landing page with gradient design
- Feature showcase
- Call-to-action button
- Real-time statistics (coming soon)

### Chat Page (`/chat`)
- Modern chat interface
- Real-time status indicator
- Message bubbles with timestamps
- Control buttons (Start, Next, Stop)
- Smooth animations and transitions

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
├── main.go                          # Entry point with routes
├── config.json                      # Configuration
├── README.md                        # Documentation
│
├── controllers/                     # MVC Controllers
│   ├── home_controller.go           # Home page
│   └── chat_controller.go           # Chat page
│
├── views/                           # MVC Views
│   ├── templates/
│   │   ├── layout.html              # Base layout
│   │   ├── home.html                # Home page template
│   │   └── chat.html                # Chat page template
│   └── static/
│       ├── css/
│       │   └── style.css            # Modern gradient design
│       └── js/
│           └── chat.js              # WebSocket client
│
├── handlers/
│   ├── ws.go                        # WebSocket handler
│   └── wsrouter/
│       ├── router.go                # Message router
│       └── handlers/
│           ├── find_match_handler.go
│           ├── send_handler.go
│           ├── next_stranger_handler.go
│           └── stop_chat_handler.go
│
├── hubs/
│   └── main_hub.go                  # Connection hub
│
├── services/
│   └── matching_service.go          # Pair matching (in-memory)
│
├── models/
│   ├── client.go
│   ├── chat_pair.go
│   ├── message.go
│   └── incoming_message.go
│
├── middlewares/
│   └── auth_middleware.go           # Simple session
│
├── interfaces/
│   └── container_interface.go       # DI interface
│
├── providers/
│   └── main_providers.go            # DI container
│
└── configuration/
    └── configuration.go             # Config loader
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

- [x] MVC architecture with controllers
- [x] Beautiful modern UI with gradients
- [x] Server-side template rendering
- [x] Responsive design for mobile
- [ ] Real-time statistics on home page
- [ ] Typing indicators
- [ ] Rate limiting
- [ ] Profanity filter
- [ ] Interest tags for better matching
- [ ] Dark mode toggle
- [ ] Sound notifications
- [ ] Video/audio chat support (WebRTC)
- [ ] Chat logs download option

## 🎨 Design Features

- **Gradient Backgrounds**: Modern purple/pink gradients
- **Glassmorphism**: Frosted glass effects on cards
- **Smooth Animations**: Fade-ins, slide-ups, hover effects
- **Message Bubbles**: Different styles for you/stranger
- **Status Indicators**: Animated dots showing connection state
- **Responsive Layout**: Mobile-first design
- **Custom Scrollbars**: Styled for better aesthetics

## 📡 WebSocket API

Connect to: `ws://localhost:8080/ws`

### Message Types

#### 1. Find Match
```json
{"type": "findMatch"}
```

#### 2. Send Message
```json
{"type": "sendMessage", "text": "Hello!"}
```

#### 3. Next Stranger
```json
{"type": "nextStranger"}
```

#### 4. Stop Chat
```json
{"type": "stopChat"}
```

## 🤝 Contributing

Feel free to submit issues and pull requests!

## 📄 License

MIT License - feel free to use this for learning or building your own chat service!

---

## 🚀 Deployment

### Deploy to Render.com (Free)

This project is ready to deploy to [Render.com](https://render.com) with zero configuration!

**Steps:**

1. **Push to GitHub:**
   ```bash
   git add .
   git commit -m "Ready for deployment"
   git push origin main
   ```

2. **Create Render Account:**
   - Go to https://render.com
   - Sign up with GitHub

3. **Deploy:**
   - Click "New +" → "Web Service"
   - Connect your GitHub repository
   - Render will auto-detect `render.yaml`
   - Click "Create Web Service"
   - Wait 2-3 minutes for build ☕

4. **Done!** 🎉
   - Your app will be live at: `https://your-app-name.onrender.com`
   - WebSocket will work at: `wss://your-app-name.onrender.com/ws`

**Note:** Free tier sleeps after 15 minutes of inactivity. First request after sleep takes ~30 seconds.

### Deploy to Railway.app

Alternative option with $5 free credit monthly:

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Initialize project
railway init

# Deploy
railway up
```

### Environment Variables

The app automatically uses `PORT` environment variable when deployed. No manual configuration needed!

**Local development:**
```bash
# Uses config.json
go run main.go
```

**Production (Render/Railway):**
```bash
# Uses PORT env variable
PORT=8080 go run main.go
```

---

**Made with ❤️ using Go, Gin, and WebSockets**
