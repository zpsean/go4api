/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package keyword

import (
    // "fmt"
    // "time"
    // "os"
    // "sort"
)


func GetContent (kwFileList []string) KWBlocks {
    var lines []string

    // to be replaced
    kwFileList, _ := WalkPath("/Users/pingzhu/Downloads/goState/testdata", ".keyword")

    for _, jsonFile := range kwFileList {
        lines, _ = readLines(jsonFile)
    }

    kwBlocks := GetBlocks(lines)
    kwBlocks = FullfillBlocks(kwBlocks, lines)

    // fmt.Println("mdBlocks: ", mdBlocks, len(mdBlocks))

    return kwBlocks
}

func FullfillBlocks (kwBlocks KWBlocks, lines []string) KWBlocks {
    for index, _ := range kwBlocks {
        // set block OriginalContent
        for i := kwBlocks[index].StartLine; i <= kwBlocks[index].EndLine; i++ {
            kwBlocks[index].OriginalContent = append(kwBlocks[index].OriginalContent, lines[i])
        }

        // set block BlockType
        kwBlocks[index].BlockType = GetBlockType(kwBlocks[index])

        // set block ParsedContent
        switch kwBlocks[index].BlockType {
        case "TestCases":
            //
        case "Settings":
            //
        case "Keywords":
            //
        case "Variables":
            //
        default:
            fmt.Println("Warning, can not recognize the block type")
        }
    }

    return kwBlocks
}

func GetBlocks (lines []string) KWBlocks {
    // Note: each block has the leading line with prefix '*** TestCases / Settings / Keywords / Variables /...''
    var blockHeaderLines []int

    linesCount := len(lines)

    // get the block header line numbers, starting from line 0
    for i, line := range lines {
        if strings.HasPrefix(strings.TrimSpace(line), "***") {
            blockHeaderLines = append(blockHeaderLines, i)
        }
    }

    var kwBlocks KWBlocks
    var kwBlock  KWBlock

    headerCount := len(blockHeaderLines)

    for i, _ := range blockHeaderLines {
        if i != headerCount - 1 {
            kwBlock = &kwBlock {
                StartLine: blockHeaderLines[i],
                EndLine: blockHeaderLines[i + 1] - 1,
            }
        } else {
            kwBlock = &kwBlock {
                StartLine: blockHeaderLines[i],
                EndLine: linesCount - 1,
            }
        }

        kwBlocks = append(kwBlocks, kwBlock)
    }

    return kwBlocks
}

func GetBlockType (kwBlock KWBlock) string {
    var blockType string

    blockTypes := []string{"TestCases", "Settings", "Keywords", "Variables"}

    for i, _ := range blockTypes {
        strings.Count(kwBlock.OriginalContent[0], blockTypes[i]) > 1 {
            blockType = blockTypes[i]
            break
        }
    }

    return blockType
}

func readLines (path string) (lines []string, err error){  
    f, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer f.Close()
 
    rd := bufio.NewReader(f)
    for {
            line, err := rd.ReadString('\n')

            line = strings.Replace(line, "\n", "", -1)
            lines = append(lines, line)

            // fmt.Println(line)
          
            if err != nil || io.EOF == err {
                break
            }  
        }

    return
}  
