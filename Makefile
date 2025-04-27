# Makefile

# VARIABLES -------------------------------------

server_cmd = "cmd/server/main.go"
server_out = "fileup_srv"

out_dir = ./go-out


# SETUP -----------------------------------------

# Create out dir if not exists
$(out_dir):
	mkdir -p $(out_dir)


# SERVER -----------------------------------------

# Build fileup server
b-server: $(out_dir)
	go build -o $(out_dir)/$(server_out) $(server_cmd)
	chmod +x $(out_dir)/$(server_out)

# Run fileup server
r-server: b-server
	$(out_dir)/$(server_out)


# TESTING ----------------------------------------

test:
	go test -vet=all -parallel=4 ./lib/*/** -v
 
# ALL --------------------------------------------

all: test b-server

