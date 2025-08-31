all: talk pdf

pdf: userguide.pdf

# Possible highlight-style: pygments (the default), kate, monochrome,
# espresso, zenburn, haddock, and tango

userguide.pdf: userguide.rst
	pandoc --highlight-style=tango $< -o $@

view: userguide.pdf
	evince userguide.pdf

talk:
	@make -C talks

clean:
	@rm -f *~ userguide.pdf
	@make -C talks clean
