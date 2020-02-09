# Golang Simple Realtime Chat - Websocket

## Installation

#### Clone the source

```bash
git clone https://github.com/taufikardiyan28/chat.git
```

#### Setup dependencies

```bash
go build
```

#### Configuration
Edit config/*.yaml file and change your port and mysql database configuration

#### Run the app
```
use default development configuration
go run main.go

or with custom configuration file
go run main.go --config=<your_config_file>
```

#### Simple Web Demo
###### open your browser and type http://localhost:<your_port>/demo?id=<your_user_id>

