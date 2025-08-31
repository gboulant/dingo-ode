all: test

build:
	go build

test: build
	./demos -d spring01
	./demos -d spring02
	./demos -d spring03
	./demos -d lorenz
	./demos -d laser01
	./demos -d laser02
	./demos -d watertank
	./demos -d volterra

test.plot: build
	PYTHONPATH=${PYTHONPATH}:.. ./demos -d laser01 -p

clean:
	rm -f ./demos *~ out.*
