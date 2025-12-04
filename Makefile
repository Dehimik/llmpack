BINARY_NAME=llmpack
INSTALL_PATH=/usr/local/bin
ENTRY_POINT=cmd/llmpack/main.go

.PHONY: all build install clean uninstall

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(ENTRY_POINT)

install: build
	@echo "Installing to $(INSTALL_PATH)..."
	sudo mv $(BINARY_NAME) $(INSTALL_PATH)
	@echo "Installed! Run '$(BINARY_NAME) --help'"

uninstall:
	@echo "Uninstalling..."
	sudo rm $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstalled."

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f context.xml context.zip