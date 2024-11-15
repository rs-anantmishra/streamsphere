# :film_strip: Streamsphere

It is a media library. At the moment you can add videos to it from yt. It supports channels, playlists and videos from yt.
> *Note: The tool is under active development. Please see open tickets for upcoming features and bug fixes.* 

## :computer: Interface
![downloading](https://github.com/user-attachments/assets/8c9654aa-6231-4bde-b144-c79d9b233592)

## :rocket: Getting Started

#### Deploy using docker compose 🐳

```
services:
  streamsphere:
    image: streamsphere/streamsphere
    container_name: streamsphere
    restart: unless-stopped
    ports:
      - 1282:3000
    environment:
      HOST_IP: "192.168.1.15"
      APPLICATION_PORT: "1282"
      CONTENT_PORT: "1288"
    volumes:
      - db-data:/app/database/db
      - content-data:/app/content

  content:
    image: streamsphere/streamsphere-content
    ports:
      - 1288:3500
    volumes:
      - content-data:/content

volumes:
  db-data:
  content-data:
```

#### Environment Variables

**HOST_IP** is the IP of you machine on the network you want to host it on.
**APPLICATION_PORT** is the port you are mapping your container services to.
**CONTENT_PORT** is the port you are mapping your content files to.

## :hammer_and_wrench: Compatibility
| CPU Architecture  | Supported |
| ------------- | ------------- |
| AMD64 | ✔️ Supported |
| ARM64 | ❗ Support to be added soon |

| Screen Size | Supported |
| ------------- | ------------- |
| 10" & below | ✖️ Unsupported |
| 10" - 14" | ✔️ Supported (untested) |
| 14" - 27" | ✔️ Supported |
| 27" & above | ✔️ Supported (untested) |

| Browsers | Supported |
| ------------- | ------------- |
| Chrome | ✔️ Supported |
| Edge | ✔️ Supported |
| Safari | ✔️ Supported |
| Firefox | ❗ Supported (has open issues) |
