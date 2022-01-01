package main

import (
	"fmt"
	"os"
	"math"
)

func main() {
	path := "/home/namll/.wine/drive_c/Program Files/Sony/EverQuest/Logs/eqlog_Oxnull_P1999Green.txt"

	file := GetLog(path)

	if file == nil { return }

	defer file.Close()

	go GetInput()

	for{ ParseLog(file) }
}

func GetInput(){
	buffer := make([]byte, 512)
	for{
		os.Stdin.Read(buffer)
		line := GetLine(buffer, 0)
		if line != nil{
			fmt.Println(string(line))
		}
	}
}

func GetLog(path string) *os.File{
	file, err := os.Open(path)
	if err != nil {
		return nil
	}

	return file
}

func GetLine(buffer []byte, offset int) []byte{
	var i, j int = 0, offset
	for i = 0; i < len(buffer); i++ {
		if buffer[i] == '\n'{
			break
		}
	}

	if i - j == len(buffer){
		return nil
	}

	return buffer[j:i]
}

func ParseLog(file *os.File){
	buffer := make([]byte, 4096)
	_, _ = file.Read(buffer)

	line := GetLine(buffer, 0)
	if line != nil{
		start := int(math.Min(float64(len(line)), 27))
		fmt.Println(string(line[start:]))
	}
}
