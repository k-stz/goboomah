# Game
You're a gopher in a 2d-burrow-maze placing bombs to clear out obstacles.

![Chain explosion](/assets/gameplay.webp)

"Chain explosion example. Go, Boom -Ah!"


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
use this to management components sprite mangement, rendering, control

Sources:
- "How to build animations with ebiten using the ECS pattern" https://co0p.github.io/posts/ecs-animation/

- "Entity Component System | Game Engine series" by "the Cherno" https://www.youtube.com/watch?v=Z-CILn2w9K0 Summary: Great problem desciption of what problem ECS solves. Starting with the struggle of inheritence based system, then an improvement using Entity with a Components vector as a field, which has a leg up on inheritence (composition over inheritence) but suffers performnace hits due to no cache hits, as its an array of pointers! 
    - Finally to a data-driven approach of ECS: associating a bunch of components by an Entity-ID. Where each components is part of an array, such that it has data locatlity and the associated lookup is problably a tree on those pure components type arrays. 


## tilemap
Use this for the stage representation

## Collision Detection / Response Libary
- Using SolarLune's Resolv library.
    
Objects that require physics-based collision detection and response will be managed using the Resolv library. Each object is assigned a bounding box, circle, or convex polygon and placed in a 2D space. This space is divided into cells, allowing for efficient collision detection. Instead of checking every object against every other object (which would have a worst-case performance of `O(n²)`), the system only checks objects within neighboring cells, significantly improving performance. 

Spatial Partitioning and Complexity: Recognizing that objects usually interact with nearby objects, we can use a common pattern to group them for efficient lookup and manipulation. This technique is known as spatial partitioning (https://gameprogrammingpatterns.com/spatial-partition.html). 

As is commmon in computer science, this optimization involves a **space-time tradeoff**: We trade **CPU cycles (time) for memory (Space)**, by organizing the collision space into a more efficiently searchable memory construct. This can reduces collision detection complexity from a native implementation needing O(n²) to an O(n log n) using quadtree and even O(1) for Spatial hashing in the best case!

> Complexity: Why `O(n log n)`: 
- if using a quadtree, inserting and retreiving each object takges `O(log n)` 
- Doing this fo rn objects takes `O(n log n)` total complexity in the worst case
- However for uniform grids or hashing it can be closer to `O(n)` in practice (n Objects each needing just `O(1)`!)

Different types of game objects can have different collision responses, even if they share the same bounding shape. To achieve this, "tags" (resolv.tags) are used. For example, an object can be set to only detect collisions with objects tagged as "TagWall." When a collision is detected, the object is displaced just enough to prevent penetration, ensuring it only touches the wall instead of passing through it. This displacement is calculated using the Minimal Translation Vector (MTV), a core feature of the library. The MTV allows for smooth interactions, such as a player character "hugging" a wall without getting stuck or tunneling through it.

# Licenses / Credit
## Assets
- The tiles are from https://kenney.nl/assets/tiny-town and in the public domain (CC0 License)
- https://github.com/MariaLetta/free-gophers-pack/tree/master
License CC0 for those in the public domain.
- Explosion effect is from Statik64 and is CC0 (public domain)
- Gopher Pictures: The Go Gopher by Renee French is licensed under the Creative Commons Attribution 4.0 License.

## Libararies
- ebitengine: " A dead simple 2D game engine for Go" by Hajime Hoshi, Apache-2.0 license 
- donburi: "ECS library for Go/Ebitengine", MIT License, (c) 2022 Yota Hamada
- resolv: "A Simple 2D Golang collision detection and resolution library" by "SolarLune", MIT License Copyright (c) 2018-2021 SolarLune

# TODOs
- [ ] add "resolv" physics engine (great examples)
- [ ] derive solid tiles from the Arena tilemap and add to "resolv" space
- [ ] Add animation, see here: https://github.com/egonelbre/gophers

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

## Golang Libraries: v3 isn't always newer than v2
Found the Animation library ganim8 and wanted to use it. Found that the example code didn't work for me, because I did a simple "go get github.com/yohamta/ganim8" which got me the v1 version.
A quick look at the pkg.go.dev registry showed me that v3 is available, but to my confusion this tag is not available on the repository...

It seems the tag was deleted/abandoned by the library and indeed v2 is the more up-to-date library. I acertained that, by comparing the package API by looking over both versions Documentation and saw that at least a single Function "isEnd" was newer on v2. Also v3 was published, according to go.dev.pkg in 2022, while v2 was published in June 2023.

So I promptly updated all my refernced to `/v2`, did a `go mod tidy` and made sure the v3 library isn't referenced anywhere in my go.mod or go.sum.