package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mattn/go-runewidth"
)

type Word struct {
	id         int
	word       string
	meaning    string
	importance string
	addedAt    string
}
type Item struct {
	id         int
	entry      string
	meaning    string
	importance string
	cycle      int
	state      int
}

func main() {
	// 设置数据库连接信息
	user, password, host, database, err := readConfig("~/.config/voc/voc.conf")
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) < 2 {
		log.Fatal("too few arguments\n")
	}
	operation := os.Args[1]
	/*
		user := "root"
		password := ""
		host := "localhost"
		database := "vocabulary"
	*/
	// 构建数据源名称 (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", user, password, host, database)

	// 连接到数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// 检查数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	fmt.Println("Successfully connected to the database")
	fmt.Println("")
	if operation == "-h" {
		fmt.Printf("          WELCOME TO USE terminalWordList!      \n")
		fmt.Printf("insert a word:      voc -iw word meaning importance\n\n")
		fmt.Printf("insert a phrase:    voc -ip phrase meaning importance\n\n")
		fmt.Printf("update importance:  voc -u table id importance\n\n")
		fmt.Printf("query a word:       voc -q word\n\n")
		fmt.Printf("show words:         voc -sw\n\n")
		fmt.Printf("show phrases:       voc -sp\n\n")
		fmt.Printf("***show order by importance: voc -sw/-sp 1***\n\n")
		fmt.Printf("delete entries:     voc -d table id offset\n\n")
		fmt.Printf("review entries:     voc -r table startIndex endIndex\n\n")
		fmt.Printf("recall entries:     voc -rc table startIndex endIndex\n\n")
		fmt.Printf("help:               voc -h\n\n")
	}
	if operation == "-iw" {
		if len(os.Args) != 5 {
			fmt.Println("Usage: voc -iw word meaning importance")
			return
		}
		word := os.Args[2]
		meaning := os.Args[3]
		importanceStr := os.Args[4]
		importance, err := strconv.Atoi(importanceStr)
		if err != nil {
			log.Fatal("Failed to convert importance to integer:", err)
		}
		insertedID, err := insertEntry(db, importance, "words", word, meaning)
		if err != nil {
			log.Fatal("Failed to insert entry:", err)
		}
		fmt.Printf("Successfully inserted entry with id: %d", insertedID)
	}

	if operation == "-ip" {
		if len(os.Args) != 5 {
			fmt.Println("Usage: voc -ip phrase meaning importance")
			return
		}
		phrase := os.Args[2]
		meaning := os.Args[3]
		importanceStr := os.Args[4]
		importance, err := strconv.Atoi(importanceStr)
		if err != nil {
			log.Fatal("Failed to convert importance to integer:", err)
		}
		insertedID, err := insertEntry(db, importance, "phrases", phrase, meaning)
		if err != nil {
			log.Fatal("Failed to insert entry:", err)
		}
		fmt.Printf("Successfully inserted entry with id: %d", insertedID)
	}
	if operation == "-u" {
		if len(os.Args) != 5 {
			fmt.Println("usage: voc -u words/phrases id importance")
			return
		}
		table := os.Args[2]
		idStr := os.Args[3]
		importanceStr := os.Args[4]
		importance, er := strconv.Atoi(importanceStr)
		if er != nil {
			log.Fatal("Failed to convert importance to integer:", er)
		}
		id, e := strconv.Atoi(idStr)
		if e != nil {
			log.Fatal("Failed to convert ID to integer:", e)
		}
		err := updateImportance(db, table, id, importance)
		if err != nil {
			log.Fatal("Failed to update importance:", err)
		}
		fmt.Println("Successfully updated importance")
	}
	if operation == "-r" {
		if len(os.Args) != 5 {
			fmt.Println("usage: voc -r words/phrases startIndex endIndex")
			return
		}
		table := os.Args[2]
		startStr := os.Args[3]
		start, err := strconv.Atoi(startStr)
		if err != nil {
			log.Fatal("Failed to convert start index to integer:", err)
		}
		endStr := os.Args[4]
		end, err := strconv.Atoi(endStr)
		if err != nil {
			log.Fatal("Failed to convert end index to integer:", err)
		}
		err = review(db, table, start, end)
		if err != nil {
			log.Fatal("Failed to review:", err)
		}
		fmt.Println("Successfully reviewed")
	}

	if operation == "-rc" {
		if len(os.Args) != 5 {
			fmt.Println("usage: voc -rc words/phrases startIndex endIndex")
			return
		}
		table := os.Args[2]
		startStr := os.Args[3]
		start, err := strconv.Atoi(startStr)
		if err != nil {
			log.Fatal("Failed to convert start index to integer:", err)
		}
		endStr := os.Args[4]
		end, err := strconv.Atoi(endStr)
		if err != nil {
			log.Fatal("Failed to convert end index to integer:", err)
		}
		err = recall(db, table, start, end)
		if err != nil {
			log.Fatal("Failed to recall:", err)
		}
		fmt.Println("Successfully recalled")
	}
	if operation == "-q" {
		word := os.Args[2]
		meaning, importance, err := queryWords(db, word)
		if err != nil {
			log.Fatal("Failed to query word:", err)
		}
		fmt.Printf("%s的意思是: %s\n", word, meaning)
		fmt.Printf("重要性: %d\n", importance)
	}
	if operation == "-sw" {
		if len(os.Args) == 2 {
			err := show(db, "words", 0)
			if err != nil {
				log.Fatal("Failed to show words:", err)
			}
		} else {
			num, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatal("Failed to convert num to integer:", err)
			}
			var showErr error
			if num == 1 {
				showErr = show(db, "words", 1)
			} else {
				showErr = show(db, "words", 0)
			}

			if showErr != nil {
				log.Fatal("Failed to show words:", showErr)
				// handle error
			}
		}
	}
	if operation == "-sp" {

		if len(os.Args) == 2 {
			err := show(db, "phrases", 0)
			if err != nil {
				log.Fatal("Failed to show phrases:", err)
			}
		} else {
			num, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatal("Failed to convert num to integer:", err)
			}
			var showErr error
			if num == 1 {
				showErr = show(db, "phrases", 1)
			} else {
				showErr = show(db, "phrases", 0)
			}

			if showErr != nil {
				log.Fatal("Failed to show phrases:", showErr)
				// handle error
			}
		}
	}
	if operation == "-d" {
		table := os.Args[2]
		idStr := os.Args[3]
		offStr := os.Args[4]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal("Failed to convert ID to integer:", err)
		}
		off, err := strconv.Atoi(offStr)
		if err != nil {
			log.Fatal("Failed to convert offset to integer:", err)
		}
		err = delete(db, table, id, off)
		if err != nil {
			log.Fatal("Failed to delete entries:", err)
		}
	}

}

