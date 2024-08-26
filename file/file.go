package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func allExcel() (all []string) {
	// 指定要遍历的文件夹路径
	folderPath := "./excel"

	// 遍历文件夹
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		// 检查是否为文件夹
		if err != nil {
			fmt.Println(err)
			return nil
		}

		// 检查是否为xlsx文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".xlsx") {
			// 获取包含路径的文件名
			all = append(all, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return
}
