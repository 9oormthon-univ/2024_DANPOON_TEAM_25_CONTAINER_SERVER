package dockerclient

const CODESERVER_DOCKERFILE = `
# Base image
FROM nixos/nix

# Install code-server and additional programs
RUN nix-channel --add https://nixos.org/channels/nixos-unstable nixpkgs && \
    nix-channel --update && \
    nix-env -iA nixpkgs.code-server \
	%s \
	&& mkdir -p /config

# Set code-server to run on container startup
CMD ["code-server", "--bind-addr", "0.0.0.0:8080", "--auth", "none", "/config"]`
