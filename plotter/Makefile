checkmpl=$(shell python3 -c "import matplotlib.pyplot" 2>/dev/null; echo $$?)

all: plot

plot:
ifneq (${checkmpl}, 0)
	$(error matplotlib not found. The test can not be executed)
endif
	@./plot.py

clean:
	@rm -f *~ out.* *.pyc
	@rm -rf __pycache__

