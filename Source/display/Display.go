package display

// TODO:
// Må ta inn full log (elevInfo og orders) fra de andre heisene for å kunne displaye alt
// Legges inn i displayOtherElevators. Nice å ha en logList i logmanagement?

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"strconv"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"

	"../logmanagement"
)

var (
	black     = color.RGBA{0x00, 0x00, 0x00, 0x00}
	blue0     = color.RGBA{0x00, 0x00, 0x1f, 0xff}
	blue1     = color.RGBA{0x00, 0x00, 0x3f, 0xff}
	darkGray  = color.RGBA{0x3f, 0x3f, 0x3f, 0xff}
	lightGray = color.RGBA{0xd8, 0xd8, 0xd8, 0x7f}
	green     = color.RGBA{0x16, 0xee, 0x50, 0x7f}
	red       = color.RGBA{0xff, 0x00, 0x00, 0x7f}
	yellow    = color.RGBA{0xef, 0xff, 0x00, 0x3f}
	white     = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

const numFloors = 4
const numButtons = 3

var btnPanel_x0 = 300 // "Start" position of Button Panel (x and y coordinate, top left corner)
var btnPanel_y0 = 70
var btnPanel_x = 160                     // Width of Button Panel
var btnPanel_y = 200                     // Height of Button Panel
var btn_size_x = btnPanel_x / numButtons // Width of button in the Button Panel
var btn_size_y = btnPanel_y / numFloors  // Height of button in the Button Panel

func Display() {
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title: "Elevator Display",
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Display up and running")
		defer w.Release()

		// Static components
		elevStatic := drawElevStatic(s, black) // Basic elevator layout
		orderExpl := drawOrderExplanation(s)
		arrow := drawArrowLeft(s, 30, 20, black, lightGray) // Arrow to use as floor indicator

		go update(w)

		var sz size.Event
		for {
			time.Sleep(20 * time.Millisecond)

			e := w.NextEvent()
			switch e := e.(type) {

			case paint.Event:
				paintScreen(w, sz, lightGray, blue0) // Paint background and border of screen in the selected colors
				displayOrderExplanations(w, orderExpl)
				//fmt.Println(&logmanagement.ElevInfo)
				displayLocalElevator(w, s, elevStatic, logmanagement.GetOrderList(), logmanagement.GetMyElevInfo(), arrow)
				displayOtherElevators(w, s, elevStatic, logmanagement.GetOtherElevInfo(), arrow)

			case size.Event: // Do not remove this
				sz = e

			case error:
				log.Print(e)
			}
		}
	})
}

func update(w screen.EventDeque) {
	for {
		time.Sleep(20 * time.Millisecond)
		if logmanagement.GetDisplayUpdates() {
			w.Send(paint.Event{})
			logmanagement.SetDisplayUpdates(false)
		}
	}

}

func displayOtherElevators(w screen.Window, s screen.Screen, elevStatic []screen.Texture, elevList []logmanagement.Elev, arrow screen.Texture) {
	for i := 0; i < len(elevList); i++ {
		//logmanagement.PrintOrderQueue(elevList[i].Orders)
		displayElevStatic(w, elevStatic, btnPanel_x0+300*(i+1))
		elevatorTitle := drawText(s, "Elevator "+strconv.Itoa(elevList[i].Id)+" overview", 155, 20) // To improve runtime: change so that this is only calculated once
		w.Copy(image.Point{btnPanel_x0 + 300*(i+1) + btn_size_x/2 - 25, btnPanel_y0 - 25}, elevatorTitle, elevatorTitle.Bounds(), screen.Src, nil)
		displayElevInfo(w, drawElevInfo(s, elevList[i]), i+1)
		displayFloorIndicator(w, arrow, elevList[i].Floor, i+1)
		displayOrders(w, s, elevList[i].Orders, btnPanel_x0+300*(i+1))
		//displayOrders(w, s, logmanagement.MyElevInfo.Orders, btnPanel_x0+300*(i+1))
	}
}

func displayLocalElevator(w screen.Window, s screen.Screen, elevStatic []screen.Texture, queue [numFloors][numButtons]logmanagement.Order, elevInfo logmanagement.Elev, arrow screen.Texture) {
	displayElevStatic(w, elevStatic, btnPanel_x0)
	displayLocalElevDynamic(w, s, queue, elevInfo, arrow, btnPanel_x0)
}

