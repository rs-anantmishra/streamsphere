# :film_strip: Streamsphere

It is a media library. It supports channels, playlists and videos. <br />
The primary goal of this application is to provide complete functionality of a media library, while being as light weight as possible.
> *Note: The tool is under active development.*

## :computer: Interface
![downloading](https://github.com/user-attachments/assets/8c9654aa-6231-4bde-b144-c79d9b233592)

## :rocket: Getting Started

### Features
- 📺 Download Channels from supported domains
- 📼 Downlad & update playlists from supported domains
- 📽️ Download Videos from supported domains
- 🔍 Search and play videos by title
- 👾 UI to navigate your media library
- 📥 Download media content that has been added to streamsphere through browser
- ✨ View tags, categories, size of media files and other details for the downloaded content
- 🎴 Light & Dark theme support

### :whale: Deploy using docker compose 
You will need [Docker](https://docs.docker.com/get-docker/) installed on your system for this.

```
services:
  streamsphere:
    image: streamsphere/streamsphere:latest
    # image: streamsphere/streamsphere:latest-arm
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
    image: streamsphere/streamsphere-content:latest
    # image: streamsphere/streamsphere-content:latest-arm
    container_name: streamsphere-content
    restart: unless-stopped
    ports:
      - 1288:3500
    volumes:
      - content-data:/content
volumes:
  db-data:
  content-data:
```
> *Note: If you are using an ARM64 machine (like a raspberry-pi), use the commented image name with **latest-arm** tag (the default uncommented image with **latest** tag is for AMD64 machines only).*

### Environment Variables

- **HOST_IP** is the IP (or hostname) of the machine on the network you want to host it on.
- **APPLICATION_PORT** is the port you are mapping your container services to.
- **CONTENT_PORT** is the port you are mapping your content files to.

### Supported Domains
-  youtube.com

## 〰️ Performance & Limitations
The application uses yt-dlp to fetch video metadata and thumbnail, if this is successful, it attempts to download the video and shows download metrics on the UI. Since the application first fetches the video metadata, the Get UI operation may feel slow as it takes a while to fetch the metadata and download the video thumbnail, even more so when its a playlist or a channel. The application does have the same limitations basically as yt-dlp as of the latest version. e.g. downloading a mix will end up in a never ending loop which keeps fetching data from yt-dlp.

## :hammer_and_wrench: Support & Compatibility

#### CPU Architecture
| AMD64 | ARM64 |
| ------------- | ------------- |
| ✔️ Supported | ✔️ Supported |

#### Screen Size
| 10" & below | 10" - 14" | 14" - 27" | 27" & above |
| ------------- | ------------- | ------------- | ------------- |
| ✖️ Unsupported | ✔️ Supported (untested) | ✔️ Supported | ✔️ Supported (untested) |

#### Browsers
| Chrome | Edge | Safari | Firefox |
| ------------- | ------------- | ------------- | ------------- |
|  ✔️ Supported | ✔️ Supported | ✔️ Supported | ❗ Supported (has open issues) |

## 🌟 Upcoming Features
Please refer this list of [upcoming work items](https://github.com/users/rs-anantmishra/projects/5) and please report any bugs if you find one!

## 🔼 Updating yt-dlp
The latest version of streamspehre (v0.1.15) includes an update to upgrade yt-dlp each time the container is stopped and started.
It can also be manually updated by following the below steps:
```
# login into the container
sudo docker exec -it streamsphere bash
```
Once inside the docker container:
```
# change directory to reach the yt-dlp binary
cd /app/utils/

# run the yt-dlp update
./yt-dlp_linux -U
```

Upcoming versions of streamsphere are planned to have a ui to update yt-dlp.

### 🎯 Design Update (Proposed)
This update will support: 
- scheduling playlist/channel updates
- enable extraction of extremely huge channels/playlists
- enable modular implementation of download filters
- improved failure management
- does not require additional docker builds

![Application Design](https://github.com/user-attachments/assets/7800f70f-f902-4cef-9c75-8e2664666cbe)


## ❓ Help & Support
Please feel free to report any [bugs](https://github.com/users/rs-anantmishra/projects/5) that you may have observed!

## 📝 License
GNU Affero General Public License v3.0
