// WebSocket Chat Client
let ws = null;
let currentState = 'disconnected'; // disconnected, searching, chatting
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    setupEventListeners();
    connect();
});

// Setup event listeners
function setupEventListeners() {
    const input = document.getElementById('messageInput');
    if (input) {
        input.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !input.disabled) {
                sendMessage();
            }
        });
    }
}

// Connect to WebSocket server
function connect() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    console.log('Connecting to:', wsUrl);
    ws = new WebSocket(wsUrl);
    
    ws.onopen = () => {
        console.log('✅ Connected to server');
        updateStatus('connected', 'Connected');
        reconnectAttempts = 0;
        
        // Enable start button
        const startBtn = document.getElementById('startBtn');
        if (startBtn) startBtn.disabled = false;
    };
    
    ws.onmessage = (event) => {
        try {
            const msg = JSON.parse(event.data);
            console.log('📨 Received:', msg);
            handleMessage(msg);
        } catch (error) {
            console.error('❌ Error parsing message:', error);
        }
    };
    
    ws.onerror = (error) => {
        console.error('❌ WebSocket error:', error);
        showSystemMessage('Connection error occurred');
    };
    
    ws.onclose = () => {
        console.log('🔌 Disconnected from server');
        updateStatus('disconnected', 'Disconnected');
        disableAllButtons();
        
        // Try to reconnect
        if (reconnectAttempts < maxReconnectAttempts) {
            reconnectAttempts++;
            const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 10000);
            console.log(`🔄 Reconnecting in ${delay}ms (attempt ${reconnectAttempts}/${maxReconnectAttempts})`);
            showSystemMessage(`Reconnecting in ${delay/1000}s...`);
            setTimeout(connect, delay);
        } else {
            showSystemMessage('Unable to connect. Please refresh the page.');
        }
    };
}

// Handle incoming messages
function handleMessage(msg) {
    switch(msg.type) {
        case 'searching':
            updateStatus('searching', 'Searching for stranger...');
            currentState = 'searching';
            showSystemMessage('🔍 Looking for a stranger to chat with...');
            
            // Disable start, enable stop
            setButtonStates({ start: false, next: false, stop: true });
            break;
            
        case 'strangerJoined':
            updateStatus('chatting', 'Chatting with stranger');
            currentState = 'chatting';
            showSystemMessage('✨ Stranger connected! Say hi!');
            
            // Enable input and buttons
            enableChatInput();
            setButtonStates({ start: false, next: true, stop: true });
            
            // Focus input
            const input = document.getElementById('messageInput');
            if (input) input.focus();
            break;
            
        case 'message':
            addStrangerMessage(msg.text, msg.timestamp);
            break;
            
        case 'strangerLeft':
            updateStatus('connected', 'Stranger left');
            currentState = 'connected';
            showSystemMessage('👋 Stranger disconnected');
            
            // Disable input, enable start button
            disableChatInput();
            setButtonStates({ start: true, next: false, stop: false });
            break;
            
        default:
            console.warn('Unknown message type:', msg.type);
    }
}

// Send message
function sendMessage() {
    const input = document.getElementById('messageInput');
    if (!input) return;
    
    const text = input.value.trim();
    
    if (text && ws && ws.readyState === WebSocket.OPEN && currentState === 'chatting') {
        ws.send(JSON.stringify({
            type: 'sendMessage',
            text: text
        }));
        
        addYourMessage(text);
        input.value = '';
        input.focus();
    }
}

// Start chatting
function startChat() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        // Clear previous messages
        const messagesDiv = document.getElementById('messages');
        if (messagesDiv) {
            messagesDiv.innerHTML = '';
        }
        
        ws.send(JSON.stringify({ type: 'findMatch' }));
        showSystemMessage('🔍 Looking for a stranger...');
    } else {
        showSystemMessage('⚠️ Not connected to server. Please wait...');
    }
}

// Next stranger
function nextStranger() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'nextStranger' }));
        showSystemMessage('🔄 Looking for a new stranger...');
        
        // Clear messages
        const messagesDiv = document.getElementById('messages');
        if (messagesDiv) {
            messagesDiv.innerHTML = '';
        }
        
        disableChatInput();
    }
}

