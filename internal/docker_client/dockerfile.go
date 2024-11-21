package dockerclient

const CODESERVER_DOCKERFILE = `
FROM nixos/nix

# Install code-server, sed, zsh and additional programs
RUN nix-channel --add https://nixos.org/channels/nixos-unstable nixpkgs && \
    nix-channel --update && \
    nix-env -iA nixpkgs.code-server \ 
             nixpkgs.zsh \
		%s \
    && mkdir -p /project && \
    echo "exec zsh" >> ~/.bashrc

# Set code-server to run on container startup
CMD ["code-server", "--bind-addr", "0.0.0.0:8080", "--auth", "none", "/project"]
	`
