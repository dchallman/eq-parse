package main

import (
	"os"
	"math"
	"regexp"
	"os/exec"
)

func GetInput(){
	for{
		ParseInput(os.Stdin)
	}
}

func GetLog(file string) *os.File{
	log, err := os.Open(file)
	if err != nil {
		return nil
	}

	return log
}

/*
	 This should be probably be changed using bufio
*/
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

func ParseInput(file *os.File){
	buffer := make([]byte, 512)
	_, _ = file.Read(buffer)

	line := GetLine(buffer, 0)
	if line != nil{
		str := string(line)

		if (ParseInputItem(str)){}
	}
}

func ParseInputItem(str string) bool{
	re := regexp.MustCompile(`item: (.+)`)

	if re.MatchString(str){
		GetItemDB(re.FindStringSubmatch(str)[1])
		return true
	}

	return false
}

func ParseLog(file *os.File){
	buffer := make([]byte, 4096)
	_, _ = file.Read(buffer)

	line := GetLine(buffer, 0)
	if line != nil{
		offset := 27
		start := int(math.Min(float64(len(line)), float64(offset)))
		str := string(line[start:])

		if ParseItemBuy(str){ return }

		if ParseItemSell(str){ return }

		if ParseYouSayBuy(str){ return }

		if ParseYouSayItem(str){ return }

		if ParseZone(str){ return }
	}
}

func ParseItemBuy(str string) bool{
	re := regexp.MustCompile(`(.+) tells you, 'That'll be (\d+ platinum)?\s?(\d+ gold)?\s?(\d+ silver)?\s?(\d+ copper)? ((for the)|per) (.+).'`)
	if re.MatchString(str){
			parts := re.FindAllStringSubmatch(str, -1)
			parts = CleanParts(parts)
			WriteBuyToDB(parts)
			return true
		}

	return false
}

func ParseItemSell(str string) bool{
		re := regexp.MustCompile(`(.+) tells you, 'I'll give you (\d+ platinum)?\s?(\d+ gold)?\s?(\d+ silver)?\s?(\d+ copper)? ((for the)|per) (.+).'`)
		if re.MatchString(str){
			parts := re.FindAllStringSubmatch(str, -1)
			parts = CleanParts(parts)
			WriteSellToDB(parts)
			return true
		}

	return false
}

func ParseYouSayBuy(str string) bool{
		re := regexp.MustCompile(`You say, '(\d+)'`)
		if re.MatchString(str){
			parts := re.FindAllStringSubmatch(str, -1)

			cmd := exec.Command("./parseBuy.sh", parts[0][1])
			cmd.Run()
			return true
		}

	return false
}

func ParseYouSayItem(str string) bool{
		re := regexp.MustCompile(`You say, '(.+)'`)
		if re.MatchString(str){
			parts := re.FindAllStringSubmatch(str, -1)
			GetItemDB(parts[0][1])
			return true
		}

	return false
}

func ParseZone(str string) bool{
		re := regexp.MustCompile(`You have entered (.+)\.`)
		if re.MatchString(str){
			parts := re.FindAllStringSubmatch(str, -1)
			SetZone(parts[0][1])
			return true
		}

	return false
}

func SetZone(str string){
	ZONE = str;
	file, _ := os.Create("zone.txt")
	file.WriteString(ZONE + "\n")

	file.Close()
}

func GetZone(){
	file, _ := os.Open("zone.txt")
	buffer := make([]byte, 512)
	_, _ = file.Read(buffer)

	line := GetLine(buffer, 0)

	if line != nil {
		ZONE = string(line)
	}
}

func CleanParts(parts [][]string) [][]string{
	re := regexp.MustCompile(`(\d+) \w`)

	if re.MatchString(parts[0][PLAT]){
		parts[0][PLAT] = re.FindStringSubmatch(parts[0][PLAT])[1]
	}else{
		parts[0][PLAT] = "0"
	}

	if re.MatchString(parts[0][GOLD]){
		parts[0][GOLD] = re.FindStringSubmatch(parts[0][GOLD])[1]
	}else{
		parts[0][GOLD] = "0"
	}

	if re.MatchString(parts[0][SILVER]){
		parts[0][SILVER] = re.FindStringSubmatch(parts[0][SILVER])[1]
	}else{
		parts[0][SILVER] = "0"
	}

	if re.MatchString(parts[0][COPPER]){
		parts[0][COPPER] = re.FindStringSubmatch(parts[0][COPPER])[1]
	}else{
		parts[0][COPPER] = "0"
	}

	return parts
}
