# span

Command `span` generates the cartesian product of config files.

Command span can be used when you want to generate multiple config files
with some variables taking different values in each generated file.
The input is a config file where the value of those variables have been
replaced with arrays of values, and a path for the output files as a
template containing reference to these variables. See the example below.

## Usage

	span <output file template> <input file>

## Example

	span "{mode}/{id|%03d}.toml" config.toml

where `config.toml` contains

	foo = 42
	id = [5, 10, 200]
	mode = ["normal", "crazy"]

will produce the following files:

	- normal/005.toml
	- normal/010.toml
	- normal/200.toml
	- crazy/005.toml
	- crazy/010.toml
	- crazy/200.toml

And (for instance) `crazy/010.toml` will contain

	foo = 42
	id = 10
	mode = "crazy"

## License

The MIT License (MIT)

	Copyright (c) 2015 Simon Leblanc
	
	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:
	
	The above copyright notice and this permission notice shall be included in
	all copies or substantial portions of the Software.
	
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
