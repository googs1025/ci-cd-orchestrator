package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// Scanner 文件扫描器接口
type Scanner interface {
	Scan(projectPath string) ([]string, error)
	IsRelevantFile(filePath string) bool
}

// FileScanner 文件扫描器实现
type FileScanner struct {
	// 要忽略的目录
	IgnoredDirs []string
	// 要忽略的文件
	IgnoredFiles []string
	// 要包含的文件扩展名
	RelevantExtensions []string
}

// NewFileScanner 创建文件扫描器实例
func NewFileScanner() Scanner {
	return &FileScanner{
		IgnoredDirs: []string{
			".git",
			".svn",
			".hg",
			"node_modules",
			"vendor",
			"dist",
			"build",
			"bin",
			"obj",
			"target",
		},
		IgnoredFiles: []string{
			".DS_Store",
		},
		RelevantExtensions: []string{
			".js",
			".jsx",
			".ts",
			".tsx",
			".go",
			".java",
			".py",
			".rs",
			".cpp",
			".c",
			".cs",
			".html",
			".css",
			".scss",
			".json",
			".yaml",
			".yml",
			".xml",
			".gradle",
			".mvn",
			".gitignore",
			".dockerignore",
			"Makefile",
			"README.md",
			"package.json",
			"go.mod",
			"requirements.txt",
			"pom.xml",
			"build.gradle",
			"Cargo.toml",
			"setup.py",
		},
	}
}

// Scan 扫描项目目录
func (s *FileScanner) Scan(projectPath string) ([]string, error) {
	var relevantFiles []string

	// 递归扫描目录
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否是目录
		if info.IsDir() {
			// 检查是否是要忽略的目录
			dirName := filepath.Base(path)
			for _, ignoredDir := range s.IgnoredDirs {
				if dirName == ignoredDir {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// 检查是否是要忽略的文件
		fileName := filepath.Base(path)
		for _, ignoredFile := range s.IgnoredFiles {
			if fileName == ignoredFile {
				return nil
			}
		}

		// 检查是否是相关文件
		if s.IsRelevantFile(path) {
			relevantFiles = append(relevantFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return relevantFiles, nil
}

// IsRelevantFile 检查文件是否与技术栈识别相关
func (s *FileScanner) IsRelevantFile(filePath string) bool {
	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, relevantExt := range s.RelevantExtensions {
		if ext == relevantExt {
			return true
		}
	}

	// 检查文件名
	fileName := strings.ToLower(filepath.Base(filePath))
	for _, relevantExt := range s.RelevantExtensions {
		if !strings.HasPrefix(relevantExt, ".") && fileName == relevantExt {
			return true
		}
	}

	return false
}
