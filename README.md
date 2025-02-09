# typi

Typi fetches youtube videos with predefined keyword (cricket) and displays them in the dashboard.

> [!note]
> I have hosted the sample application at [typi.infosec.org.in](https://typi.infosec.org.in/). Please check it out.
>
> sample config:
> - API Keys - 3
> - frequency - 5 minutes (minimum time to allow quotas to reset for one day with 3 keys)
> - topic - cricket

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Components](#components)
    - [typi-database](#typi-database)
    - [typi-api](#typi-api)
    - [typi-cron](#typi-cron)
    - [typi-ui](#typi-ui)
    - [typi-proxy](#typi-proxy)

## Prerequisites
I have used docker compose to deploy the application. Please follow [this](https://docs.docker.com/engine/install/ubuntu/) to install docker engine on your machine.

## Installation

Clone the repository

```bash
git clone https://github.com/TheAlpha16/typi.git
cd yourproject
```

You can use docker compose to deploy the application. 

```bash
docker compose up -d
```

The application binds to port 80. Please make sure that it is not in use. If you want to change the port, you can do so by changing the bind port for proxy in the docker-compose.yml file.

```yaml
typi-proxy:
build:
    context: ./proxy
container_name: typi-proxy
ports:
    - "8888:80" # Change this to your desired port (8888)
```

## Components

The application has 4 components:

### typi-database

- Postgres database to store the fetched video data and required configuration. It has 2 tables -

```sql
CREATE TABLE IF NOT EXISTS videos (
    vid text PRIMARY KEY,
    title text,
    description text,
    thumbnail text,
    published_at text
);

CREATE TABLE IF NOT EXISTS config (
    key text PRIMARY KEY,
    value text
);
```

### typi-api

- A go fiber application serving the HTTP API to fetch data stored from the database. It has 2 endpoints -

1. `GET /api/ping` - Health check endpoint.

2. ` GET /api/videos` - Fetches videos from the database.

#### parameters

- `page` - Page number to fetch the data. Default is 1.
- `per_page` - Number of videos per page. Default is 10.

#### response

This is a typical response from the API.

```json
{
    "page": 1, // current page number
    "per_page": 1, // number of videos per page
    "status": "success",
    "total_pages": 13, // total number of pages
    "total_videos": 50, // total number of videos
    "videos": [
        {
            "vid": "9DSf0ot_zTs", // unique video id
            "title": "India vs England, 2nd ODI | Live Cricket Match Today | IND vs ENG Live Match | India England Match",
            "description": "India vs England Live #livestream #cricket #trending #indvseng #indiavsenglandlive #livestream #cricketlive #livescore #cricket ...",
            "published_at": "2025-02-09T13:12:53Z",
            "thumbnail": "https://i.ytimg.com/vi/9DSf0ot_zTs/default.jpg"
        }
    ]
}
```

- [logrus](https://github.com/sirupsen/logrus) in combination with [lumberjack](https://github.com/natefinch/lumberjack) handle logging in the application. 

- Database queries are handled by [gorm](https://gorm.io/).

```go
package models

type Video struct {
	VID         string `gorm:"primaryKey;column:vid" json:"vid"`
	Title       string `gorm:"column:title" json:"title"`
	Description string `gorm:"column:description" json:"description"`
	PublishedAt string `gorm:"column:published_at" json:"published_at"`
	Thumbnail   string `gorm:"column:thumbnail" json:"thumbnail"`
}
```

### typi-cron

- A go application that fetches the videos from youtube with the predefined keyword (cricket) and stores them in the database. 

- It uses [cron](https://github.com/robfig/cron) to schedule the job every n seconds.

frequency of the fetch can be set with `FREQUENCY_IN_SECONDS` (defaults 300) environment variable.

```go
package config

...
var FREQUENCY_IN_SECONDS = os.Getenv("FREQUENCY_IN_SECONDS")

func init() {
	if FREQUENCY_IN_SECONDS == "" {
		FREQUENCY_IN_SECONDS = "300"
	}
}
```

```yaml
typi-cron:
build:
    context: ./cron
container_name: typi-cron
environment:
    ...
    - FREQUENCY_IN_SECONDS=60 # Change this to your desired frequency
```

- It uses [youtube v3 api](https://developers.google.com/youtube/v3) to fetch the videos.

```go
package fetcher

import (
    ...
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	client *youtube.Service
)
```

- fetcher requires `YT_API_KEYS` env var to be set. It accepts comma seperated API keys. You can get the keys from [here](https://developers.google.com/youtube/registering_an_application).

> [!note]
> Keys present in the file are just placeholders. I have attached the actual keys in the email.

```yaml
typi-cron:
build:
    context: ./cron
container_name: typi-cron
environment:
    ...
    - YT_API_KEYS=AIzaSyDOcTfdT1Z6QWrp8qqDjRl_c1sTabpfaiw,AIzaSyA79nnWhz3RnoPQaxvfcs8sb4xPN8nqOu4
```

- API keys are automatically rotated if the quota is exceeded.

```go
package keyrings

...

var (
	apiKeys  []string
	keyIndex int = 0
	mu       sync.Mutex
)

func GetKey() string {
	mu.Lock()
	defer mu.Unlock()

	currentKey := apiKeys[keyIndex]
	keyIndex = (keyIndex + 1) % len(apiKeys)
	return currentKey
}
```

### typi-ui

- A NextJS application that displays the fetched videos in the dashboard.

- It uses [tailwindcss](https://tailwindcss.com/) for styling.

- videos are fetched from the API at `GET /api/videos`. 

- Zustand is used for state management.

```ts
...

interface VideoState {
    videos: { [page: number]: Video[] };
    currentPage: number;
    totalPages: number;
    isLoading: boolean;
    searchQuery: string;
    filteredVideos: Video[];
    fetchVideos: (page: number, perPage: number) => Promise<void>;
    setSearchQuery: (query: string) => void;
}

export const useVideoStore = create<VideoState>()(
    persist(
        (set, get) => ({
            ...
        }),
        {
            name: "video-storage",
        }
    )
);

const filterVideos = (videos: { [page: number]: Video[] }, query: string): Video[] => {
    ...
};

const flattenVideos = (videos: { [page: number]: Video[] }): Video[] => {
    ...
};
```

- It has a search bar to filter the videos based on the specified keyword.

- UI lazy loads the videos as the user scrolls down the page.

### typi-proxy

- A nginx proxy to route the requests to the respective services.

```nginx
server {
    listen 80;
    server_name localhost;

    location /api {
        proxy_pass http://typi-api:4601;
    }

    location / {
        proxy_pass http://typi-ui:4602;
    }
}
```
