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
https://github.com/MariaLetta/free-gophers-pack/tree/master
License CC0 for those is public domain.

The Go Gopher by Renee French is licensed under the Creative Commons Attribution 4.0 License.

# Gamedev questions

## Using ECS how to model the Game level?
That is the game takes place in an arena on 2d grid. That is rendered
based on a 2d-slice. But I don't see how other projects do it, it seems to be obscued