func displayLocalElevDynamic(w screen.Window, s screen.Screen, queue [numFloors][numButtons]logmanagement.Order, elevInfo logmanagement.Elev, arrow screen.Texture, start_x int) {
	// Display Elevator title with correct Id
	elevatorTitle := drawText(s, "Elevator "+strconv.Itoa(elevInfo.Id)+" overview", 155, 20) // To improve runtime: change so that this is only calculated once
	w.Copy(image.Point{start_x + btn_size_x/2 - 25, btnPanel_y0 - 25}, elevatorTitle, elevatorTitle.Bounds(), screen.Src, nil)
	displayElevInfo(w, drawElevInfo(s, elevInfo), 0)
	displayFloorIndicator(w, arrow, elevInfo.Floor, 0)
	displayOrders(w, s, queue, btnPanel_x0)
}

func displayElevStatic(w screen.Window, list []screen.Texture, start_x int) {
	// Draw Elevator (rectangle)
	w.Copy(image.Point{start_x, btnPanel_y0}, list[0], list[0].Bounds(), screen.Src, nil)

	offset := int(float64(btnPanel_y0) + (numFloors-0.5)*float64(btn_size_y)) // Bottom value of y

	// Draw button Types
	w.Copy(image.Point{start_x + btn_size_x/2 - 10, offset + btn_size_y/2 + 5}, list[1], list[1].Bounds(), screen.Src, nil)
	w.Copy(image.Point{start_x + btn_size_x*3/2 - 18, offset + btn_size_y/2 + 5}, list[2], list[2].Bounds(), screen.Src, nil)
	w.Copy(image.Point{start_x + btn_size_x*5/2 - 15, offset + btn_size_y/2 + 5}, list[3], list[3].Bounds(), screen.Src, nil)

	// Draw floor numbers (for n floors)
	for i := 0; i < numFloors; i++ {
		w.Copy(image.Point{start_x - 15, offset - 10 - i*btn_size_y}, list[i+4], list[i+4].Bounds(), screen.Src, nil)
	}
}

func drawElevStatic(s screen.Screen, color color.RGBA) []screen.Texture {
	list := make([]screen.Texture, 0)
	buttonPanel := drawButtonPanel(s, color)
	UP, DOWN, CAB := drawButtonTypes(s)
	floorNumbers := drawElevatorFloorNumbers(s)

	list = append(list, buttonPanel)
	list = append(list, UP)
	list = append(list, DOWN)
	list = append(list, CAB)
	list = append(list, floorNumbers...)

	return list
}

func displayOrders(w screen.Window, s screen.Screen, queue [numFloors][numButtons]logmanagement.Order, start_x int) {
	for i := 0; i < numFloors; i++ {
		for j := 0; j < numButtons; j++ {
			color := getOrderColor(queue[i][j])
			//fmt.Println(s)
			button := drawButton(s, btn_size_x, btn_size_y, color)
			displayButton(w, button, i, j, start_x)
		}
	}
}

func getOrderColor(order logmanagement.Order) color.RGBA {
	if order.Finished == true {
		return red
	}
	if order.Status > 0 {
		return green
	}
	if order.Status == 0 {
		return yellow
	}
	return black
}

func displayButton(w screen.Window, button screen.Texture, floor int, btnType int, start_x int) {
	w.Copy(image.Point{start_x + btnType*btn_size_x + 1, btnPanel_y0 + (numFloors-floor-1)*btn_size_y + 1}, button, button.Bounds(), screen.Src, nil)
}

func drawButton(s screen.Screen, x int, y int, color color.RGBA) screen.Texture {
	//fmt.Println(s)
	/*fmt.Println(x)
	fmt.Println(y)
	fmt.Println(color)*/
	size0 := image.Point{x - 1, y - 1} // -1 to avoid painting over the white lines in the button panel
	temp, _ := s.NewBuffer(size0)
	m := temp.RGBA()
	b := temp.Bounds()

	for i := b.Min.X + 1; i < b.Max.X-1; i++ {
		for j := b.Min.Y + 1; j < b.Max.Y-1; j++ {
			m.SetRGBA(i, j, color)
		}
	}

	rect, _ := s.NewTexture(size0)
	rect.Upload(image.Point{}, temp, temp.Bounds())
	return rect
}

