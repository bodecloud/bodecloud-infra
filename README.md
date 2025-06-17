# Media Stack Docker Image

This repository contains a Docker image for running a media stack application. The image is designed to be simple to use and can be shared with others for a small open source project.

## About Configuration Files

This Docker image includes the configuration files from the `./configs` directory in the repository. These files are copied into the image at build time and will be available at `/etc/media-stack/configs` in the container.

When users run the container, they can:

- Use the included configs by not mounting anything
- Override the configs by mounting their own directory with `-v ./their-configs:/etc/media-stack/configs`

## Building the Docker Image

You can build the Docker image locally using the following command:

```bash
docker build -t media-stack:latest .
```

## Running the Container

To run the container with the configuration files included in the image:

```bash
docker run -d \
  --name media-stack \
  -p 8080:8080 \
  media-stack:latest
```

To override with your own configuration files:

```bash
docker run -d \
  --name media-stack \
  -v ./configs:/etc/media-stack/configs \
  -p 8080:8080 \
  media-stack:latest
```

This will:

- Mount your local `./configs` directory to `/etc/media-stack/configs` in the container
- Expose port 8080 from the container to your host
- Run the container in detached mode

## Environment Variables

You can customize the container by setting the following environment variables:

- `TZ`: Timezone (default: UTC)
- `PUID`: User ID (default: 1000)
- `PGID`: Group ID (default: 1000)

Example:

```bash
docker run -d \
  --name media-stack \
  -p 8080:8080 \
  -e TZ=America/New_York \
  -e PUID=1001 \
  -e PGID=1001 \
  media-stack:latest
```

## Publishing Your Image

### Option 1: Docker Hub

1. Create an account on [Docker Hub](https://hub.docker.com/)
2. Log in to Docker Hub from your terminal:

   ```bash
   docker login
   ```

3. Tag your image:

   ```bash
   docker tag media-stack:latest yourusername/media-stack:latest
   ```

4. Push the image to Docker Hub:

   ```bash
   docker push yourusername/media-stack:latest
   ```

### Option 2: GitHub Container Registry (GHCR)

1. Create a Personal Access Token (PAT) on GitHub with the `write:packages` scope
2. Log in to GHCR:

   ```bash
   echo $GITHUB_PAT | docker login ghcr.io -u USERNAME --password-stdin
   ```

3. Tag your image:

   ```bash
   docker tag media-stack:latest ghcr.io/yourusername/media-stack:latest
   ```

4. Push the image to GHCR:

   ```bash
   docker push ghcr.io/yourusername/media-stack:latest
   ```

## Using Your Published Image

Once published, others can use your image with:

```bash
# For Docker Hub - using the configs included in the image
docker run -d -p 8080:8080 yourusername/media-stack:latest

# For GHCR - using the configs included in the image
docker run -d -p 8080:8080 ghcr.io/yourusername/media-stack:latest

# To override with their own configs
docker run -d -v ./their-configs:/etc/media-stack/configs -p 8080:8080 ghcr.io/yourusername/media-stack:latest
```

## License

[Add your license information here]
