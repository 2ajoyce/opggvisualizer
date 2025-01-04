# opggvisualizer

[op.gg](https://op.gg) is an invaluable resource for League of Legends players. Building on that, opggvisualizer enables players to visualize their op.gg statistics over a longer time, and at a more granular level. It uses a sqlite database, a grafana dashboard, a Go backend, and a cron triggered daily update.

This project would not be possible without

- [op.gg](https://op.gg)
- [ddragon](https://riot-api-libraries.readthedocs.io/en/latest/ddragon.html)

This application makes a single api call to both services every 24 hours. Please do not abuse their services by increasing that threshhold without additional caching.

### Screenshot

![alt text](docs/images/image.png)
*...I will not be accepting feedback on these stats*

### System Requirements

- Docker
- docker-compose
- Go
- Make (optional, but recommended)

### Environment Variables

Several environment variables are required. Create a `.env` file in the root directory with the following contents.

```
SUMMONER_ID=<<SUMMONER_ID>>
DATABASE_PATH=/opggvisualizer_data/data.db
GF_SECURITY_ADMIN_PASSWORD=admin
```

Replace `<<SUMMONER_ID>>` with the summoner id obtained from the op.gg http call.

Replace the admin password with a more secure password if desired. Otherwise you will be prompted to change it on initial login.

### Running the app

Build: `make build`

Initialize Data `make init`

Run: `make run`

### Grafana

The Grafana dashboard can be accessed at http://localhost:3000

The initial credentials are

- Username: `admin`
- Password: `admin`

### Refresh Cycle

By default the application will trigger an update every hour. This is set by the cron timing in `./docker-compose.yml`. This **may** trigger a data pull.

By default the application will pull fresh data once every 24 hours. When triggered, the current time is checked against a timestamp in the `fetch` database table. These timestamps are saved independently for Games and Champions. The timestamp is updated on successful fetches.

## C4 Diagrams

### Context Diagram

```mermaid
    C4Context
    %% title System Context Diagram

    Person(user, "Gamer")
    System(opggvisualizer, "opggvisualizer", "Visualize OP.GG data")

    BiRel(user, opggvisualizer, "Uses")
    UpdateRelStyle(user, opggvisualizer, $offsetY="0", $offsetX="10")
    UpdateLayoutConfig($c4ShapeInRow="1", $c4BoundaryInRow="1")
```

### Container Diagram

```mermaid
    C4Container
    %% title Container Diagram

    Person(user, "Gamer")

    Container_Boundary(containers, "Containers"){
        Container(grafana, "Grafana", "Docker Container", "Visualizes data from the SQLite database")
        Container(opggvisualizer, "opggvisualizer", "Docker Container", "Runs the Go application")
        Container(cron, "Cron", "Docker Container", "Runs Hourly")
    }

    Container_Boundary(volumes, "Volumes"){
        ContainerDb(sqlite, "SQLite", "Database", "Stores game and champion data")
    }
    Rel(user, grafana, "Uses")
    Rel(opggvisualizer, sqlite, "Reads/Writes")
    Rel(grafana, sqlite, "Reads")
    Rel(cron, opggvisualizer, "Triggers data refresh")

    UpdateLayoutConfig($c4ShapeInRow="3", $c4BoundaryInRow="1")
    UpdateRelStyle(cron, opggvisualizer, $offsetY="-15", $offsetX="-30")
    UpdateRelStyle(opggvisualizer, sqlite, $offsetY="0", $offsetX="40")
    UpdateRelStyle(grafana, sqlite, $offsetY="0", $offsetX="15")
```

### Component Diagram

```mermaid
    C4Component
    %% title Component Diagram

    Container_Boundary(opggvisualizer, "opggvisualizer", "Docker Container", "Runs the Go application") {
        Component(cli, "CLI", "Go", "Handles command-line interactions")
        Component(db, "Database", "Go", "Manages database interactions")
        Component(api, "API", "Go", "Exposes HTTP endpoints")
        Component(client, "Client", "Go", "Fetches data from external APIs")
        Component(config, "Config", "Go", "Loads configuration settings")
    }

    Container_Boundary(volumes, "Volumes"){
        ContainerDb(sqlite, "SQLite", "Database", "Stores game and champion data")
    }
    Rel(cli, db,"Wipe Data")
    Rel(cli, client,"Fetch Data")
    Rel(cli, api,"Start / Stop")
    Rel(api, client,"Refresh Data")
    Rel(client, db,"Read/Write")
    BiRel(db, sqlite, "Uses")

    UpdateLayoutConfig($c4ShapeInRow="2", $c4BoundaryInRow="2")
    UpdateRelStyle(cli, db, $offsetY="-20", $offsetX="-30")
    UpdateRelStyle(cli, client, $offsetY="-20", $offsetX="-20")
    UpdateRelStyle(cli, api, $offsetY="0", $offsetX="-70")
    UpdateRelStyle(api, client, $offsetY="20", $offsetX="-40")
    UpdateRelStyle(client, db, $offsetY="0", $offsetX="10")
    UpdateRelStyle(db, sqlite, $offsetY="-10", $offsetX="-15")
```

### Code Diagram

```mermaid
    C4Component
    %% title Code Diagram

    Boundary(cli, "CLI", "Go", "Handles command-line interactions") {
        Component(cli.go, "cli.go", "Go", "Main entry point for CLI commands")
        Boundary(cli_sub, "") {
            Component(cli_api.go, "cli / api.go", "Go", "CLI commands for API server")
            Component(cli_db.go, "cli / db.go", "Go", "CLI commands for database operations")
            Component(cli_fetch.go, "cli / fetch.go", "Go", "CLI commands for fetching data")
        }
    }

    Boundary(api, "API", "Go", "Exposes HTTP endpoints") {
        Component(api.go, "api.go", "Go", "Create, Start, Stop the server")
    }

    Boundary(client, "Client", "Go", "Fetches data from external APIs") {
        Component(client.go, "client.go", "Go", "Multi client functions")
        Boundary(client_sub, "") {
            Component(fetchGames.go, "client / fetchGames.go", "Go", "Fetches game data")
            Component(fetchChampions.go, "client / fetchChampions.go", "Go", "Fetches champion data")
        }
    }
    Boundary(db, "Database", "Go", "Manages database interactions") {
        Component(db.go, "db.go", "Go", "Multi resource database functions")
        Boundary(db_sub, "") {
            Component(games.go, "games.go", "Go", "DB functions for games and participants")
            Component(champions.go, "champions.go", "Go", "DB functions for champions")
        }
    }

    Boundary(config, "Config", "Go", "Loads configuration settings") {
        Component(config.go, "config.go", "Go", "Main entry point for configuration settings")
    }


    UpdateLayoutConfig($c4ShapeInRow="6", $c4BoundaryInRow="2")
```
