package main

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MAX_MESSAGES = 100
)

var chat_token = "chat"

type Message struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Image string `json:"image"` // New field for image URL
	Time  string `json:"time"`
}

type ChatServer struct {
	messages      *ring.Ring
	messagesMux   sync.RWMutex
	announcedIPs  map[string]bool
	serverRunning bool
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		messages:      ring.New(MAX_MESSAGES),
		announcedIPs:  make(map[string]bool),
		serverRunning: true,
	}
}

func (s *ChatServer) addMessage(msg Message) {
	s.messagesMux.Lock()
	defer s.messagesMux.Unlock()

	s.messages.Value = msg
	s.messages = s.messages.Next()
}

func (s *ChatServer) getMessages() []Message {
	s.messagesMux.RLock()
	defer s.messagesMux.RUnlock()

	var messages []Message
	s.messages.Do(func(x interface{}) {
		if x != nil {
			messages = append(messages, x.(Message))
		}
	})
	return messages
}

func (s *ChatServer) clearMessages() {
	s.messagesMux.Lock()
	s.messages = ring.New(MAX_MESSAGES)
	s.messagesMux.Unlock()
	s.addMessage(Message{
		Name: "System",
		Text: "Message history has been cleared",
		Time: time.Now().Format("15:04:05"),
	})
}

func (s *ChatServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	// Check token
	if len(r.URL.Path) < len(chat_token) || r.URL.Path[1:len(chat_token)+1] != chat_token {
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}

	// Handle new IP announcement
	clientIP := r.RemoteAddr
	if strings.Contains(r.RemoteAddr, ":") {
		clientIPparts := strings.Split(r.RemoteAddr, ":")
		clientIP = clientIPparts[0]
	}

	if !s.announcedIPs[clientIP] {
		s.announcedIPs[clientIP] = true
		s.addMessage(Message{
			Name: "System",
			Text: fmt.Sprintf("New user joined from IP: %s", clientIP),
			Time: time.Now().Format("15:04:05"),
		})
	}

	path := r.URL.Path[len(chat_token)+1:]

	switch {
	case r.Method == "GET" && path == "/shutdown":
		s.serverRunning = false
		fmt.Fprintf(w, "Server shutting down...")
		go func() {
			time.Sleep(time.Second)
			os.Exit(0)
		}()

	case r.Method == "GET" && path == "/clear":
		s.clearMessages()
		fmt.Fprintf(w, "Messages cleared")

	case r.Method == "GET" && path == "/messages":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.getMessages())

	case r.Method == "POST" && path == "":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		message := r.FormValue("message")
		image := r.FormValue("image") // New field for image URL

		if name != "" && (message != "" || image != "") {
			s.addMessage(Message{
				Name:  html.EscapeString(name),
				Text:  html.EscapeString(message),
				Image: html.EscapeString(image),
				Time:  time.Now().Format("15:04:05"),
			})
		}

		http.Redirect(w, r, "/"+chat_token, http.StatusSeeOther)

	default:
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, getHTML())
	}
}

