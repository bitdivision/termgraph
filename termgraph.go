package main

//TODO
// - Add vertical graphs
// - Add colours

import (
    "encoding/csv"
    "flag"
    "fmt"
    //"errors"
    "log"
    "os"
    "strconv"
    "syscall"
    "unsafe"
)

var tick rune = '▇'
var sm_tick rune = '|'

//Gets the dimensions of the current terminal. Taken from the crypto package
func getSize(fd int) (width, height int, err error) {
    var dimensions [4]uint16

    if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dimensions)),     0, 0, 0); err != 0 {
        return -1, -1, err
    }
    return int(dimensions[1]), int(dimensions[0]), nil
}


func convertStrings(values []string) []float64 {
    output := make([]float64, len(values))
    for i, val := range values {
        output[i], _ = strconv.ParseFloat(val, 64)
        //if err != nil {log.Fatalf("Could not convert string values to floats: %v", err)}
    }
    return output
}

func printGraphs(headers []string, values []string, width int) {

    max_head := 0
    //Find the size of the largest header string to align graphs
    for _, val := range headers {
        if max_head < len(val) {
            max_head = len(val)
        }
    }


    new_vals := convertStrings(values)
    //Find max value to calculate scale
    max_val, min_val := 0.0, 0.0
    max_val_len := 0
    for _, val := range new_vals {
        if max_val < val {
            max_val = val
            max_val_len = len(fmt.Sprintf("%v", val))
        }
        if min_val > val {
            min_val = val
        }
    }

    //Use the max header size and given width to calc graph width
    //Currently scales to terminal width - 1
    graph_width:= width-(max_head+2)-(max_val_len + 3)
    //Calculate scale
    scale := float64(graph_width) / (max_val - min_val)
    
    //Loop through the values and print them
    for i, val := range new_vals {
        size := int(scale * val)
        fmt.Printf("%s: ", headers[i])

        //Print the requisite number of space to align (max_head+2)-len
        no_spaces := (max_head) - len(headers[i])
        for j := 0; j < no_spaces; j++ {
            fmt.Printf(" ")
        }

        if size == 0 {
            fmt.Printf("%v", sm_tick)
        } else {
            for v := 0; v < size; v++ {
                fmt.Printf("%c", tick)
            }
        }

        //Now print the actual values
        fmt.Printf("  %v\n", val)
    }

}

func main() {
    //Get the width of the terminal
    term_width,_,_ := getSize(0)

    //orientation := flag.String("Orientation", "v", "The orientation of the graph (h or v)")

    width := flag.Int("Width", term_width, "The width of the entire graph including labels and values")

    flag.Parse()


    //Check if a filename was provided
    filename := flag.Arg(0)
    if filename == "" {
        log.Fatal("No filename specified. Exiting...")
    }

    //Open the file
    csvFile, err := os.Open(filename)
    if err != nil {
        log.Fatalf("Could not open file: %v", err)
    }
    defer func() {
        csvFile.Close()
    }()

    //Open the file and read in the data
    r := csv.NewReader(csvFile)
    lines, err := r.ReadAll()
    if err != nil {
        log.Fatalf("Could not read specified file: %v ", err)
    }

    //Make sure only two lines in CSV file: header and values
    if len(lines) != 2 {
        log.Fatal("CSV format incorrect. Should be two lines: headers and values")
    }

    headers, values := lines[0], lines[1]

    if len(headers) != len(values) {
        log.Fatal("Data mismatch: different number of headers and values")
    }

    //Display the graph
    printGraphs(headers, values, *width)
    
}

