# Project Nineveh

The goal of this project is to do to the Nintendo Entertainment System what the Arcade Learning Environment ([ALE](https://github.com/mgbellemare/Arcade-Learning-Environment))
did for the Atari; make it viable to write agents to play classic Nintendo games.
 
Nineveh is a Modification of the [Nintengo](https://github.com/nwidger/nintengo) emulator (built on the same commit history, with many of the original files included). The goals I have been working towards are:
-  Emulation core uncoupled from rendering and sound generation modules for fast emulation with minimal library dependencies.
-  Python development is supported through ctypes.

The system is designed with the intent of easy interfacing with the OpenAI [Gym](https://github.com/openai/gym). Unlike the ALE, which appears to have initially been desinged
to integrate with C++ agents out of the box. Further extension would be needed to enable its control via pipes (which seems to be something of a pain to implement)

At present, can either be launched in human playable mode or in headless mode. All agents will operate in headless mode. Each observation outputs a jpeg file; the option to output the entire
run to the current point as a gif is also present.

At time of writing, There are some unresolved irregularities with the output of the pixels, and frames will occasionally be dropped in headless mode. Resolving this issue is a high priority.

Currently using Castlevania for my testing and ongoing proof of concept (see the Alucard repository). Ideally, would like to have a suite of games to match those in the Atari 57 benchmark.
Will likely extend the Alucard proof of concept to also include Super Mario Brothers 3 before pivoting back to Wintermute (which will necessitate adding final fantasy)