// Stop chat
function stopChat() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'stopChat' }));
        showSystemMessage('🛑 Chat stopped');
        
        currentState = 'connected';
        disableChatInput();
        setButtonStates({ start: true, next: false, stop: false });
    }
}

// Clear chat
function clearChat() {
    const messagesDiv = document.getElementById('messages');
    if (messagesDiv) {
        messagesDiv.innerHTML = '<div class="welcome-message">' +
            '<div class="welcome-icon">👋</div>' +
            '<h2>Welcome to Anonymous Chat!</h2>' +
            '<p>Click "Start Chatting" to be paired with a random stranger</p>' +
            '</div>';
    }
}

// UI Helper Functions
function updateStatus(state, text) {
    const statusBadge = document.getElementById('statusBadge');
    if (!statusBadge) return;
    
    // Remove old classes
    statusBadge.classList.remove('connected', 'searching', 'chatting');
    
    // Add new class
    if (state !== 'disconnected') {
        statusBadge.classList.add(state);
    }
    
    // Update text
    const statusText = statusBadge.querySelector('.status-text');
    if (statusText) {
        statusText.textContent = text;
    }
}

function addMessage(className, text, timestamp) {
    const messagesDiv = document.getElementById('messages');
    if (!messagesDiv) return;
    
    // Remove welcome message if exists
    const welcome = messagesDiv.querySelector('.welcome-message');
    if (welcome) {
        welcome.remove();
    }
    
    const msgDiv = document.createElement('div');
    msgDiv.className = `message ${className}`;
    
    const bubble = document.createElement('div');
    bubble.className = 'message-bubble';
    
    const textP = document.createElement('p');
    textP.className = 'message-text';
    textP.textContent = text;
    bubble.appendChild(textP);
    
    if (timestamp) {
        const time = new Date(timestamp);
        const timeDiv = document.createElement('div');
        timeDiv.className = 'message-time';
        timeDiv.textContent = time.toLocaleTimeString();
        bubble.appendChild(timeDiv);
    }
    
    msgDiv.appendChild(bubble);
    messagesDiv.appendChild(msgDiv);
    
    // Scroll to bottom
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

function addYourMessage(text) {
    addMessage('you', text, new Date().toISOString());
}

function addStrangerMessage(text, timestamp) {
    addMessage('stranger', text, timestamp);
}

function showSystemMessage(text) {
    const messagesDiv = document.getElementById('messages');
    if (!messagesDiv) return;
    
    // Remove welcome message if exists
    const welcome = messagesDiv.querySelector('.welcome-message');
    if (welcome) {
        welcome.remove();
    }
    
    const msgDiv = document.createElement('div');
    msgDiv.className = 'system-message';
    
    const bubble = document.createElement('div');
    bubble.className = 'system-bubble';
    bubble.textContent = text;
    
    msgDiv.appendChild(bubble);
    messagesDiv.appendChild(msgDiv);
    
    // Scroll to bottom
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

function enableChatInput() {
    const input = document.getElementById('messageInput');
    const sendBtn = document.getElementById('sendBtn');
    
    if (input) {
        input.disabled = false;
        input.placeholder = 'Type your message...';
    }
    if (sendBtn) sendBtn.disabled = false;
}

function disableChatInput() {
    const input = document.getElementById('messageInput');
    const sendBtn = document.getElementById('sendBtn');
    
    if (input) {
        input.disabled = true;
        input.placeholder = 'Start a chat to send messages...';
        input.value = '';
    }
    if (sendBtn) sendBtn.disabled = true;
}

function setButtonStates(states) {
    const buttons = {
        start: document.getElementById('startBtn'),
        next: document.getElementById('nextBtn'),
        stop: document.getElementById('stopBtn')
    };
    
    for (const [key, enabled] of Object.entries(states)) {
        if (buttons[key]) {
            buttons[key].disabled = !enabled;
        }
    }
}

function disableAllButtons() {
    setButtonStates({ start: false, next: false, stop: false });
    disableChatInput();
}
