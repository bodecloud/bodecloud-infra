FROM alpine

# Install necessary packages, including Docker and Docker Compose
RUN apk add --no-cache \
    bash \
    ca-certificates \
    curl \
    docker \
    docker-compose \
    tzdata

# Create directory for configs
RUN mkdir -p /etc/media-stack/arbitrary-scripts
RUN mkdir -p /etc/media-stack/configs
RUN mkdir -p /etc/media-stack/data/media/anime
RUN mkdir -p /etc/media-stack/data/media/books
RUN mkdir -p /etc/media-stack/data/media/comics
RUN mkdir -p /etc/media-stack/data/media/games
RUN mkdir -p /etc/media-stack/data/media/movies
RUN mkdir -p /etc/media-stack/data/media/music/videos
RUN mkdir -p /etc/media-stack/data/media/other
RUN mkdir -p /etc/media-stack/data/media/special
RUN mkdir -p /etc/media-stack/data/media/tv
RUN mkdir -p /etc/media-stack/data/watch
RUN mkdir -p /etc/media-stack/downloads
RUN mkdir -p /etc/media-stack/downloads/incomplete
RUN mkdir -p /etc/media-stack/secrets
RUN echo 'Your Docker services/containers configuration/runtime files will be stored here. Map your docker-compose.yml volumes here for persistent data storage across container restarts and rebuilds.' > /etc/media-stack/configs/about-this-folder.txt
RUN echo 'Place any custom scripts you want to run within the container in this directory. Ensure they are executable (chmod +x).' > /etc/media-stack/arbitrary-scripts/about-this-folder.txt
RUN echo 'Anime media files should be placed in this directory. Refer to docker-compose.yml for volume mapping.' > /etc/media-stack/data/media/anime/about-this-folder.txt
RUN echo 'Book media files (ebooks, audiobooks) should be placed in this directory. Refer to docker-compose.yml for volume mapping.' > /etc/media-stack/data/media/books/about-this-folder.txt
RUN echo 'Comic media files should be placed in this directory. Refer to docker-compose.yml for volume mapping.' > /etc/media-stack/data/media/comics/about-this-folder.txt
RUN echo 'Game files should be placed in this directory.  This is not currently used in docker-compose.yml, but is here for future use.' > /etc/media-stack/data/media/games/about-this-folder.txt
RUN echo 'Movie media files should be placed in this directory. Refer to docker-compose.yml for volume mapping.' > /etc/media-stack/data/media/movies/about-this-folder.txt
RUN echo 'Music video files should be placed in this directory. Refer to docker-compose.yml for volume mapping.' > /etc/media-stack/data/media/music/videos/about-this-folder.txt
RUN echo 'Other media files should be placed in this directory. This is not currently used in docker-compose.yml, but is here for future use.' > /etc/media-stack/data/media/other/about-this-folder.txt
RUN echo 'Special media files should be placed in this directory. This is not currently used in docker-compose.yml, but is here for future use.' > /etc/media-stack/data/media/special/about-this-folder.txt
RUN echo 'TV show media files should be placed in this directory. Refer to docker-compose.yml for volume mapping.' > /etc/media-stack/data/media/tv/about-this-folder.txt
RUN echo 'Files to be watched by services like Radarr, Sonarr, etc. should be placed in this directory. This is not currently used in docker-compose.yml, but is here for future use.' > /etc/media-stack/data/watch/about-this-folder.txt
RUN echo 'Downloaded files will be placed in this directory. Refer to docker-compose.yml for volume mapping.' > /etc/media-stack/downloads/about-this-folder.txt
RUN echo 'Incomplete downloads will be placed in this directory. Refer to docker-compose.yml for volume mapping.' > /etc/media-stack/downloads/incomplete/about-this-folder.txt
RUN echo 'Secret files (e.g., API keys, passwords) should be placed in this directory and referenced in your docker-compose.yml using Docker secrets.  These are not mapped as volumes, but are available to be used as docker secrets.' > /etc/media-stack/secrets/about-this-folder.txt

# Copy configs directory to the image
COPY configs /etc/media-stack/configs
COPY arbitrary-scripts /etc/media-stack/arbitrary-scripts
COPY docker-compose.yml /etc/media-stack/docker-compose.yml

# Create a non-root user and group to run the application
# Create docker group with ID 112. This is referenced in docker-compose.yml with PGID
RUN set -e && \
    (grep -q "^docker:" /etc/group || addgroup -g 112 docker) && \
    (grep -q "^mediastack:" /etc/group || addgroup -g 1000 mediastack) && \
    (grep -q "^mediastack:" /etc/passwd || adduser -D -u 1000 -G mediastack mediastack && addgroup mediastack docker) && \
    chown -R mediastack:mediastack /etc/media-stack || true

# Set environment variables
ENV TZ=UTC \
    PUID=1000 \
    PGID=112

# Switch to non-root user
USER mediastack

# Set volume for configs
VOLUME ["/etc/media-stack/configs"]
VOLUME ["/etc/media-stack/arbitrary-scripts"]
VOLUME ["/etc/media-stack/data"]
VOLUME ["/etc/media-stack/downloads"]
VOLUME ["/etc/media-stack/secrets"]

WORKDIR /etc/media-stack
ENTRYPOINT ["docker", "compose", "-f", "/etc/media-stack/docker-compose.yml", "up", "-d", "--build", "--remove-orphans"]
