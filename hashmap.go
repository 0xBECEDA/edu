package main

import(
    "fmt"
)

type hash_map_cell struct {
    Age int
    Name string
    NextCell *hash_map_cell
}

var table_cap = 10
var table = make( []*hash_map_cell, table_cap)

func hash_key(key int) int {
    return key % table_cap
}

// проходится по списку значений и либо добавляет ячейку,
// либо изменяет значение ячейки, чей ключ соответствует заданному
func addValueToList( list *hash_map_cell, key int, value string ) {

    if list.Age == key {
        list.Name = value
        return

    } else {
        if list.NextCell == nil {
            newCell := hash_map_cell{Age: key, Name: value, NextCell: nil}
            list.NextCell = &newCell
            return
        }
        addValueToList( list.NextCell, key, value )
    }
    return
}

// добавляет значение в таблицу
// или изменяет уже существующее, если заданный ключ присутствует
func addValueToTable(key int, value string) {
    indx := hash_key(key)
    hashCell := table[indx]

    if	hashCell == nil {
        newCell := hash_map_cell{Age: key, Name: value, NextCell: nil}
        table[indx] = &newCell
        return
    }

    addValueToList(hashCell, key, value)
    table[indx] = hashCell
    return
}

// ищет значение в списке
func foundValueInList(list *hash_map_cell, key int) (string, int) {

    if list.Age == key {
        return list.Name, 0

    } else {
        if list.NextCell == nil {
            return "No value found \n", -1
        }
        return foundValueInList(list.NextCell, key)
    }
}

// ищет значение в таблице
// возвращает значение или сообщение об ошибке
func getValueFromTable(key int) (string, int)  {

    indx := hash_key(key)
    hashCell := table[indx]

    if hashCell == nil {
        return "No value found \n", -1
    }
    return foundValueInList(hashCell, key)
}

// удаляет ячейку из списка
func deleteValueFromList(list *hash_map_cell, key int) (string, int) {

    if list == nil {
        return "No such value \n", -1

    } else if list.NextCell.Age == key {
        list.NextCell = list.NextCell.NextCell
        return "Value successfully deleted \n ", 0
    }
    return deleteValueFromList(list.NextCell, key)
}

// удаляет значение из мапы по ключу
func deleteValue(key int) (string, int) {

    indx := hash_key(key)
    hashCell := table[indx]

    if hashCell == nil {
        return "No such value \n", -1

    } else if hashCell.Age == key {
        table[indx] = hashCell.NextCell
        return "Value successfully deleted \n ", 0
    }
    return deleteValueFromList(hashCell, key)
}

func main () {
    // тестируем добавление значений, хэш ключей которых
    // получается одинаковым
    addValueToTable(32, "Pushkin" )
    val, _ := getValueFromTable(32)
    fmt.Printf(" getValueFromTable: key is %d, value is %s \n", 32, val)

    addValueToTable(42, "Leonid" )
    val2, _ := getValueFromTable(42)
    fmt.Printf(" getValueFromTable: key is %d, value is %s \n", 42, val2)

    addValueToTable(12, "Ivan" )
    val3, _ := getValueFromTable(12)
    fmt.Printf(" getValueFromTable: key is %d, value is %s \n", 12, val3)

    // добавляем значение с ключом, хэш которого будет отличным от предыдущих
    addValueToTable(10, "Stefan" )
    val4, _ := getValueFromTable(10)
    fmt.Printf(" getValueFromTable: key is %d, value is %s \n", 10, val4)

    // удаляем значение из середины цепочки
    msg, _ := deleteValue(42)
    fmt.Printf(" deleteValue: key is %d, msg: %s \n", 42, msg)

    // проверяем, что значение действительно удалено
    val2, _ = getValueFromTable(42)
    fmt.Printf(" getValueFromTable: key is %d, value is %s \n", 42, val2)
}
