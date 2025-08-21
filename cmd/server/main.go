package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"collaborative-markdown-editor/internal/client"
	"collaborative-markdown-editor/internal/hub"
	"collaborative-markdown-editor/internal/user"
)

var h *hub.Hub

func main() {
	// Create hub
	h = hub.NewHub()
	go h.Run()

	// HTTP routes
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/room/", serveRoom)



	fmt.Println("Server starting on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

// serveHome serves the home page
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Serve the home page template
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CollabWrite - Real-time Collaborative Markdown Editor</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .container {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 20px;
            padding: 3rem;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
            text-align: center;
            max-width: 500px;
            width: 90%;
        }

        .logo {
            width: 80px;
            height: 80px;
            margin: 0 auto 2rem;
            background: linear-gradient(45deg, #667eea, #764ba2);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-size: 2rem;
            font-weight: bold;
        }

        h1 {
            color: #2d3748;
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
            font-weight: 700;
        }

        .subtitle {
            color: #718096;
            margin-bottom: 2rem;
            font-size: 1.1rem;
        }

        .button {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 1rem 2rem;
            border-radius: 12px;
            font-size: 1.1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            margin: 0.5rem;
            min-width: 200px;
        }

        .button:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgba(102, 126, 234, 0.3);
        }

        .button.secondary {
            background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
        }

        .input-group {
            margin: 1.5rem 0;
        }

        .input-wrapper {
            position: relative;
            margin-bottom: 1rem;
        }

        input[type="text"] {
            width: 100%;
            padding: 1rem 1rem 1rem 3rem;
            border: 2px solid #e2e8f0;
            border-radius: 12px;
            font-size: 1rem;
            transition: all 0.3s ease;
            background: white;
        }

        input[type="text"]:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }

        .input-icon {
            position: absolute;
            left: 1rem;
            top: 50%;
            transform: translateY(-50%);
            color: #a0aec0;
            font-size: 1.2rem;
        }

        .features {
            display: flex;
            justify-content: center;
            gap: 2rem;
            margin-top: 2rem;
            padding-top: 2rem;
            border-top: 1px solid #e2e8f0;
        }

        .feature {
            text-align: center;
        }

        .feature-icon {
            font-size: 1.5rem;
            margin-bottom: 0.5rem;
        }

        .feature-text {
            font-size: 0.9rem;
            color: #718096;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">📝</div>
        <h1>CollabWrite</h1>
        <p class="subtitle">Real-time collaborative markdown editing made simple</p>

        <div class="button secondary" onclick="createRoom()">✨ Create New Room</div>

        <div class="input-group">
            <div class="input-wrapper">
                <span class="input-icon">🔗</span>
                <input type="text" id="roomId" placeholder="Enter room code to join">
            </div>
            <button class="button" onclick="joinRoom()">Join Room</button>
        </div>

        <div class="features">
            <div class="feature">
                <div class="feature-icon">⚡</div>
                <div class="feature-text">Real-time sync</div>
            </div>
            <div class="feature">
                <div class="feature-icon">👥</div>
                <div class="feature-text">Multi-user</div>
            </div>
            <div class="feature">
                <div class="feature-icon">📝</div>
                <div class="feature-text">Markdown</div>
            </div>
        </div>
    </div>

    <script>
        function createRoom() {
            // Generate a random room ID
            const roomId = Math.random().toString(36).substring(2, 8);
            window.location.href = '/room/' + roomId;
        }

        function joinRoom() {
            const roomId = document.getElementById('roomId').value.trim();
            if (roomId) {
                window.location.href = '/room/' + roomId;
            } else {
                alert('Please enter a room code');
            }
        }

        // Add enter key support
        document.getElementById('roomId').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                joinRoom();
            }
        });
    </script>
