build: gorth.go
	go build

run: build
	./gorth

run_001: build
	./gorth -src gorth_source/001.gorth

test_001: build # Test - empty source file argument
	./gorth -src

test_bad_file: build # Test - bad source file
	./gorth -src very/bad/path/to/file.gorth

clean:
	rm gorth

