<<<<<<< HEAD
![Haxe logo](http://haxe.org/img/haxe2/logo.png)
# [Haxe - The Cross-Platform Toolkit](http://haxe.org)

[![Build Status](https://travis-ci.org/HaxeFoundation/haxe.png?branch=development)](https://travis-ci.org/HaxeFoundation/haxe)

Haxe is an open source toolkit that allows you to easily build cross-platform tools and applications that target many mainstream platforms. The Haxe toolkit includes:

 * **The Haxe programming language**, a modern, high-level, strictly-typed programming language
 * **The Haxe cross-compiler**, a state-of-the-art, lightning-speed compiler for many targets
 * **The Haxe standard library**, a complete, cross-platform library of common functionality

Haxe allows you to compile for the following targets:

 * C++
 * C#
 * Flash
 * Java
 * JavaScript
 * NekoVM
 * PHP

You can try Haxe directly from your browser at [try.haxe.org](http://try.haxe.org)!

For more information about Haxe, head to the [offical Haxe website](http://haxe.org).

## License

The Haxe project has several licenses, covering different parts of the projects.

 * The Haxe compiler is released under the GNU General Public License version 2 or any later version.
 * The Haxe libraries are released under a "two-clause" BSD license.
 * The Neko runtime is licensed under the GNU Lesser General Public License version 2.1 or any later version.

For the complete Haxe licenses, please see http://haxe.org/doc/license or [extra/LICENSE.txt](extra/LICENSE.txt).

## Installing Haxe

The latest stable release is [Haxe v3.1.3](http://haxe.org/download). Pre-built binaries are available for your platform:

 * **[Windows installer](http://haxe.org/file/haxe-3.1.3-win.exe)**
 * **[Windows binaries](http://haxe.org/file/haxe-3.1.3-win.zip)**
 * **[OSX installer](http://haxe.org/file/haxe-3.1.3-osx-installer.pkg)**
 * **[OSX binaries](http://haxe.org/file/haxe-3.1.3-osx.tar.gz)**
 * **[Linux 32-bit binaries](http://haxe.org/file/haxe-3.1.3-linux32.tar.gz)**
 * **[Linux 64-bit binaries](http://haxe.org/file/haxe-3.1.3-linux64.tar.gz)**

Automated development builds are available from [build.haxe.org](http://build.haxe.org).

## Building from source

 1. Clone the repository using git. Be sure to initialize and fetch the submodules.

        git clone --recursive git://github.com/HaxeFoundation/haxe.git
        cd haxe

 2. Follow the [documentation on building Haxe for your platform](http://haxe.org/doc/build).

## Using Haxe

For information on on using Haxe, consult the [Haxe documentation](http://haxe.org/doc):

 * [Haxe introduction](http://haxe.org/doc/intro), an introduction to the Haxe toolkit
 * [Haxe language reference](http://haxe.org/ref), an overview of the Haxe programming language
 * [Haxe API](http://api.haxe.org/), a reference for the Haxe standard and native APIs
 * [Haxelib](http://lib.haxe.org/), a repository of Haxe libraries for a variety of needs

## Community

You can get help and talk with fellow Haxers from around the world via:

 * the [official Haxe Google Group](https://groups.google.com/forum/#!forum/haxelang)
 * the [Haxe IRC chatroom](http://unic0rn.github.io/tiramisu/haxe/), #haxe on chat.freenode.net

## Version compatibility

Haxe   | neko
----   | -----
2.*    | 1.*
3.0.0  | 2.0.0
3.1.3  | 2.0.0
=======
# TARDIS Go -> Haxe transpiler

#### Haxe -> JavaScript / ActionScript / Java / C++ / C# / PHP / Neko

[![Build Status](https://travis-ci.org/tardisgo/tardisgo.png?branch=master)](https://travis-ci.org/tardisgo/tardisgo)
[![GoDoc](https://godoc.org/github.com/tardisgo/tardisgo?status.png)](https://godoc.org/github.com/tardisgo/tardisgo)
[![status](https://sourcegraph.com/api/repos/github.com/tardisgo/tardisgo/badges/status.png)](https://sourcegraph.com/github.com/tardisgo/tardisgo)

## Objectives:
The objective of this project is to enable the same [Go](http://golang.org) code to be re-deployed in  as many different execution environments as possible, thus saving development time and effort. 
The long-term vision is to provide a framework that makes it easy to target many languages as part of this project.

The first language targeted is [Haxe](http://haxe.org), because the Haxe compiler generates 7 other languages and is already well-proven for making multi-platform client-side applications, mostly games. 

Planned current use cases: 
- For the Go community: write a library in Go and call it from  existing Haxe, JavaScript, ActionScript, Java, C++, C# or PHP applications. 
- For the Haxe community: provide access to the portable elements of Go's extensive libraries and open-source code base.
- Write a multi-platform client-side application in a mixture of Go and Haxe, using [OpenFL](http://openfl.org) / [Lime](https://github.com/openfl/lime) or [Kha](http://kha.ktxsoftware.com/) to target a sub-set of: 
Windows,
Mac,
Linux,
iOS,
Android,
BlackBerry,
Tizen,
Emscripten,
HTML5,
webOS,
Flash,
Xbox and PlayStation.

For more background and on-line examples see the links from: http://tardisgo.github.io/

## Project status: a working proof of concept
####  DEMONSTRABLE, EXPERIMENTAL, INCOMPLETE,  UN-OPTIMIZED

> "Premature optimization is the root of all evil (or at least most of it) in programming." - Donald Knuth

All of the core [Go language specification](http://golang.org/ref/spec) is implemented, including single-threaded goroutines and channels. However the packages "unsafe" and "reflect", which are mentioned in the core specification, are not currently supported. 

Goroutines are implemented as co-operatively scheduled co-routines. Other goroutines are automatically scheduled every time there is a function call or a channel operation, so loops without calls or channel operations will never give up control. The empty function tardisgolib.Gosched() provides a way to give up control without including the full Go runtime.  

Some parts of the Go standard library work, as you can see in the [example TARDIS Go code](http://github.com/tardisgo/tardisgo-samples), but the bulk has not been  tested or implemented yet. If the standard package is not mentioned in the notes below, please assume it does not work. So fmt.Println("Hello world!") will not transpile, instead use the go builtin function: println("Hello world!").  

Some standard Go library packages do not call any runtime C or assembler functions and will probably work OK (though their tests still need to be rewritten and run to validate their correctness), these include:
- errors
- unicode
- unicode/utf8 
- unicode/utf16
- sort
- container/heap
- container/list
- container/ring

Other standard libray packages make limited use of runtime C or assembler functions without using the actual Go "runtime" or "os" packages. These limited runtime functions have been emulated for the following packages (though their tests still need to be rewritten and run to validate their correctness). To use these packages, their corresponding runtime functions need to be included as follows:
```
include ( 
	"bytes" 
	_ "github.com/tardisgo/tardisgo/golibruntime/bytes"
	
	"strings"  // but see issue #19 re JS and Flash
	_ "github.com/tardisgo/tardisgo/golibruntime/strings"
	
	"sync" // partial: only WaitGroup is known to work 
	_ "github.com/tardisgo/tardisgo/golibruntime/sync"
	
	"sync/atomic""
	_ "github.com/tardisgo/tardisgo/golibruntime/sync/atomic"
	
	"math"
	"strconv"  // uses the math package
	_ "github.com/tardisgo/tardisgo/golibruntime/math"
)
```
At present, standard library packages which rely on the Go "runtime", "os", "reflect" or "unsafe" packages are not implemented (although some OSX test code is in the golibruntime tree).

A start has been made on the automated integration with Haxe libraries, but this is currently incomplete see: https://github.com/tardisgo/gohaxelib

TARDIS Go specific runtime functions are available in [tardisgolib](https://github.com/tardisgo/tardisgo/tree/master/tardisgolib):
```
import "github.com/tardisgo/tardisgo/tardisgolib" // runtime functions for TARDIS Go
```

The code is developed on OS X Mavericks‎, using Go 1.2.1 and Haxe 3.1.1. The target platforms tested are Ubuntu 13.10 32-bit & 64-bit, and Windows 7 32-bit. The 64-bit platforms work fine, but compilation to the C# target fails on Win-7; and PHP is flakey, espeecially the 32-bit version (but you probably knew that).

## Installation and use:
 
Dependencies:
```
go get code.google.com/p/go.tools
```

TARDIS Go:
```
go get github.com/tardisgo/tardisgo
```

If tardisgo is not installing and there is a green "build:passing" icon at the top of this page, please e-mail [Elliott](https://github.com/elliott5)!

To translate Go to Haxe, from the directory containing your .go files type the command line: 
```
tardisgo yourfilename.go 
``` 
A single Go.hx file will be created in the tardis subdirectory.

To run your transpiled code you will first need to install [Haxe](http://haxe.org).

Then to run the tardis/Go.hx file generated above, type the command line: 
```
haxe -main tardis.Go --interp
```
... or whatever [Haxe compilation options](http://haxe.org/doc/compiler) you want to use. 
See the [tgoall.sh](https://github.com/tardisgo/tardisgo-samples/blob/master/scripts/tgoall.sh) script for simple examples.

To run cross-target command-line tests as quickly as possible, the "-testall" flag  concurrently runs the Haxe compiler and executes the resulting code for all supported targets (with compiler output suppressed and results appearing in the order they complete):
```
tardisgo -testall myprogram.go
```

PHP specific issues:
* to compile for PHP you currently need to add the haxe compilation option "--php-prefix tgo" to avoid name conflicts
* very long PHP class/file names may cause name resolution problems on some platforms

## Next steps:
Please go to http://github.com/tardisgo/tardisgo-samples for example Go code modified to work with tardisgo.

For public help or discussion please go to the [Google Group](https://groups.google.com/d/forum/tardisgo); or feel free to e-mail [Elliott](https://github.com/elliott5) direct to discuss any issues if you prefer.

The documentation is sparse at present, if there is some aspect of the system that you want to know more about, please let [Elliott](https://github.com/elliott5) know and he will prioritise that area.

If you transpile your own code using TARDIS Go, please report the bugs that you find here, so that they can be fixed.

## Future plans:

Development priorities:
- For all Go standard libraries, report testing and implementation status
- Improve integration with Haxe code and libraries, automating as far as possible
- Improve currently poor execution speeds and update benchmarking results
- Research and publish the best methods to use TARDIS Go to create multi-platform client-side applications
- Improve debug and profiling capabilities
- Add command line flags to control options
- Publish more explanation and documentation
- Move more of the runtime into Go (rather than Haxe) to make it more portable 
- Implement other target languages

If you would like to get involved in helping the project to advance, that would be wonderful. However, please contact [Elliott](https://github.com/elliott5) or discuss your plans in the [tardisgo](https://groups.google.com/d/forum/tardisgo) forum before writing any substantial amounts of code to avoid any conflicts. 

## License:

MIT license, please see the license file.
>>>>>>> upstream/master