</body>
</html>`

	t := template.Must(template.New("home").Parse(tmpl))
	t.Execute(w, nil)
}

// serveRoom serves the editor page for a specific room
func serveRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Path[len("/room/"):]
	if roomID == "" {
		http.Error(w, "Room ID required", http.StatusBadRequest)
		return
	}

	// Get current content of the room
	currentContent := h.GetRoomContent(roomID)

	// Serve the editor template
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Room: {{.RoomID}} - CollabWrite</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f8fafc;
            height: 100vh;
            overflow: hidden;
        }

        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 1rem 2rem;
            display: flex;
            justify-content: between;
            align-items: center;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }

        .header h1 {
            font-size: 1.5rem;
            font-weight: 600;
        }

        .share-link {
            background: rgba(255,255,255,0.2);
            padding: 0.5rem 1rem;
            border-radius: 8px;
            font-family: 'Monaco', 'Menlo', monospace;
            font-size: 0.9rem;
            backdrop-filter: blur(10px);
        }

        .main-container {
            display: flex;
            height: calc(100vh - 73px);
            position: relative;
        }

        .users-sidebar {
            width: 250px;
            background: white;
            border-left: 1px solid #e2e8f0;
            padding: 1rem;
            overflow-y: auto;
        }

        .users-header {
            font-weight: 600;
            color: #2d3748;
            margin-bottom: 1rem;
            font-size: 1.1rem;
        }

        .user-item {
            display: flex;
            align-items: center;
            padding: 0.5rem;
            border-radius: 8px;
            margin-bottom: 0.5rem;
            transition: all 0.2s ease;
        }

        .user-item:hover {
            background: #f8fafc;
        }

        .user-avatar {
            width: 32px;
            height: 32px;
            border-radius: 50%;
            margin-right: 0.75rem;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
            font-size: 0.8rem;
        }

        .user-name {
            font-weight: 500;
            color: #2d3748;
        }

        .user-status {
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background: #48bb78;
            margin-left: auto;
        }

        .cursor-indicator {
            position: absolute;
            height: 20px;
            width: 2px;
            background: red;
            animation: blink 1s infinite;
            z-index: 100;
            pointer-events: none;
        }

        @keyframes blink {
            0%, 50% { opacity: 1; }
            51%, 100% { opacity: 0; }
        }

        .editor-panel {
            flex: 1;
            padding: 2rem;
            background: white;
            border-right: 1px solid #e2e8f0;
            display: flex;
            flex-direction: column;
        }

        .preview-panel {
            flex: 1;
            padding: 2rem;
            background: #f8fafc;
            overflow: auto;
            border-left: 1px solid #e2e8f0;
        }

        .panel-header {
            display: flex;
            align-items: center;
            margin-bottom: 1rem;
            padding-bottom: 0.5rem;
            border-bottom: 1px solid #e2e8f0;
        }

        .panel-title {
            font-size: 1.1rem;
            font-weight: 600;
            color: #2d3748;
            margin-left: 0.5rem;
        }

        .panel-icon {
            font-size: 1.2rem;
        }

        textarea {
            flex: 1;
            width: 100%;
            padding: 1rem;
            border: 2px solid #e2e8f0;
            border-radius: 8px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 14px;
            resize: none;
            outline: none;
            transition: all 0.3s ease;
            background: #fafafa;
            line-height: 1.6;
        }

        textarea:focus {
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
            background: white;
        }

        .preview {
            line-height: 1.7;
            color: #2d3748;
        }

        .status {
            position: fixed;
            bottom: 20px;
            right: 20px;
            background: rgba(255,255,255,0.95);
            color: #2d3748;
            padding: 0.5rem 1rem;
            border-radius: 20px;
            font-size: 0.8rem;
            font-weight: 500;
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            backdrop-filter: blur(10px);
            border: 1px solid rgba(0,0,0,0.1);
            z-index: 1000;
        }

        .status.connected {
            background: rgba(34,197,94,0.1);
            color: #16a34a;
            border-color: rgba(34,197,94,0.2);
        }

        .status.disconnected {
            background: rgba(239,68,68,0.1);
            color: #dc2626;
            border-color: rgba(239,68,68,0.2);
        }

        /* Markdown styling */
        .markdown-body {
            font-size: 16px;
            line-height: 1.6;
        }

        .markdown-body h1 {
            font-size: 2em;
            margin-top: 0.67em;
            margin-bottom: 0.67em;
            color: #1a202c;
            border-bottom: 1px solid #e2e8f0;
            padding-bottom: 0.3em;
        }

        .markdown-body h2 {
            font-size: 1.5em;
            margin-top: 0.83em;
            margin-bottom: 0.83em;
            color: #2d3748;
        }

        .markdown-body h3 {
            font-size: 1.17em;
            margin-top: 1em;
            margin-bottom: 1em;
            color: #4a5568;
        }

        .markdown-body p {
            margin-top: 1em;
            margin-bottom: 1em;
            color: #2d3748;
        }

        .markdown-body code {
            background: #f7fafc;
            padding: 0.2em 0.4em;
            border-radius: 3px;
            font-family: 'Monaco', 'Menlo', monospace;
            font-size: 0.9em;
            border: 1px solid #e2e8f0;
        }

        .markdown-body pre {
            background: #f7fafc;
            padding: 1rem;
            border-radius: 6px;
            overflow: auto;
            border: 1px solid #e2e8f0;
        }

        .markdown-body blockquote {
            border-left: 4px solid #667eea;
            padding-left: 1rem;
            margin: 1rem 0;
            color: #718096;
            font-style: italic;
        }

        .markdown-body ul, .markdown-body ol {
            margin: 1rem 0;
            padding-left: 2rem;
        }

        .markdown-body li {
            margin: 0.5rem 0;
        }

        .copy-btn {
            background: transparent;
            border: none;
            cursor: pointer;
            padding: 0.25rem;
            border-radius: 4px;
            transition: all 0.2s ease;
        }

        .copy-btn:hover {
            background: rgba(255,255,255,0.3);
        }

        @media (max-width: 768px) {
            .main-container {
                flex-direction: column;
                height: calc(100vh - 120px);
            }

            .editor-panel {
                height: 50vh;
                border-right: none;
                border-bottom: 1px solid #e2e8f0;
            }

            .preview-panel {
                height: 50vh;
                border-left: none;
                border-top: 1px solid #e2e8f0;
            }

            .header {
                flex-direction: column;
                gap: 1rem;
                text-align: center;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>🖊️ Room: {{.RoomID}}</h1>
        <div class="share-link">
            <button class="copy-btn" onclick="copyLink()">📋</button>
            {{.Host}}/room/{{.RoomID}}
        </div>
    </div>

    <div class="main-container">
        <div class="editor-panel">
            <div class="panel-header">
                <span class="panel-icon">📝</span>
                <span class="panel-title">Editor</span>
            </div>
            <textarea id="editor" placeholder="Start typing Markdown here...&#10;&#10;**Bold text**&#10;*Italic text*&#10;&#10;### Lists&#10;- Item 1&#10;- Item 2&#10;&#10;### Code&#10;&#96;inline code&#96;&#10;&#10;### Links&#10;[Google](https://google.com)">{{.Content}}</textarea>
        </div>
        <div class="preview-panel">
            <div class="panel-header">
                <span class="panel-icon">👁️</span>
                <span class="panel-title">Preview</span>
            </div>
            <div id="preview" class="markdown-body preview"></div>
        </div>
        <div class="users-sidebar">
            <div class="users-header">👥 Active Users</div>
            <div id="users-list">
                <!-- Users will be populated by JavaScript -->
            </div>
        </div>
    </div>

    <div class="status" id="status">
        <span>🔌</span> Connecting...
    </div>

    <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
    <script>
        const editor = document.getElementById('editor');
        const preview = document.getElementById('preview');
        const status = document.getElementById('status');
        const usersList = document.getElementById('users-list');
        const roomID = '{{.RoomID}}';

        let ws;
        let reconnectInterval;
        let lastContent = '';
        let currentUser = null;
        let cursorUpdateInterval;

        function connect() {
            ws = new WebSocket('ws://' + window.location.host + '/ws/' + roomID);

            ws.onopen = function(event) {
                status.innerHTML = '<span>✅</span> Connected';
                status.className = 'status connected';
                clearInterval(reconnectInterval);

                // Request user info and room users
                fetchUsers();
                startCursorUpdates();
            };

            ws.onmessage = function(event) {
                try {
                    // Try to parse as JSON first (for user list updates)
                    const data = JSON.parse(event.data);
                    if (data.type === 'userList') {
                        // This is a user list update
                        updateUsersList(data.users);
                        return;
                    }
                } catch (e) {
                    // Not JSON, treat as regular content
                    const content = event.data;
                    if (content !== editor.value) {
                        const cursorPos = editor.selectionStart;
                        editor.value = content;
                        lastContent = content;
                        updatePreview();

                        // Restore cursor position approximately
                        if (cursorPos <= content.length) {
                            editor.selectionStart = editor.selectionEnd = cursorPos;
                        }
                    }
                }
            };

            ws.onclose = function(event) {
                status.innerHTML = '<span>❌</span> Disconnected';
                status.className = 'status disconnected';
                clearInterval(cursorUpdateInterval);

                // Try to reconnect every 3 seconds
                reconnectInterval = setInterval(connect, 3000);
            };

            ws.onerror = function(error) {
                console.error('WebSocket error:', error);
            };
        }

        function updatePreview() {
            preview.innerHTML = marked.parse(editor.value);
        }

        function copyLink() {
            navigator.clipboard.writeText(window.location.href).then(function() {
                showNotification('Room link copied to clipboard!');
            });
        }

        function showNotification(message) {
            const notification = document.createElement('div');
            notification.style.cssText = 'position: fixed; top: 20px; right: 20px; background: #48bb78; color: white; padding: 0.5rem 1rem; border-radius: 8px; z-index: 1000;';
            notification.textContent = message;
            document.body.appendChild(notification);

            setTimeout(() => {
                notification.remove();
            }, 3000);
        }

        let typingTimer;
        const doneTypingInterval = 300;

        editor.addEventListener('input', function() {
            clearTimeout(typingTimer);
            typingTimer = setTimeout(doneTyping, doneTypingInterval);
            updatePreview();
        });

        function doneTyping() {
            const content = editor.value;
            if (content !== lastContent && ws && ws.readyState === WebSocket.OPEN) {
                ws.send(content);
                lastContent = content;
            }
        }

        function fetchUsers() {
            // Start with empty list - users will be populated via WebSocket messages
            updateUsersList([]);
        }

        function updateUsersList(users) {
            usersList.innerHTML = '';
            users.forEach(user => {
                const userElement = document.createElement('div');
                userElement.className = 'user-item';
                userElement.innerHTML = '<div class="user-avatar" style="background: ' + user.color + '">' +
                    user.username.charAt(0).toUpperCase() +
                    '</div><span class="user-name">' + user.username + '</span><div class="user-status"></div>';
                usersList.appendChild(userElement);
            });
        }

        function startCursorUpdates() {
            cursorUpdateInterval = setInterval(() => {
                if (ws && ws.readyState === WebSocket.OPEN) {
                    // Send cursor position update
                    const cursorPos = editor.selectionStart;
                    // In a real implementation, you'd send this to the server
                    console.log('Cursor position:', cursorPos);
                }
            }, 100);
        }

        // Initialize
        connect();
        updatePreview();
        lastContent = editor.value;

        // Auto-save cursor position
        editor.addEventListener('keyup', function() {
            localStorage.setItem('cursorPos_' + roomID, editor.selectionStart);
        });

        // Restore cursor position
        const savedCursorPos = localStorage.getItem('cursorPos_' + roomID);
        if (savedCursorPos) {
            setTimeout(() => {
                editor.selectionStart = editor.selectionEnd = parseInt(savedCursorPos);
            }, 100);
        }

        // Add keyboard shortcuts
        document.addEventListener('keydown', function(e) {
            if (e.ctrlKey || e.metaKey) {
                switch(e.key) {
                    case 's':
                        e.preventDefault();
                        copyLink();
                        break;
                    case 'b':
                        e.preventDefault();
                        wrapText('**', '**');
                        break;
                    case 'i':
                        e.preventDefault();
                        wrapText('*', '*');
                        break;
                    case 'k':
                        e.preventDefault();
                        wrapText('[', '](url)');
                        break;
                }
            }
        });

        function wrapText(before, after) {
            const start = editor.selectionStart;
            const end = editor.selectionEnd;
            const selectedText = editor.value.substring(start, end);
            const newText = before + selectedText + after;

            editor.setRangeText(newText, start, end);
            editor.selectionStart = editor.selectionEnd = start + newText.length;
            editor.focus();
            updatePreview();
            doneTyping();
        }
    </script>
</body>
</html>`

	data := struct {
		RoomID string
		Content string
		Host    string
	}{
		RoomID:  roomID,
		Content: currentContent,
		Host:    r.Host,
	}

	t := template.Must(template.New("room").Parse(tmpl))
	t.Execute(w, data)
}