func displayFloorIndicator(w screen.Window, arrow screen.Texture, floor int, elevNumber int) { // Could be better to change elevNumber to start_x
	pos_0_x := btnPanel_x0 + numButtons*btn_size_x + 5 + 300*elevNumber
	pos_0_y := int(float64(btnPanel_y0) + 3.5*float64(btn_size_y) - 10)
	w.Copy(image.Point{pos_0_x, pos_0_y - floor*btn_size_y}, arrow, arrow.Bounds(), screen.Src, nil)
}

func drawArrowLeft(s screen.Screen, x int, y int, color color.RGBA, backgroundColor color.RGBA) screen.Texture {
	size0 := image.Point{x, y}
	temp, _ := s.NewBuffer(size0)
	m := temp.RGBA()
	b := temp.Bounds()

	for i := (b.Max.X) / 3; i < b.Max.X; i++ {
		for j := b.Min.Y; j < b.Max.Y; j++ {
			if j > (b.Max.Y)/4 && j < (b.Max.Y)*3/4 {
				m.SetRGBA(i, j, color)
			} else {
				m.SetRGBA(i, j, backgroundColor)
			}

		}
	}

	for i := b.Min.X; i < (b.Max.X)/3; i++ {
		for j := b.Min.Y; j < (b.Max.Y)/2; j++ {
			if i+j > (b.Max.Y)/2 {
				m.SetRGBA(i, j, color)
			} else {
				m.SetRGBA(i, j, backgroundColor)
			}
		}
	}

	for i := b.Min.X; i < (b.Max.X)/3; i++ {
		for j := (b.Max.Y) / 2; j < b.Max.Y; j++ {
			if i+(b.Max.X)/3 > j {
				m.SetRGBA(i, j, color)
			} else {
				m.SetRGBA(i, j, backgroundColor)
			}
		}
	}

	arrow, _ := s.NewTexture(size0)
	arrow.Upload(image.Point{}, temp, temp.Bounds())
	return arrow
}

func displayElevInfo(w screen.Window, elevInfo []screen.Texture, elevNum int) {
	end_y := int(float64(btnPanel_y0) + numFloors*float64(btn_size_y))
	for i := 0; i < len(elevInfo); i++ {
		w.Copy(image.Point{btnPanel_x0 + 300*elevNum, end_y + btn_size_y + 20*i}, elevInfo[i], elevInfo[i].Bounds(), screen.Src, nil)
	}
}

func drawElevInfo(s screen.Screen, elev logmanagement.Elev) []screen.Texture {
	elevInfo := make([]screen.Texture, 0)

	id := drawText(s, "ID: "+strconv.Itoa(elev.Id), 160, 20)
	floor := drawText(s, "Floor: "+strconv.Itoa(elev.Floor), 160, 20)
	state := drawText(s, "State: "+convState2String(elev.State), 160, 20)
	currentOrder := drawText(s, "CurrentOrder: "+convOrder2String(elev.CurrentOrder), 160, 20)

	elevInfo = append(elevInfo, id)
	elevInfo = append(elevInfo, floor)
	elevInfo = append(elevInfo, state)
	elevInfo = append(elevInfo, currentOrder)

	return elevInfo
}

func displayOrderExplanations(w screen.Window, orderExpl []screen.Texture) {
	w.Copy(image.Point{10, 70}, orderExpl[0], orderExpl[0].Bounds(), screen.Src, nil)
	w.Copy(image.Point{10, 10}, orderExpl[1], orderExpl[1].Bounds(), screen.Src, nil)
	w.Copy(image.Point{10, 30}, orderExpl[2], orderExpl[2].Bounds(), screen.Src, nil)
	w.Copy(image.Point{10, 50}, orderExpl[3], orderExpl[3].Bounds(), screen.Src, nil)
}

func drawOrderExplanation(s screen.Screen) []screen.Texture {
	text := make([]screen.Texture, 0)

	blackOrder := drawText(s, "Black: No order for this button", 250, 20)
	redOrder := drawText(s, "Red: Order finished", 250, 20)
	yellowOrder := drawText(s, "Yellow: Pending Order", 250, 20)
	greenOrder := drawText(s, "Green: Order being executed", 250, 20)

	text = append(text, blackOrder)
	text = append(text, redOrder)
	text = append(text, yellowOrder)
	text = append(text, greenOrder)
	return text
}

