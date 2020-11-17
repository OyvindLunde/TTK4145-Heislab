# Go Display for n elevators with m floors
**This module contains functionality for creating a display that shows n elevators with m floors, used in the NTNU course TTK4145 - Real-Time Programming.**

The display was created to make it easier to visualize and debug the system, as the simulator provided by the course is quite tricky to use to check if the elevators are synchronized, e.g. see which orders are active for the different elevators or check whether they register each other orders, and so on.

## About the Module
The core of the module is based around the go package shiny, see https://github.com/golang/exp/tree/master/shiny. The module was created by taking inspiration from the examples provided by Nigeltao at https://github.com/golang/exp/tree/master/shiny/example. Note that the amount of this kind of functionality - as far as my googling skills go - is very sparse for Go, and to my knowledge these are the absolute best examples online. Nigeltaos examples provide a lot more advanced funcionality that can be useful in other projects.

### Colors - RGBA
The module uses RGBA to color the display. RGBA works almost exactly as RGB, just with the extra 'alpha' attribute whose only use is to make the color transparent. This is pretty much useless in this module, so if you want to add new colors, simply set the first three values exatly as you would with RGB, and set the 'A' value to whatever, e.g. 0x00 or 0xff (it doesn't matter which value you give it).

### Coordinate system
The coordinate system is defined such that (0,0) is at the top left corner, contrary to "normal" coordinate systems which uses the bottom left corner as origin. This also means that y increases downwards, and is important to remember when navigating in the display.


## Til den glade koker
Note that the display interacts with the elevator software, while the simulator provided by the course interacts with the elevator hardware. This means that this display code probably won't work directly for your system. You will most likely need to "tune" and alter the code a bit to make it work properly.  

There is probably not much point in using the Display in the beginning of the project. It is probably most useful when you have a working single elevator, and want to test your system with multiple elevators and check if they communicate correctly. This, combined with a COVID-19 induced quarantine, was actually my motivation for creating the display.

### Modifying the code to fit your system
**Static components**<br/>
The static components of the system are displayed with the functions drawElevStatic() and displayElevStatic() to create and display the elevator, along with its floors and buttons. The static part of the module is thus designed such that any system should be able to display the basis (static) components of the elevator(s), and thus it is recommended to change these functions as little as possible.<br/>
The elevator title could also be a part of these functions, but it is not as it requires the elevators Id, which is different for each elevator. However, since the ID most likely (depending on how you design your system) will be constant, you could modify the static functions to also create and display the elevator title, instead of having it in the dynamic part.

**Dynamic components**<br/>
The dynamic components require the following input from the system:
- The OrderQueue from each elevator. In our system, each order has these attributes:
  - Floor 
  - ButtonType
  - Status
  - Finished
  - Confirmed

The first four attributes are needed to tell the module where to color a button, and which color it should have. Other attributes than ours will work, you'll just have to modify getOrderColor() to change how the different buttons are colored (and then modify drawOrderExplanations()).

- ElevInfo from each elevator. In our system, ElevInfo has these attributes: 
  - Id
  - Floor
  - CurrentOrder
  - State<br/>

Modify drawElevInfo() if you want to add or remove some attributes from the display. Floor needs to be given to the module in some way in order to set the floor indicator correctly.

- An "Updates" variable. This variable is used to tell the Display that something (one of the attributes above) has changed, and needs to be redrawn. You should not continuously redraw the display as the display will be constantly flickering. The goroutine updates() uses the Updates variable to tell the Display to redraw the (updated) display by generating a paint event.


The sizes of the system are dynamic, i.e. it can display an arbitrary amount of floors and elevators. The default size of the screen is set to display two elevators, which can be adjusted by changing the "width" attribute in the NewWindow() function (in the beginning of Display()). The number of buttons is fixed.