func insertEntry(db *sql.DB, importance int, table, insertion, meaning string) (int64, error) {
	replace := "entry"
	if table == "words" {
		replace = "word"
	}
	var minID *int64
	// 查询最小可用的ID
	err := db.QueryRow("SELECT MIN(id+1) as nextID FROM " + table + " t WHERE NOT EXISTS (SELECT 1 FROM " + table + " WHERE id = t.id + 1)").Scan(&minID)
	if err != nil {
		return 0, err
	}
	if minID == nil {
		tmp := int64(1)
		minID = &tmp
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := fmt.Sprintf("INSERT INTO %s (id, %s, meaning, importance,added_at) VALUES (?, ?, ?, ?,?)", table, replace)
	result, err := db.Exec(query, minID, insertion, meaning, importance, currentTime)

	if err != nil {
		return 0, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func queryWords(db *sql.DB, word string) (string, int, error) {
	var meaning string
	var importance int
	err := db.QueryRow("SELECT meaning, importance FROM words WHERE word = ?", word).Scan(&meaning, &importance)

	if err != nil {
		return "", -1, err
	}

	return meaning, importance, nil
}

func updateImportance(db *sql.DB, table string, id int, importance int) error {
	query := fmt.Sprintf("UPDATE %s SET importance = ? WHERE id = ?", table)
	_, err := db.Exec(query, importance, id)
	if err != nil {
		return err
	}
	return nil
}

func show(db *sql.DB, table string, opt int) error {
	// Prepare the query string
	var query string
	if opt == 0 {
		query = fmt.Sprintf("SELECT * from %s", table)
	} else {
		query = fmt.Sprintf("SELECT * FROM %s order by importance desc", table)
	}

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Store the entries in a slice to calculate max column widths
	var items []Item
	for rows.Next() {
		var id int
		var entry, meaning, addedAt string
		var importance int
		var cycles int
		var state int
		err = rows.Scan(&id, &entry, &meaning, &importance, &addedAt, &cycles, &state)
		if err != nil {
			return err
		}

		//cycles, err = strconv.Atoi(cycles)
		item := Item{
			id:         id,
			entry:      entry,
			meaning:    meaning,
			importance: fmt.Sprintf("%d", importance),
			cycle:      cycles,
			state:      state,
		}
		items = append(items, item)
	}

	// Calculate the max width for each column
	maxID, maxWord, maxMeaning, maxImportance := 2, 10, 7, 10
	for _, item := range items {
		idLen := len(strconv.Itoa(item.id))
		wordLen := widthOfString(item.entry)
		meaningLen := widthOfString(item.meaning)
		importanceLen := len(item.importance)

		if idLen > maxID {
			maxID = idLen
		}
		if wordLen > maxWord {
			maxWord = wordLen
		}
		if meaningLen > maxMeaning {
			maxMeaning = meaningLen
		}
		if importanceLen > maxImportance {
			maxImportance = importanceLen
		}
	}

	// Print the table header with specified widths

	fmt.Println(strings.Repeat("-", maxID+maxWord+maxMeaning+maxImportance+29))
	headerFormat := fmt.Sprintf("| %%-%ds | %%-%ds | %%-%ds | %%-%ds | cycle | state |\n", maxID, maxWord, maxMeaning, maxImportance)
	fmt.Printf(headerFormat, "ID", "Word/Phrase", "Meaning", "Importance")

	// Print the separator line
	fmt.Println(strings.Repeat("-", maxID+maxWord+maxMeaning+maxImportance+29))
	// Iterate through the entries and print the data
	for _, item := range items {
		// Parse the addedAt string into a time.Time value

		rowFormat := fmt.Sprintf("| %%-%dd | %%-%ds | %%-%ds | %%-%ds | %%-5d | %%-5d |\n", maxID, maxWord, maxMeaning-numberOfChinese(item.meaning), maxImportance)
		fmt.Printf(rowFormat, item.id, item.entry, item.meaning, item.importance, item.cycle, item.state)
	}

	fmt.Println(strings.Repeat("-", maxID+maxWord+maxMeaning+maxImportance+29))
	return nil
}

func delete(db *sql.DB, table string, id int, off int) error {
	// Prepare the delete statement
	query := fmt.Sprintf("DELETE FROM %s WHERE id >= ? AND id <= ?", table)

	// Execute the delete statement
	_, err := db.Exec(query, id, id+off)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted entries with ID between %d and %d in table %s.\n", id, id+off, table)
	return nil
}

func readConfig(filename string) (string, string, string, string, error) {
	if strings.HasPrefix(filename, "~") {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		filename = strings.Replace(filename, "~", usr.HomeDir, 1)
	}
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", "", "", "", fmt.Errorf("Failed to read the config file: %v", err)
	}

	lines := strings.Split(string(config), "\n")

	var user, password, host, database string
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		keyValue := strings.Split(line, ":")
		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		switch key {
		case "user":
			user = value
		case "password":
			password = value
		case "host":
			host = value
		case "database":
			database = value
		default:
			return "", "", "", "", fmt.Errorf("Unknown key in config file: %s", key)
		}
	}

	return user, password, host, database, nil
}

func getStringWidth(s string) int {
	return runewidth.StringWidth(s)
}

func widthOfString(s string) int {
	width := 0
	for _, r := range s {
		if r < 128 {
			width += 1
		} else {
			width += 2
		}
	}
	return width
}

func numberOfChinese(s string) int {
	width := 0
	for _, r := range s {
		if r >= 128 {
			width += 1
		}
	}
	return width
}

func findCycleIndex(timePassed int) int {
	// 时间区间（单位：分钟）
	timeIntervals := []int{5, 30, 24 * 60, 48 * 60, 96 * 60, 168 * 60, 360 * 60, 744 * 60}

	// 使用二分法查找时间区间
	left, right := 0, len(timeIntervals)
	for left < right {
		mid := left + (right-left)/2
		if timePassed > timeIntervals[mid] {
			left = mid + 1
		} else {
			right = mid
		}
	}

	return left
}

func review(db *sql.DB, table string, start int, end int) error {
	// 查询 id 介于 start 和 end 之间的行
	rows, err := db.Query(fmt.Sprintf("SELECT id, added_at, cycle, state FROM %s WHERE id BETWEEN ? AND ?", table), start, end)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		var id int
		var addedAt string
		var cycle int
		var state int

		err := rows.Scan(&id, &addedAt, &cycle, &state)
		if err != nil {
			return err
		}

		// 计算 added_at 到现在过去的分钟数
		addedAtTime, err := time.Parse("2006-01-02 15:04:05", addedAt)
		if err != nil {
			return err
		}
		timePassed := int(time.Since(addedAtTime).Minutes())

		// 调用 findCycleIndex 函数获取 index
		index := findCycleIndex(timePassed)
		_, err = db.Exec(fmt.Sprintf("UPDATE %s SET cycle = cycle + 1  WHERE id = ?", table), id)
		if err != nil {
			return err
		}
		// 如果 cycle > state，更新 state 列为 index
		if cycle+1 > state && index > state {
			_, err := db.Exec(fmt.Sprintf("UPDATE %s SET state = state + 1 WHERE id = ?", table), id)
			if err != nil {
				return err
			}
		}
	}

	// 检查遍历过程中是否发生错误
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func recall(db *sql.DB, table string, start int, end int) error {
	query := fmt.Sprintf("UPDATE %s SET cycle = cycle - 1  WHERE id between ? and ?", table)
	_, err := db.Exec(query, start, end)
	if err != nil {
		return err
	}
	return nil
}
