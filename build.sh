# Navigate to cmd directory
# cd cmd/app

# Build the project
go build -o translate-tool

# Move the binary to the go bin directory
mv translate-tool /Users/pulsar/go/bin

# Update permissions
chmod +x /Users/pulsar/go/bin/translate-tool

# Refresh zsh
source ~/.zshrc

# Inform the user that the build is complete
echo "Build complete. Run 'translate-tool -h' to use the app."

# Navigate back to the root directory
# cd ../../
