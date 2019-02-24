# What is this?

Space Elevator is a visualization tool for exploring elevator algorithms.

![Alt text](https://raw.githubusercontent.com/jaredfolkins/spaceelevator/master/screenshot.png "Space Elevator Simulation")

# What algorithms are available?

Nearest Car (NC): Elevator calls are assigned to the elevator best placed to answer that call according to three criteria that are used to compute a figure of suitability (FS) for each elevator. (1) If an elevator is moving towards a call, and the call is in the same direction, FS = (N + 2) - d, where N is one less than the number of floors in the building, and d is the distance in floors between the elevator and the passenger call. (2) If the elevator is moving towards the call, but the call is in the opposite direction, FS = (N + 1) - d.  (3) If the elevator is moving away from the point of call, FS = 1. The elevator with the highest FS for each call is sent to answer it. The search for the "nearest car" is performed continuously until each call is serviced.

https://www.quora.com/Is-there-any-public-elevator-scheduling-algorithm-standard

# What is the supported Golang Version?

`>= 1.11`

# How do I compile this?

From the project root run the following build script. The build script also creates a space a spaceelevator.tar.gz of all the binaries compressed to help with distribution.

With go modules you can place this application anywhere on your file system.

`export GO111MODULE=on`

`./build.sh`

# How do I execute this?

After compilation, you should have a bin folder with three binaries, select the one for your operating system.

### Windows 64bit
`./bin/spaceelevator.exe`

### OSX 64bit

`./bin/spaceelevator.darwin`

### Linux 64bit

`./bin/spaceelevator.amd64`

# What should I see?

Space Elevator should automatically open a browser window on port http//localhost:8989

# Why?

I was asked to implement a program that used an elevator algorithm for a job interview once and I decided I’d go big. I also really suck at whiteboard challenges. My career has lead me to know a ton of things across many domains in computer science. But it means that I only keep the gist of a solution in my head. I need a reference to read and the ability to do research and explore the problem domain with my fingers and face in front of a computer.

# Isn’t this a waste of time?

No.

# Then why’d you do it?

I like to practice the `N+1 Favorable Outcome Philosophy`.

For my kids or the students I mentor and help, the idea to is solve problems in ways that maximize the potential for N+1 FOP events. It is a philosophy that tries to align your efforts to maximize on as many favorable outcomes as possible.

In this case, if I would have submitted a lesser tool, that would have been absolutely fine. But I reasoned that many times, a job interview is a complete crapshoot. Especially a technical one. I’ve had great interviews and horrible interviews. ¯\_(ツ)_/¯ 

And if the interview doesn’t go well it is a complete sunk cost of my time.

So I figured that if I invest a little bit more time upfront, I could create a project that I could use to teach others from. Especially my small children as they approach the age where this interests them. The thinking is that it would have the potential to offer longer term yield. And it worked! Now I have my two olders wanting to turn this simulation into a game! :-)

# Why should I learn this?

Space Elevator is essentially teaching you about how a scheduler works when it is distributing jobs. A very common computer science problem. In fact, right now, on the device that you are operating from, there are jobs being scheduled using very similar techniques! TMYK

# Conclusion

At this time it works and is fun to interact with and observe!

# TODO

There are several things that could be improved

### Error handling

	I would create a channel via the scheduler to aggregate all the errors from the elevators

### Testing

	The code is pretty modular so wrapping more tests around it would be ok.

### O(n)

	Some of the seeking functions/methods are not optimized.

### Add more algorithms

	It would be cool to explore a few more with my kids.

