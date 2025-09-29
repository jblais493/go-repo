// package air

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// )

// const airTemplate = `
// root = "."
// tmp_dir = "tmp"

// [build]
//   cmd = "go build -o ./bin/main .main"
//   bin = "bin/main"
//   full_bin = ""
//   include_ext = ["go", "tpl", "tmpl", "html"]
//   include_dir = ["cmd", "internal", "ui"]
//   exclude_dir = ["tmp"]

// [log]
//   time = true
// `
