# Project Nineveh

The goal of this project is to do to the Nintendo Entertainment System what the Arcade Learning Environment ([ALE](https://github.com/mgbellemare/Arcade-Learning-Environment))
did for the Atari; make it viable to write agents to play classic Nintendo games.
 
Nineveh is a Modification of the [Nintengo](https://github.com/nwidger/nintengo) emulator. The goals I have been working towards are:
-  Emulation core uncoupled from rendering and sound generation modules for fast emulation with minimal library dependencies.
-  Python development is supported through ctypes.
Both of these have been accomplished at present, in a rudimentary fashion. Current product is very, very rough, requiring considerable polish

May wind up orphaned if I can find a better alternative to this