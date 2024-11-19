# KISS WebChat

A lightweight, minimalistic web chat application built with Go, designed to be simple, efficient, and easy to deploy. This project demonstrates basic web server functionality, message handling, and user interaction with a clean UI.

## Features

- **Real-time chat:** Supports text and optional image sharing.
- **Message history:** Stores the last 100 messages in memory.
- **System notifications:** Automatically announces new users.
- **Admin controls:** Clear message history or shut down the server via simple commands.
- **HTML-based interface:** Accessible from any browser without additional software.
- **Cookie-based username retention:** Saves user names for convenience.

## Getting Started

### Prerequisites

- **Go (Golang)**: Version 1.18 or later.
- A terminal/command prompt for running the server.
- Optionally, a valid SSL certificate for HTTPS support.

### Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/DmitryStrunevsky/kisschat.git
   cd kisschat
   ```

2. Build the server:
   ```bash
   go build -o kisschat
   ```

3. Run the server:
   ```bash
   ./kisschat [port] [token]
   ```
   - Replace `[port]` with the desired port (default is `8000`).
   - Replace `[token]` with a custom token for admin actions (default is `chat`).

4. Open the chat interface in your browser:
   ```
   http://localhost:8000/chat
   ```

### Admin Actions

- **Clear Messages:** `/chat/clear`
- **Shutdown Server:** `/chat/shutdown`

Access these via your browser or tools like `curl`.

### HTTPS Setup (Optional)

For secure connections, provide an SSL certificate and key:
```bash
./kisschat [port] [token]
```
Ensure the files `server.crt` and `server.key` are in the same directory as the binary.

## Development

### Code Structure

- `ChatServer`: Core chat server implementation, manages message storage and user connections.
- `Message`: Represents a chat message, with support for text, images, and timestamps.
- Handlers: HTTP endpoints for chat interaction (`/messages`, `/clear`, `/shutdown`).

### Customize

- Modify the `MAX_MESSAGES` constant to adjust message history length.
- Update the HTML and CSS in the `getHTML()` function for UI changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Feel free to contribute, report issues, or suggest features via the [issues page](https://github.com/DmitryStrunevsky/kisschat/issues).
