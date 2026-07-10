package main

// run_context is meant to store all the cli args and flags into a single object
type RunContext struct {
	// args passed to the cli, program name excluded
	Args []string

	// format used for the run
	Format string

	// weither the version flag is present
	Version bool

	// weither the help flag is
	Help bool
}