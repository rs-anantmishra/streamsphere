FROM node:lts-bookworm-slim AS uibase

# update system and install required packages
RUN <<EOF
apt update -y
apt install -y wget sqlite3
npm install -g @angular/cli
EOF

# copy code to build
WORKDIR /app-code/
COPY . .

# build angular ui
RUN <<EOF
cd ui
npm install
ng build --configuration production 
EOF

# copy ui build to release directory
WORKDIR /app/browser
RUN cp /app-code/ui/dist/ui/browser/ /app/ -rv

# get latest yt-dlp
WORKDIR /app/utils
RUN <<EOF
wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux
chmod +x /app/utils/yt-dlp_linux
EOF

# make directories and copy required files
WORKDIR /app
RUN <<EOF
mkdir ./database/scripts ./database/db -p
cp /app-code/database/scripts/* /app/database/scripts/
cp /app-code/.env /app/.env
cp /app-code/entrypoint.sh /app/entrypoint.sh
mkdir application content
EOF

# init database
WORKDIR /app/database/db
RUN <<EOF
sqlite3 streamsphere.db < ../scripts/create.sql | bash
sqlite3 streamsphere.db < ../scripts/seed.sql | bash
EOF

# stage 2 starts here for application build 
FROM golang:1.22.8-bookworm AS appbase

COPY --from=uibase /app-code/ /app-code/
COPY --from=uibase /app/ /app/

WORKDIR /app-code
RUN <<EOF
cd cmd
go build
cp cmd /app/application/streamsphere -rv
chmod +x /app/application/streamsphere
EOF

# stage 3 copy application to bookworm lts as final
FROM debian:bookworm-slim AS final
RUN <<EOF
apt update -y && apt upgrade -y
apt install -y gettext-base
EOF

COPY --from=appbase /app/ /app/

EXPOSE 3000
# make entrypoint script executable
WORKDIR /app/
RUN <<EOF
chmod +x /app/entrypoint.sh | bash 
EOF

# CMD ["/bin/bash"]
CMD ["./entrypoint.sh"]
