package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"github.com/nwaples/rardecode"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func mockUnPack() {

	s1 := "test/unpack/deploy.zip"
	d1 := "test/unpack"
	UnPackage(s1, d1)

	s2 := "test/unpack/example.zip"
	d2 := "test/unpack/example"
	UnPackage(s2, d2)

	s3 := "test/unpack/fin.tar.gz"
	d3 := "test/unpack/fingz"
	UnPackage(s3, d3)

	s4 := "test/unpack/fin.tar"
	d4 := "test/unpack/fintar"
	UnPackage(s4, d4)

	s5 := "test/unpack/IM.rar"
	d5 := "test/unpack/IM"
	UnPackage(s5, d5)
}

func main() {
	if len(os.Args) < 3 {
		println("Usage: unpackage source destination")
		println("Example:")
		println("解压test.rar到E盘我的文件夹")
		println(".\\unpack.exe .\\test.rar  E:\\我的文件夹\\")
		println(".\\unpack.exe .\\test.zip E:\\我的文件夹\\")
		println(".\\unpack.exe .\\test.tar.gz E:\\我的文件夹\\")
		println("解压test.tar到当前目录")
		println(".\\unpack.exe .\\test.tar .")

		return
	}

	sou := os.Args[1]
	des := os.Args[2]

	UnPackage(sou, des)
}

func UnPackage(sou string, des string) {
	if strings.Contains(sou, "zip") {
		unzip(sou, des)
	} else if strings.Contains(sou, "tar") {
		untar(sou, des)
	} else if strings.Contains(sou, "rar") {
		unrar(sou, des)
	} else {
		println("Not support file type")
	}
}

func unrar(sou string, des string) {
	// 打开RAR文件
	file, err := os.Open(sou)
	if err != nil {
		fmt.Printf("打开RAR文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 创建RAR reader
	rr, err := rardecode.NewReader(file, "") // 空字符串表示无密码
	if err != nil {
		fmt.Printf("创建RAR reader失败: %v\n", err)
		return
	}

	// 创建目标目录
	if err := os.MkdirAll(des, 0755); err != nil {
		fmt.Printf("创建目标目录失败: %v\n", err)
		return
	}

	// 遍历RAR文件中的所有文件
	for {
		header, err := rr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("读取RAR文件失败: %v\n", err)
			return
		}

		// 构建目标路径
		path := filepath.Join(des, header.Name)

		if header.IsDir {
			// 创建目录
			if err := os.MkdirAll(path, 0755); err != nil {
				fmt.Printf("创建目录失败: %v\n", err)
				continue
			}
		} else {
			// 创建目标文件的父目录
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				fmt.Printf("创建目录失败: %v\n", err)
				continue
			}

			// 创建目标文件
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				fmt.Printf("创建文件失败: %v\n", err)
				continue
			}

			// 复制文件内容
			if _, err := io.Copy(file, rr); err != nil {
				file.Close()
				fmt.Printf("复制文件内容失败: %v\n", err)
				continue
			}
			file.Close()
		}
	}

	fmt.Println("RAR文件解压完成")
}

func untar(sou string, des string) {
	// 打开tar文件
	file, err := os.Open(sou)
	if err != nil {
		fmt.Printf("打开tar文件失败: %v\n", err)
		return
	}
	defer file.Close()

	var reader *tar.Reader

	// 根据文件扩展名判断是否需要gzip解压
	if strings.HasSuffix(sou, ".gz") || strings.HasSuffix(sou, ".tgz") {
		// 创建gzip reader
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			fmt.Printf("创建gzip reader失败: %v\n", err)
			return
		}
		defer gzipReader.Close()
		reader = tar.NewReader(gzipReader)
	} else {
		reader = tar.NewReader(file)
	}

	// 创建目标目录
	if err := os.MkdirAll(des, 0755); err != nil {
		fmt.Printf("创建目标目录失败: %v\n", err)
		return
	}

	// 遍历tar文件中的所有文件和目录
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break // 到达文件末尾
		}
		if err != nil {
			fmt.Printf("读取tar文件失败: %v\n", err)
			return
		}

		// 构建目标路径
		path := filepath.Join(des, header.Name)

		switch header.Typeflag {
		case tar.TypeDir: // 目录
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				fmt.Printf("创建目录失败: %v\n", err)
				continue
			}
		case tar.TypeReg: // 普通文件
			// 创建目标文件的父目录
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				fmt.Printf("创建目录失败: %v\n", err)
				continue
			}

			// 创建目标文件
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				fmt.Printf("创建文件失败: %v\n", err)
				continue
			}

			// 复制文件内容
			if _, err := io.Copy(file, reader); err != nil {
				file.Close()
				fmt.Printf("复制文件内容失败: %v\n", err)
				continue
			}
			file.Close()
		}
	}
}

func unzip(sou string, des string) {
	// 打开zip文件
	reader, err := zip.OpenReader(sou)
	if err != nil {
		fmt.Printf("打开zip文件失败: %v\n", err)
		return
	}
	defer reader.Close()

	// 创建目标目录
	if err := os.MkdirAll(des, 0755); err != nil {
		fmt.Printf("创建目标目录失败: %v\n", err)
		return
	}

	// 遍历zip文件中的所有文件和目录
	for _, file := range reader.File {
		// 构建目标路径
		path := filepath.Join(des, file.Name)

		// 如果是目录，创建它
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		// 创建目标文件的父目录
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			fmt.Printf("创建目录失败: %v\n", err)
			continue
		}

		// 创建目标文件
		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fmt.Printf("创建文件失败: %v\n", err)
			continue
		}

		// 打开zip中的源文件
		srcFile, err := file.Open()
		if err != nil {
			dstFile.Close()
			fmt.Printf("打开zip中的文件失败: %v\n", err)
			continue
		}

		// 复制文件内容
		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()

		if err != nil {
			fmt.Printf("复制文件内容失败: %v\n", err)
		}
	}
}
