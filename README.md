# Game
You're a gopher in a 2d-burrow-maze placing bombs to clear out obstacles.

## Goals
- Use Ebitengine to create the 2d game
- Implement Gameloop, implement loading any sprites
- Add Player: Movable sprite
- Add layered sprites (background, foreground)
- 2d-Grid levels
- Add player and allow to move only along the grid collision detection
- add bombs and blast radius
- Add enemies with very simple AI
- Add Level Editor: based on textfile to load level
- Input via joystick
- Allow to render game in browser (WASM?)

# Concepts

## ECS - Entity component system
use this for spirte mangement, rendering, control

Sources:
- "How to build animations with ebiten using the ECS pattern" https://co0p.github.io/posts/ecs-animation/

- "Entity Component System | Game Engine series" by the Cherno https://www.youtube.com/watch?v=Z-CILn2w9K0 Summary: Great problem desciption of what problem ECS solves. Starting with the struggle of inheritence based system, then an improvement of Entity with a Components vector as a field, which has a leg up on inheritence (composition over inheritence) but suffers performnace hits due to no cache hits! Finally to a data-driven approach of ECS: associating a bunch of components by an Entity-ID. Where each components is part of an array, such that it has data locatlity and the associated lookup is problably a tree on those pure components type arrays. 


## tilemap
Use this for the stage representation

# Assets
- The tiles are from https://kenney.nl/assets/tiny-town and in the public domain (CC0 License)
- https://github.com/MariaLetta/free-gophers-pack/tree/master
License CC0 for those in the public domain.

- Gopher Pictures: The Go Gopher by Renee French is licensed under the Creative Commons Attribution 4.0 License.

# Dev Journal

## Donburi ECS System: Many layers of indirection
My first experience with the donburi ECS system was a daunting one, simply adopting it for a hello world example was very tough. First I took the example code for the "platformer" and removed all indirection to make it run in a single file. Still not fully understanding the modules involved (entities, components, systems).

Finally I've got a simple object rendered (a measily brown rectangle) I was embolden to try to emulate the donburi example project layout. This was even harder. After a few hours of working thourgh example code following the many branching paths accross multiple folders and files I got to adopt it to my code and rendered, yet again, a single measily brown rectangle.

What blocked me was that I wanted to start with the game arena as the background. But the platformer example didn't have an object corresponding to a "background". It instead placed platforms, floating platforms, walls and ramps directly into the scene. I looked at a game that did: "airplanes" by m110. But carefully reading through the code, finding my self being n-layers deep in indirections I found many gameobjects but not the background.
I gave up trying to find a solution and match it and simply decided its just another "Object" and created a "Arena" Archtype for it and successfuly reandered, you guessed it, a meassly brown rectangle as a placeholder for the glorious arena to follow hopefully soon.

With a better intuition of how the donburi ecs works, I'm looking forward
to building the game with it.

## Using ECS how to model the Game level?
That is the game takes place in an arena on 2d grid. That is rendered
based on a 2d-slice. But I don't see how other projects do it, it seems to be obscured.

I decided to create a 2d-Slice of ints, which I then map to `ebiten.Images`. I dislike this approach as it goes against the ECS pitch where you retrieve a contigious array of the thing you want to process and then iterate over it without pointer jumps, reducing cache misses.
Also I'm afraid that deriving tiles from a 2d-array TileID instead of having the tiles them instances of some "Tile"-Archtype, will later make it harder to built complex interactions with them, like placing items or rearranging them dynamically for creative gameplay mechanics.

## ECS: Open Source examples to the rescue 
I found myself when stuck on how to ECS for my gamedev that I would lean havily on the examples on `donburi`, like the Platformer or the `airplanes`-game by `m110` (link: https://github.com/m110/airplanes).
Though it was painful to undestand the code at first, even the simplest examples being spread over multiple files. I steadily got the hang of it.

What helped was aiming at a focus on a single feature first, like rendering the background level ("the arena") and if its too complex break it down to even smalle parts. Then move on to the next. The very first step was the hardest as the boilerplate is enormous, but the great thing about donburi's ECS is, that once start getting it, it scales greatly. 

What I love the most about ECS, I think, is that now I can much clearly imagine how to implement any feature I can think of. I love the feeling when your skills and knowledge suddenly catch up with what seemed to be insurmountable before, but it was "just" around the corner all along. I had this experience many times in IT, and I try to instill it in others.

