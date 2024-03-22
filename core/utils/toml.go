package utils

import "os"

// func CreateToml(tomlPath string) {
// 	tree, err := toml.Load("")
// 	if err != nil {
// 		fmt.Println("Error while creating empty Toml tree:", err)
// 		return
// 	}
// 	subtree, err := toml.Load("")
// 	subtree.Set("key1", "value1 - %s")
// 	subtree.Set("key2", "42")
// 	subtree.Set("key3", true)

// 	tree.Set("report_config", subtree)

// 	file, err := os.Create(tomlPath)
// 	if err != nil {
// 		fmt.Println("Error while creating Toml file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	_, err = tree.WriteTo(file)
// 	if err != nil {
// 		fmt.Println("Error while writing to Toml file:", err)
// 		return
// 	}

// 	fmt.Println("Toml file created successfully.")
// }

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
