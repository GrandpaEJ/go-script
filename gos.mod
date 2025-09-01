module main

go 1.21

# Go-Script module configuration
gos_version "1.0.0"

# Dependencies
require (
    # Standard Go modules work automatically
)

# Go-Script specific configuration
config {
    default_package "main"
    output_dir "./generated"
    module_paths ["./modules", "./lib"]
}
