package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	var insensitive bool
	var useBox bool
	var up int
	var down int
	var bgColorFlag string
	var fgColorFlag string

	flag.BoolVar(&insensitive, "i", false, "Case insensitive search")
	flag.BoolVar(&useBox, "box", false, "Use box highlighting")
	flag.BoolVar(&useBox, "b", false, "Use box highlighting (shorthand)")
	flag.IntVar(&up, "up", 0, "Lines of context to show after match")
	flag.IntVar(&down, "down", 0, "Lines of context to show before match")
	flag.StringVar(&bgColorFlag, "bg", "", "Background color of the match\nColors: [red, green, yellow, blue, magenta, cyan, white, black, pink]")
	flag.StringVar(&fgColorFlag, "fg", "", "Foreground color of the match\nColors: [red, green, yellow, blue, magenta, cyan, white, black, pink]")

	flag.Usage = func() {
		fmt.Printf("Usage: grepz [-i] [-box|-b] [-up num] [-down num] <search_term> [<input_file>]\n")
		fmt.Printf("Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		return
	}

	searchTerm := args[0]

	validColors := []string{"red", "green", "yellow", "blue", "magenta", "cyan", "white", "black", "pink"}

	if !contains(validColors, bgColorFlag) {
		fmt.Printf("Error: %s is not a valid background color\n\n", bgColorFlag)
		flag.Usage()
		return
	}
	if !contains(validColors, fgColorFlag) {
		fmt.Printf("Error: %s is not a valid foreground color.\n\n", fgColorFlag)
		flag.Usage()
		return
	}

	input, err := getInput(args)
	if err != nil {
		log.Fatal(err)
	}

	fileContent, err := getFileContent(input)
	if err != nil {
		log.Fatal(err)
	}

	regex, err := regexp.Compile(getRegexPattern(searchTerm, insensitive))
	if err != nil {
		log.Fatal(err)
	}
	
	scanner := bufio.NewScanner(strings.NewReader(fileContent))
	counter_total := up + down + 1
	counter_tmp := counter_total
	var chunk_lines []t_chunk_line
	var is_printed bool
	match_found := false
	for scanner.Scan() {
		line := scanner.Text()
		match := regex.FindString(line)
		if (len(match)>0)&&(counter_tmp >= 0) {
			match_found = true
			var line_highlighted string
			if useBox{
				line_highlighted = highlightMatchColorBox(line, match, bgColorFlag, fgColorFlag)
			} else {
				line_highlighted = strings.Replace(line, match, highlightMatchColor(match, bgColorFlag, fgColorFlag), -1)
			}
			chunk_lines = append(chunk_lines, t_chunk_line{true, line_highlighted})
			counter_tmp = counter_total
		} else if match_found{
			chunk_lines = append(chunk_lines, t_chunk_line{false, line})
		}
		if (counter_tmp == 0) {
			chunk_lines, is_printed = print_chunk(chunk_lines, up, down)
			if (is_printed && ((up>0) || (down>0))) {
				rgbFg := getColorCode("green")
				highlighted := fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", rgbFg.R, rgbFg.G, rgbFg.B, "------------------------------------------------------------")
				fmt.Println(highlighted)
			}
			counter_tmp = counter_total
		}
		counter_tmp = counter_tmp - 1
	}

	if len(chunk_lines) >= up {
		chunk_lines, _ = print_chunk(chunk_lines, up, down)
	}

	if !match_found{
		fmt.Println("No matches found")
	}
}

type t_chunk_line struct {
	isMatch bool
	line string
}

func print_chunk(chunk []t_chunk_line, up int, down int) ([]t_chunk_line, bool) {
	is_printed := false
	if len(chunk) > 0 {
		if (up == 0)&&(down == 0) {
			for _, item := range chunk {
				fmt.Println(item.line)
			}
			is_printed = true
		} else {
			index_first_match := -1
			for j:=0; j<len(chunk); j++ {
				item := chunk[j]
				if(item.isMatch) {
					index_first_match = j
					break
				}
			}
			index_last_match := -1
			for j:=len(chunk)-1; j>=0; j-- {
				item := chunk[j]
				if(item.isMatch) {
					index_last_match = j
					break
				}
			}
			if (index_first_match >= 0)&&(index_last_match >= 0) {
				var index int
				var max int
				if (index_first_match - up) < 0 {
					index = 0
				} else {
					index = index_first_match - up
				}
				max = index_last_match + down + 1
				if max > len(chunk) {
					max = len(chunk)
				}
				for j:=index; j<max; j++ {
					fmt.Println(chunk[j].line)
				}
				is_printed = true
			}
		}
		// Remove all items in chunk
		return chunk[(len(chunk)-up):], is_printed
	} else {
		return chunk, is_printed
	}
}


