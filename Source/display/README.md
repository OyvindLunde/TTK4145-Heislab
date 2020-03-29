# Go Display for n elevators with m floors
**This module contains functionality for creating a display that shows n elevators with m floors, used in the NTNU course TTK4145 - Real-Time Programming.**

The display was created to make it easier to visualize and troubleshoot the system, as the simulator provided by the course is quite tricky to use to check if the elevators are synchronized, e.g. see which orders are active for the different elevators or check whether they registers each other orders, and so on. The picture below shows an example of the Display on a running system of 2 elevators with 5 floors.

![Elevator Display](images/ElevatorDisplay_prototype.PNG)

## Module content
The core of the module is based around the go package shiny, see https://github.com/golang/exp/tree/master/shiny. The module was created by taking inspiration from the examples provided by Nigeltao at https://github.com/golang/exp/tree/master/shiny/example. Note that the amount of this kind of functionality - as far as my googling skills go - is very sparse for Go, and to my knowledge these are the absolute best examples online. Nigeltaos examples provide a lot more advanced funcionality that can be useful in other projects.

### Colors - RGBA
The module uses RGBA to color the display. RGBA works almost exactly as RGB, just with the extra 'alpha' attribute whose only use is to make the color transparent. This is pretty much useless in this module, so if you want to add new colors, simply set the first three values exatly as you would with RGB, and set the 'A' value to whatever, e.g. 0x00 or 0xff (it doesn't matter which value you give it).

### Coordinate system
The coordinate system is defined such that (0,0) is at the top left corner, contrary to "normal" coordinate systems which uses the bottom left corner as origin. This also means that y increases downwards, and is important to remember when navigating in the display.


## Til den glade koker
Note that the display interacts with the elevator software, while the simulator provided by the course interacts with the elevator hardware. This means that this display code probably won't work directly for your system. You will most likely need to "tune" and alter the code a bit to make it work properly. If your solution to the task differs a lot from ours, it might be best to start with the file "DisplayStatic.go", found in the ExampleCode folder, and build on this code according to your own system [TODO: Oppdater DisplayStatic.go]. 

There is probably not much point in using the Display in the beginning of the project. It is probably most useful when you have a working (single) elevator, and want to test your system with multiple elevators and check if they communicate correctly. This, combined with a COVID-19 induced quarantine, was my motivation for creating the display.

### Modifying the code to fit your system
**Static components**<br/>
The static components of the system are displayed with the functions drawElevStatic() and displayElevStatic() to create and display the elevator, along with its floors and buttons. The static part of the module is thus designed such that any system should be able to display the basis (static) components of the elevator(s), and thus it is recommended to change these functions as little as possible.<br/>
The elevator title could also be a part of these functions, but it is not as it requires the elevators Id. However, since the ID most likely (depending on how you design your system) will be constant, you could modify the static functions to also creat and display the elevator title.

**Dynamic components**<br/>
The dynamic components require the following input from the system:
- The OrderQueue from each elevator. In our system, each order has these attributes:
  - Floor 
  - ButtonType
  - Status
  - Finished

All of these are needed to tell the module where to color a button, and which color to should have. Other attributes than ours will work, but you'll have to modify getOrderColor() to change how the different buttons are colored (and then modify displayOrderExplanations()).

- ElevInfo from each elevator. In our system, ElevInfo has these attributes: 
  - Id
  - Floor
  - CurrentOrder
  - State<br/>

Modify drawElevInfo() if you want to add or remove some attributes from the display. Floor needs to be given to the module in some way in order to set the floor indicator correctly.

- An "Updates" variable. This variable is used to tell the Display that something (one of the attributes above) has changed, and needs to be redrawn. You can not continuously redraw the display as it will make it crash. The goroutine updates() is used to tell the Display to redraw the (updated) display.


The sizes of the system are dynamic, i.e. it can display an arbitrary amount of floors and elevators. However, the display may not be wide enough to display more than 3-4 elevators, depending on the size of your screen. The number of buttons is fixed.

