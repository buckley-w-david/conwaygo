# conwaygo

Conway's Game of Life implementation

## Usage
```bash
$ go build cmd/life/life.go
$ ./life -h
Usage of ./life:
  -debug
      Display in debug mode
  -delay float
      Seconds between updates (default 0.2)
```

When the binary in run, you will be presented with a blank screen.  
Click on the terminal where you wish the seed live cells. You can do this at anytime, even after you have started the simulation.


### Keybindings
| Key         | Description                          |
| ----------- | ------------------------------------ |
| `space`     | Start and stop the simulation        |
| `Page Up`   | Increase the speed of the simulation |
| `Page Down` | Decrease the spped of the simulation |
| `ctrl`+`d`  | Enable debug mode.                   |
| arrow keys  | Move the viewport                    | 

In debug mode instead of displaying live cells, each tile will show the number living cells surrounding it.

## Example
![](https://github.com/buckley-w-david/conwaygo/raw/1279507e3138381d4c9dc362b632564aea9b481e/resources/conway.gif)

*Note* This gif is a little out of date, now that the code is split up into a `pkg` and `cmd` directory, build the project by running `go build cmd/life/life.go`
