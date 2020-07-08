package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "regexp"
)

type Program struct {
    size         int
    instructions []byte
    at           int
}

func check(e error) {
    if e != nil {
        fmt.Println(e)
        panic(e)
    }
}

func execute(code string) {
    var program = new(Program)
    program.size = 30000
    program.instructions = make([]byte, program.size, program.size)
    program.at = 0

    executeWith(program, code)

    fmt.Println()
    fmt.Println("(END)")
}

func executeWith(program *Program, code string) {
    var loopStart = -1
    var loopEnd = -1
    var ignore = 0
    var skipClosingLoop = 0

    for pos, char := range code {
        if ignore == 1 {
            if char == '[' {
                skipClosingLoop += 1
            } else if char == ']' {
                if skipClosingLoop != 0 {
                    skipClosingLoop -= 1
                    continue
                }

                loopEnd = pos
                ignore = 0
                if loopStart == loopEnd {
                    loopStart = -1
                    loopEnd = -1
                    continue
                }
                loop := code[loopStart:loopEnd]
                for program.instructions[program.at] > 0 {
                    executeWith(program, loop)
                }
            }
            continue
        }

        if char == '+' {
            program.instructions[program.at] += 1
        } else if char == '-' {
            program.instructions[program.at] -= 1
        } else if char == '>' {
            if program.at == program.size-1 {
                program.at = 0
            } else {
                program.at += 1
            }
        } else if char == '<' {
            if program.at == 0 {
                program.at = program.size - 1
            } else {
                program.at -= 1
            }
        } else if char == '.' {
            fmt.Printf("%c", rune(program.instructions[program.at]))
        } else if char == ',' {
            fmt.Print("input: ")
            reader := bufio.NewReader(os.Stdin)
            input, _, err := reader.ReadRune()
            check(err)
            program.instructions[program.at] = byte(input)
        } else if char == '[' {
            loopStart = pos + 1
            ignore = 1
        }
    }
}

func main() {
    if len(os.Args) > 1 {
        file, err := ioutil.ReadFile(os.Args[1])
        check(err)

        code := string(file)
        re := regexp.MustCompile(`\r?\n| |[a-zA-Z0-9]`)
        code = re.ReplaceAllString(code, "")

        fmt.Println(code)
        fmt.Println()

        execute(code)
    } else {
        log.Fatal("You must specify a file to execute")
    }
}
