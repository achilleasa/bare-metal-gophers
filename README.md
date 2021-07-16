# Bare Metal Gophers

This repository contains the demo code for the GolangUK2017 talk 
"Bare Metal Gophers: Can you write a kernel in Go?".

The slides for the talk are available [here](https://speakerdeck.com/achilleasa/bare-metal-gophers-can-you-write-an-os-kernel-in-go).

## Requirements 

The project Makefile contains targets for building a bootable ISO containing
the compiled code. [Virtualbox](https://www.virtualbox.org/) or [qemu](http://www.qemu-project.org/)
can be used to boot the generated ISO (`make run-qemu` or `make run-vbox`). A 
debugging target is also provided (`make gdb`) that boots the ISO using qemu 
and then spawns an interactive gdb session. 

The code was developed using Go 1.8. However, you should be able to compile it
with any Go version that supports native cross-compilation.

## Building on Linux

To compile the demo code when running on Linux you need a fairly recent version of:

- xorriso
- grub
- nasm 
- gcc (for GNU ld)
- go (1.6+; recommended: 1.8)
- gdb (optional; only install if you want to debug the demo code)

The above dependencies can be installed using the appropriate package manager 
for each particular Linux distribution.

## Building on OSX

To properly link the kernel object files so that the bootloader can pick up the 
multi-boot signature we need to be able to control the linker configuration. For
the time being this is only possible when using GNU ld ([lld](https://lld.llvm.org/) 
is a potential alternative but doesn't yet fully support linker scripts).

You can still build the kernel using [vagrant](https://www.vagrantup.com/). For
this purpose, a Vagrantfile is provided so all you need to do is just install 
vagrant on your machine and run `vagrant up` before running any of the following 
make commands.

To debug the demo code, you need to install gdb. If using
[brew](https://brew.sh/) you can do this by running `brew install gdb
--with-all-targets`. This is an optional step and is not required for running
the demo code. 

## Supported make targets 

The provided Makefile will work on both Linux and OSX (using vagrant) targets.
When running under OSX, the Makefile will ensure that all build-related
commands actually run inside the vagrant box. The following targets are
supported:

- `run-vbox`: compile the code, build an ISO file and run it using virtualbox.
- `run-qemu`: compile the code, build an ISO file and run it using qemu.
- `gdb`: compile the code, build an ISO file, run it using qemu and start a gdb session.
- `iso`: compile the code and build a bootable ISO using grub as the bootloader.
- `kernel`: compile the code into an elf binary.

## How does it look?

If you are curious what the demo looks like but do not want to actually build
and run it, here is a screen capture of it running under virtualbox:

![running under virtualbox](img/bare-metal-gophers.gif)
