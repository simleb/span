# span

[![GoDoc](https://godoc.org/github.com/simleb/span?status.svg)](http://godoc.org/github.com/simleb/span)
[![Coverage Status](https://img.shields.io/coveralls/simleb/span.svg)](https://coveralls.io/r/simleb/span)
[![Build Status](https://drone.io/github.com/simleb/span/status.png)](https://drone.io/github.com/simleb/span/latest)

Command `span` generates the cartesian product of config files.

Command `span` can be used when you want to generate multiple config files
with some variables taking different values in each generated file.
In particular, it is often useful to produce the cartesian product of
multiple variables when generating config files, but multiple variables
can also be bound such that their values change synchronously. Lastly,
some variables can contain templates that are rendered using the values
of other variables being expanded (e.g. to generate a unique output path).

Since the generated config files must each have a unique path, the output
path of each file must contain the values of the expanded variables.
Thus the output path must be provided as a template containing placeholders
for each variable to be expanded. These variables are looked for by name
and must be arrays that will be replaced by their contained elements
in each of the generated files. The syntax for placeholders is:

	{variable_name}
	{var|%03d}        # limited printf-style formatting is supported
	{section.subsection.other_var}
	{"stupid .section\n name".var}

If the output path refer to non-existing directories, they will be created.
Placeholders can also appear in directory names. The format of both the input
and output files are determined by their extension (they can differ).
Supported file formats: TOML and JSON.

When two or more variables must be expanded together, that is, not as a
cartesian product, they can be bound with `-b` (or `--bind`) followed by a
comma separated list of variable names. The variables reffered to must
be arrays of the same length. This option can be used multiple times.

String variables containing placeholders as described above can be rendered
using the values of the expanded variables. Use `-r` (or `--render`) followed
by the name of the variable. This option can be used multiple times.

## Usage

	span [-r var] [-b var1,var2[,varN ...]] output_file_template input_file

## Example

	span -r output.dir -b width,height "config/{simulation.mode}/{width}x{height}_{size|%03d}.toml" config.toml

where `config.toml` contains

	id = 42
	size = [5, 10, 200]
	width = [640, 800]
	height = [480, 600]
	[simulation]
		mode = ["normal", "crazy"]
	[output]
		dir = "data/{width}x{height}_{size|%03d}_{simulation.mode}{id}.dat"

will produce the following files:

- `config/normal/640x480_005.toml`
- `config/normal/640x480_010.toml`
- `config/normal/640x480_200.toml`
- `config/normal/800x600_005.toml`
- `config/normal/800x600_010.toml`
- `config/normal/800x600_200.toml`
- `config/crazy/640x480_005.toml`
- `config/crazy/640x480_010.toml`
- `config/crazy/640x480_200.toml`
- `config/crazy/800x600_005.toml`
- `config/crazy/800x600_010.toml`
- `config/crazy/800x600_200.toml`

And (for instance) `config/crazy/640x480_010.toml` will contain

	height = 480
	id = 42
	size = 10
	width = 640
	[output]
	  dir = "data/640x480_010_crazy42.dat"
	[simulation]
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