//Colors Code//

type RGB struct {
    R int
    G int
    B int
}

func getColorCode(colorOption string) RGB {
    // Return the color code corresponding to the background color option
    switch colorOption {
    case "red":
        return RGB{255, 0, 0}
    case "green":
        return RGB{0, 255, 0}
    case "yellow":
        return RGB{255, 255, 0}
    case "blue":
        return RGB{0, 0, 255}
    case "magenta":
        return RGB{255, 0, 255}
    case "cyan":
        return RGB{0, 255, 255}
    case "white":
        return RGB{255, 255, 255}
	case "black":
		return RGB{0, 0, 0}
	case "pink":
		return RGB{255, 20, 147}
    default:
        return RGB{255, 0, 255} // magenta // default for no box mode
    }
}

func highlightMatchColor(match, bgColor, fgColor string) string {
    var highlighted string

	rgbBg := getColorCode(bgColor)
	if fgColor == "" {
		fgColor = "white"
	}
	rgbFg := getColorCode(fgColor)

	highlighted = fmt.Sprintf("\x1b[48;2;%d;%d;%dm\x1b[38;2;%d;%d;%dm\033[1m%s\x1b[0m", rgbBg.R, rgbBg.G, rgbBg.B, rgbFg.R, rgbFg.G, rgbFg.B, match)

    return highlighted
}

func highlightMatchColorBox(line string, match string, bgColor string, fgColor string) string {

	prev_lines := line[:strings.Index(line, match)]
	post_lines := line[strings.Index(line, match)+len(match):]
	// count the length of the line before the match
	lineLenBeforeMatch := len(prev_lines)
	padding := strings.Repeat(" ", lineLenBeforeMatch)

	// Use red color if bgColor is empty
	if bgColor == "" {
		bgColor = "red"
	}
	if fgColor == "" {
		fgColor = "white"
	}
	rgbBg := getColorCode(bgColor)
	rgbFg := getColorCode(fgColor)

	var highlighted strings.Builder
	highlighted.WriteString(fmt.Sprintf("%s\033[38;2;%d;%d;%dm┏%s┓\033[0m\n", padding, rgbBg.R, rgbBg.G, rgbBg.B, strings.Repeat("━", len(match))))
	highlighted.WriteString(fmt.Sprintf("%s\033[38;2;%d;%d;%dm┃\033[38;2;%d;%d;%dm%s\033[38;2;%d;%d;%dm┃\033[0m%s\n",
	prev_lines, rgbBg.R, rgbBg.G, rgbBg.B,
	rgbFg.R, rgbFg.G, rgbFg.B, match,
	rgbBg.R, rgbBg.G, rgbBg.B, post_lines))
	highlighted.WriteString(fmt.Sprintf("%s\033[38;2;%d;%d;%dm┗%s┛\033[0m", padding, rgbBg.R, rgbBg.G, rgbBg.B, strings.Repeat("━", len(match))))

	return highlighted.String()
}

func contains(colors []string, color string) bool {
	for _, c := range colors {
		if strings.ToLower(c) == strings.ToLower(color) || color == "" {
			return true
		}
	}
	return false
}

//*Colors Code*//

func getInput(args []string) (string, error) {
	if len(args) == 2 {
		return args[1], nil
	} else if !isPipe() {
		return "", fmt.Errorf("No input specified")
	} else {
		return readPipe(), nil
	}
}


func getRegexPattern(searchTerm string, insensitive bool) string {
	if insensitive {
		return "(?i)" + searchTerm
	} else {
		return searchTerm
	}
}

func isPipe() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func readPipe() string {
	var sb strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}
	return sb.String()
}

func getFileContent(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("No file specified")
	}

	stat, err := os.Stat(input)
	if err == nil && !stat.IsDir() {
		bytes, err := ioutil.ReadFile(input)
		if err == nil {
			return string(bytes), nil
		}
	}

	return input, nil
}