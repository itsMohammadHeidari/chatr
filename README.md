# Chatr: A TCP-Based TUI Chatbot

This repository contains **Chatr**, a Text User Interface (TUI) chatbot that communicates over TCP. Built in Go, Chatr allows multiple clients to connect to a server and exchange messages in real-time through a command-line TUI. This project serves as the final assignment for the Network lecture under the guidance of **Professor [Hamid Haj Seyyed Javadi](https://www.researchgate.net/profile/Hamid-Haj-Seyyed-Javadi)**.

---

## Table of Contents

- [Chatr: A TCP-Based TUI Chatbot](#chatr-a-tcp-based-tui-chatbot)
  - [Table of Contents](#table-of-contents)
  - [Project Overview](#project-overview)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Running the Server](#running-the-server)
    - [Running the Client](#running-the-client)
    - [Environment Variables](#environment-variables)
    - [TCP Settings](#tcp-settings)
    - [Interacting via TUI](#interacting-via-tui)
  - [Contributing](#contributing)
  - [License](#license)
  - [Contact](#contact)

---

## Project Overview

**Chatr** is a client-server chat application written in Go that uses **raw TCP sockets** to send and receive messages. The project demonstrates practical networking concepts by:

- Establishing a persistent TCP connection between a server and multiple clients.
- Broadcasting messages from one client to all connected clients.
- Dynamically tracking and displaying online users in a TUI.
- Providing a fluid command-line interface experience.

Key features include:

- **TUI-based chat client** to send and receive messages.
- **Server broadcasting** to all connected clients in real-time.
- **User list management** (join and leave notifications).
- **Configurable** through environment variables or command-line flags.

This project is especially relevant to the **Network** lecture, as it highlights low-level TCP socket management and concurrency.

---

## Prerequisites

Before installing or running this project, make sure you have the following:

1. **Go (1.21+ recommended)**  
    You can check your Go version with:

    ```bash
    go version
    ```

2. **Git** (to clone the repository).
3. **Basic knowledge of TCP networking** (ports, hosts, firewalls).

---

## Installation

Follow these steps to set up Chatr on your local machine:

1. **Clone the repository**:

    ```bash
    git clone https://github.com/itsmohammadheidari/chatr.git
    cd chatr
    ```

2. **Initialize Go modules (optional if already present)**:

    ```bash
    go mod tidy
    ```

    This will download and verify all the required dependencies listed in `go.mod` and `go.sum`.

3. **(Optional) Create a `.env` file** in the project root if you want to override default configurations via environment variables. For example:

    ```env
    HOST=127.0.0.1
    PORT=8080
    USERNAME=MyUsername
    ```

    > **Note**: If you don'    t create a `.env` file, Chatr will use the defaults (127.0.0.1:8080) or the values provided via flags.

---

## Usage

### Running the Server

1. **Navigate** to the project directory:

    ```bash
    cd chatr
    ```

2. **Start the server**:

    ```bash
    go run main.go server
    ```

    Alternatively, if you have built the binary (e.g., `go build -o chatr`):

    ```bash
    ./chatr server
    ```

    By default, the server listens on `127.0.0.1:8080`. You can override this via flags or environment variables.

3. **Server logs**: You'll see logs about new connections, disconnections, and server status on your terminal.

### Running the Client

1. **In a new terminal**, navigate to the project directory again (or from anywhere if you have the binary in your `PATH`).

2. **Run the client**:

    ```bash
    go run main.go client
    ```

    Or using the built binary:

    ```bash
    ./chatr client
    ```

3. **Provide a username** (optional). If you want a custom username, you can either:
   - **Set it in `.env`**: `USERNAME=yourname`
   - **Use a flag**:

    ```bash
    go run main.go client --username yourname
    ```

    When the client starts, it will connect to the server and you will see a TUI displaying online users and chat messages.

### Environment Variables

Chatr supports the following environment variables (through `.env` or your system environment):

- `HOST` (default: `127.0.0.1`)
- `PORT` (default: `8080`)
- `USERNAME` (default: `Guest` when using the client)

Each environment variable can also be overridden by command-line flags:

- `--host` or `-H`
- `--port` or `-P`
- `--username` or `-u`

### TCP Settings

- **Default Host**: `127.0.0.1`  
  This means both client and server run on the localhost.
- **Default Port**: `8080`  
  If port 8080 is in use, choose another port (e.g., `--port 9090`) and ensure that both server and client match the same port.

If running on different machines, ensure that:

1. The **server** machine's firewall allows inbound connections on the chosen port.
2. The **client** can reach the **server** over the network.

### Interacting via TUI

Once the client is running, you'll see:

- A **User List** panel on the left (or top, depending on your screen).
- The **Chat Box** displaying messages from you and other users.
- An **Input Field** at the bottom where you can type messages.

**Sending Messages**:

- Type your message in the input field.
- Press **Enter** to send the message.
- Watch the conversation scroll up in the Chat Box.

**Exiting**:

- Press **Ctrl + C** in the terminal that is running the client (or server).
- Or close the terminal window.

---

## Contributing

Contributions are highly encouraged! If you have any suggestions for improvements or encounter any bugs, please feel free to open an issue or submit a pull request.

---

## License

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)  
This project is licensed under the GNU General Public License v3.0. For more information, refer to the [LICENSE](LICENSE) file.

---

## Contact

For questions, feedback, or contributions, reach out via:

- **Project Maintainer**: [Mohammad Heidari](https://github.com/itsMohammadHeidari)
- **Email**: <itsMohammadHeidari@gmail.com>

Thank you for using **Chatr**!
