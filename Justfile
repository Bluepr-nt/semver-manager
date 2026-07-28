# Build automation for semver-manager

binary := "smgr"
src_dir := "src"
cmd_pkg := "./cmd/smgr"

# Build the CLI binary at the repository root
build:
    cd {{src_dir}} && go build -o ../{{binary}} {{cmd_pkg}}
