#!/bin/bash

# Install Go on macOS
echo "Installing Go..."

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

# Install Go
brew install go

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc

# Reload shell
source ~/.zshrc

# Verify installation
go version

echo "Go installation complete!"

