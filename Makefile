all: test

test.pkg:
	@go test ./...

test.cmd:
	@make -C demos test

test: test.pkg test.cmd

# -----------------------------------------------------------------------
# A plotter base on the python matplotlib library is used for drawing
# the timeseries. Have a look in the plotter folder that is a python
# package.
plot:
	make -C plotter plot

# -----------------------------------------------------------------------
# Documentation
README.pdf: README.md
	pandoc $< -o $@

doc:
	@make -C admin/doc all

# -----------------------------------------------------------------------
# Clean the workspace
clean:
	@rm -f *~ out.* README.pdf
	@make -C plotter clean
	@make -C demos clean
	@make -C admin/doc clean


