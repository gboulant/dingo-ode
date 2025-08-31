all: talk01.html

#OPTIONS=--self-contained --mathml
OPTIONS=--standalone --mathjax

RST2HTML=pandoc ${OPTIONS} --css=./css/pandoc.css -t slidy

srcfiles=talk01.rst

talk01.html: ${srcfiles} css/pandoc.css
	@$(RST2HTML) ${srcfiles} -o $@

clean:
	@rm -f *.html
	@find . -name "*~" | xargs rm -f
