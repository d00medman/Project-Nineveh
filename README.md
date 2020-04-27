# Project Nineveh

The goal of this project is to do to the Nintendo Entertainment System what the Arcade Learning Environment ([ALE](https://github.com/mgbellemare/Arcade-Learning-Environment))
did for the Atari; make it viable to write agents to play classic Nintendo games.
 
Nineveh is a Modification of the [Nintengo](https://github.com/nwidger/nintengo) emulator (built on the same commit history, with many of the original files included). The goals I have been working towards are:
-  Emulation core uncoupled from rendering and sound generation modules for fast emulation with minimal library dependencies.
-  Python development is supported through ctypes.

The system is designed with the intent of easy interfacing with the OpenAI [Gym](https://github.com/openai/gym). Unlike the ALE, which appears to have initially been desinged
to integrate with C++ agents out of the box. Further extension would be needed to enable its control via pipes (which seems to be something of a pain to implement)

At present, can either be launched in human playable mode or in headless mode. Human playable mode is for debugging purposes, and will likely end up stripped out of the final product to avoid legal trouble. All agents will operate in headless mode. Each observation outputs a jpeg file; the option to output the entire
run to the current point as a gif is also present.

I named it Nineveh as a slant portmanteau of **NIN**tendo **E**n**V**ironm**E**nt, and I like to give my work names which are classical references. I might need to change the name to something drier like "NES Learning Environment" when I start to publicize it.

Currently have rudimentary reward and game over outputs for 
- Mario Brothers
- Donkey Kong
- Castlevania

Still need to figure out a method for delivering relevant roms (as well as a unified framework for the games)

At present, working on a proof of concept using the Arcade game Donkey Kong.