func getHTML() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>KISS WebChat</title>
    <style>
        body {
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            font-family: Arial, sans-serif;
        }
        .message {
            padding: 10px;
            margin: 5px 0;
            background: #f0f0f0;
            border-radius: 5px;
        }
        .message.system {
            background: #e8f5e9;
            border-left: 3px solid #4caf50;
            font-size: 0.5em;
            padding: 5px;
        }
        .name {
            font-weight: bold;
        }
        .name.system {
            color: #2e7d32;
        }
        .time {
            color: #666;
            font-size: 0.8em;
        }
        form {
            margin: 20px 0;
        }
        input[type="text"] {
            padding: 5px;
            margin-right: 10px;
        }
        input[type="submit"], button {
            padding: 5px 15px;
        }
        #messages {
            height: 400px;
            overflow-y: auto;
            border: 1px solid #ccc;
            padding: 10px;
        }
        #controls {
            display: flex;
            gap: 10px;
            margin-top: 10px;
        }
        .admin-btn {
            border: none;
            border-radius: 3px;
            cursor: pointer;
            color: white;
            font-weight: bold;
        }
        .kill-btn {
            background: #ff5252;
        }
        .kill-btn:hover {
            background: #ff1744;
        }
        .clear-btn {
            background: #ffa726;
        }
        .clear-btn:hover {
            background: #fb8c00;
        }
        .image {
            max-width: 100%;
            height: auto;
            border-radius: 5px;
            margin-top: 5px;
        }
    </style>
    <script>
        // Load saved username from cookie
        document.addEventListener('DOMContentLoaded', () => {
            const savedName = getCookie('username');
            if (savedName) {
                document.querySelector('input[name="name"]').value = savedName;
            }
        });

        function getCookie(name) {
            const value = ` + "`" + `; ${document.cookie}` + "`" + `;
            const parts = value.split(` + "`" + `; ${name}=` + "`" + `);
            if (parts.length === 2) return parts.pop().split(';').shift();
        }

        function setCookie(name, value) {
            document.cookie = ` + "`" + `${name}=${value};path=/;max-age=31536000` + "`" + `;
        }

        function escapeHtml(unsafe) {
            return unsafe
                .replace(/&/g, "&amp;")
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;")
                .replace(/"/g, "&quot;")
                .replace(/'/g, "&#039;");
        }

        function scrollToBottom() {
            const messages = document.getElementById('messages');
            messages.scrollTop = messages.scrollHeight;
        }

        function updateMessages() {
            fetch('` + chat_token + `/messages')
                .then(response => response.json())
                .then(messages => {
                    const messagesDiv = document.getElementById('messages');
                    if(JSON.stringify(messages) === JSON.stringify(window.messages)) return;
                    window.messages = messages;
                    messagesDiv.innerHTML = messages.map(msg => {
                        const isSystem = msg.name === 'System';
                        const messageContent = msg.text && ` + "`" + `<span class="text">${escapeHtml(msg.text)}</span>` + "`" + `;
                        const imageContent = msg.image && ` + "`" + `<img src="${escapeHtml(msg.image)}" alt="Image" class="image">` + "`" + `;
                        return ` + "`" + `
                            <div class="message ${isSystem ? 'system' : ''}">
                                <span class="name ${isSystem ? 'system' : ''}">${escapeHtml(msg.name)}:</span>
                                ${messageContent || ''}
                                ${imageContent || ''}
                                <span class="time">[${escapeHtml(msg.time)}]</span>
                            </div>
                        ` + "`" + `;
                    }).join('');
                    scrollToBottom();
                });
        }

        async function killServer() {
            if (confirm('Are you sure you want to shut down the server?')) {
                try {
                    await fetch('/` + chat_token + `/shutdown');
                    alert('Server is shutting down...');
                } catch (e) {
                    // Error is expected as server shuts down
                }
            }
        }

        async function clearMessages() {
            if (confirm('Are you sure you want to clear all messages?')) {
                try {
                    await fetch('/` + chat_token + `/clear');
                    updateMessages();
                } catch (e) {
                    alert('Failed to clear messages');
                }
            }
        }

        // Update messages every 2 seconds
        setInterval(updateMessages, 2000);

        // Submit form without page reload
        document.addEventListener('DOMContentLoaded', () => {
            const form = document.querySelector('form');
            form.onsubmit = async (e) => {
                e.preventDefault();
                const formData = new FormData(form);

                // Save username to cookie
                const username = formData.get('name');
                setCookie('username', username);

                await fetch('/` + chat_token + `', {
                    method: 'POST',
                    body: new URLSearchParams(formData)
                });

                // Only reset message field, keep the name
                document.querySelector('input[name="message"]').value = '';
                document.querySelector('input[name="image"]').value = '';
                updateMessages();
            };
        });
    </script>
</head>
<body onload="updateMessages()">
    <h1>KISS WebChat</h1>
    <div id="messages"></div>
    <form>
        <input type="text" name="name" placeholder="Your name" required>
        <input type="text" name="message" placeholder="Your message" required>
        <input type="text" name="image" placeholder="Image URL (optional)">
        <input type="submit" value="Send">
    </form>
    <div id="controls">
        <button onclick="clearMessages()" class="admin-btn clear-btn">Clear Messages</button>
        <button onclick="killServer()" class="admin-btn kill-btn">Shutdown Server</button>
    </div>
</body>
</html>`
}

func main() {
	port := 8000
	if len(os.Args) > 1 {
		if p, err := strconv.Atoi(os.Args[1]); err == nil {
			port = p
		}
	}

	if len(os.Args) > 2 {
		chat_token = os.Args[2]
	}

	server := NewChatServer()

	fmt.Printf("\nChat URL: https://localhost:%d/%s\n", port, chat_token)
	fmt.Printf("Admin token for shutdown/clear: %s\n\n", chat_token)

	http.HandleFunc("/", server.handleRequest)

	if port != 8000 {
		// Configure TLS
		err := http.ListenAndServeTLS(fmt.Sprintf(":%d", port), "server.crt", "server.key", nil)
		if err != nil {
			log.Fatal("ListenAndServeTLS: ", err)
		}
	} else {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