func drawButtonPanel(s screen.Screen, color color.RGBA) screen.Texture {
	size0 := image.Point{btnPanel_x, btnPanel_y}
	temp, _ := s.NewBuffer(size0)
	m := temp.RGBA()
	b := temp.Bounds()

	for i := b.Min.X; i < b.Max.X; i++ {
		for j := b.Min.Y; j < b.Max.Y; j++ {
			m.SetRGBA(i, j, color)
		}
	}

	drawHorizontalLines(m, numFloors-1, white)
	drawVerticalLines(m, numButtons-1, white)

	rect, _ := s.NewTexture(size0)
	rect.Upload(image.Point{}, temp, temp.Bounds())
	return rect
}

func drawButtonTypes(s screen.Screen) (screen.Texture, screen.Texture, screen.Texture) {
	UP := drawText(s, "UP", 19, 20)
	DOWN := drawText(s, "DOWN", 35, 20)
	CAB := drawText(s, "CAB", 27, 20)

	return UP, DOWN, CAB
}

func drawElevatorFloorNumbers(s screen.Screen) []screen.Texture {
	list := []screen.Texture{}
	for i := 0; i < numFloors; i++ {
		floorNumber := drawText(s, strconv.Itoa(i), 10, 20)
		list = append(list, floorNumber)
	}

	return list
}

func paintScreen(w screen.Window, sz size.Event, backgroundColor color.RGBA, borderColor color.RGBA) {
	const inset = 10
	//fmt.Println(sz)
	for _, r := range imageutil.Border(sz.Bounds(), inset) {
		w.Fill(r, borderColor, screen.Src) // Paint border of screen
	}
	w.Fill(sz.Bounds().Inset(inset), backgroundColor, screen.Src) // Paint screen
}

// Functions for converting different data types to strings that can be displayed

func convState2String(state int) string {
	if state == 0 {
		return "INIT"
	}
	if state == 1 {
		return "IDLE"
	}
	if state == 2 {
		return "EXECUTE"
	}
	if state == 3 {
		return "RESET"
	}
	if state == 4 {
		return "LOST"
	}
	return "NOT CONNECTED"
}

func convOrder2String(order logmanagement.Order) string {
	if order.ButtonType == 0 {
		return strconv.Itoa(order.Floor) + "UP"
	}
	if order.ButtonType == 1 {
		return strconv.Itoa(order.Floor) + "DOWN"
	}
	if order.ButtonType == 2 {
		return strconv.Itoa(order.Floor) + "CAB"
	}
	return "NONE"
}

// The most basic functions for drawing text and lines

func drawText(s screen.Screen, text string, x_size int, y_size int) screen.Texture {
	floor0 := image.Point{x_size, y_size}
	f0, err := s.NewBuffer(floor0)

	drawRGBA(f0.RGBA(), text)

	f01, err := s.NewTexture(floor0)
	if err != nil {
		log.Fatal(err)
	}
	f01.Upload(image.Point{}, f0, f0.Bounds())
	defer f0.Release()
	return f01
}

func drawRGBA(m *image.RGBA, str string) {
	draw.Draw(m, m.Bounds(), image.White, image.Point{}, draw.Src)

	d := font.Drawer{
		Dst:  m,
		Src:  image.Black,
		Face: inconsolata.Regular8x16,
		Dot: fixed.Point26_6{
			Y: inconsolata.Regular8x16.Metrics().Ascent,
		},
	}
	//d.DrawString(fmt.Sprint(tp))
	d.DrawString(str)
}

func drawHorizontalLines(m *image.RGBA, num int, color color.RGBA) {
	b := m.Bounds()
	intervall := (b.Max.Y - b.Min.Y) / (num + 1)
	for i := 0; i <= num+1; i++ {
		drawHorizontalLine(m, intervall*i, color)
	}
	//drawHorizontalLine(m, b.Max.Y-1, color) // Supposed to draw a bottom line, but had no effect. Don't know why.
}

func drawHorizontalLine(m *image.RGBA, y int, color color.RGBA) {
	b := m.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		m.SetRGBA(x, y, color)
	}
}

func drawVerticalLines(m *image.RGBA, num int, color color.RGBA) {
	b := m.Bounds()
	intervall := (b.Max.X - b.Min.X) / (num + 1)
	for i := 0; i <= num+1; i++ {
		drawVerticalLine(m, intervall*i, color)
	}
}

func drawVerticalLine(m *image.RGBA, x int, color color.RGBA) {
	b := m.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		m.SetRGBA(x, y, color)
	}
}
