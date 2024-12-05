package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Plugin struct {
	workingDir string

	StoragePath   string
	Src           string
	Dockerignores []string

	/*
		Dest     string
		DestType string
	*/
}

// Exec executes the plugin step
func (p Plugin) Exec() error {
	workingDir, _ := os.Getwd()
	p.workingDir = workingDir
	fmt.Printf("Working dir: %s\n", workingDir)
	if p.Src != "" {
		err := p.copySrc2StoragePath()
		if err != nil {
			return fmt.Errorf("copy src to storage path error: %s", err)
		}
		err = p.handleDockerignoreFile()
		if err != nil {
			return fmt.Errorf("handle .dockerignore error: %s", err)
		}
		return nil
	}
	/*
		if p.Dest != "" {
			return doExec("cp", "-rf", abs(p.StoragePath), abs(p.Dest))
		}
		fmt.Println("src and desc both are blank, skip execute.")
	*/
	fmt.Println("src is blank, skip execute.")
	return nil
}

// copySrc2StoragePath 复制src到storagePath
func (p Plugin) copySrc2StoragePath() error {
	err := doExec("cp", "-rf", abs(p.Src)+string(filepath.Separator)+".", abs(p.StoragePath))
	if err != nil {
		return err
	}
	repositoriesFilePath := filepath.Join(abs(p.StoragePath), "image/overlay2/repositories.json")
	_, err = os.Stat(repositoriesFilePath)
	if os.IsNotExist(err) {
		fmt.Printf("\u001B[31m%s is not exists.\u001B[0m\n", repositoriesFilePath)
		return nil
	}
	err = doExec("jq", ".", repositoriesFilePath)
	if err != nil {
		return doExec("cat", repositoriesFilePath)
	}
	return nil
}

// handleDockerignoreFile 处理.dockerignore文件
func (p Plugin) handleDockerignoreFile() error {
	// 1. 处理忽略的文件
	ignores := []string{}
	// 1.1 若storage-path在工作目录下，忽略它
	if !filepath.IsAbs(p.StoragePath) {
		ignores = append(ignores, p.StoragePath)
	} else if strings.HasPrefix(p.StoragePath, p.workingDir) {
		rel, _ := filepath.Rel(p.workingDir, abs(p.StoragePath))
		ignores = append(ignores, rel)
	}
	// 1.2 额外忽略路径
	for _, ignore := range p.Dockerignores {
		ignores = append(ignores, strings.Split(ignore, ",")...)
	}
	if len(ignores) <= 0 {
		return nil
	}

	// 2. 写入忽略路径
	var file *os.File
	// 2.1 判断文件是否存在, 不存在则创建
	fp := filepath.Join(p.workingDir, ".dockerignore")
	_, err := os.Stat(fp)
	if os.IsNotExist(err) {
		file, err = os.Create(fp)
		if err != nil {
			return err
		}
		fmt.Printf("+ create file: %s\n", fp)
	}
	// 2.2 以追加模式打开文件
	file, err = os.OpenFile(fp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// 2.3 将文本写入文件
	_, err = file.WriteString("\n# --- plugin auto append begin ---\n")
	_, err = file.WriteString(strings.Join(ignores, "\n"))
	_, err = file.WriteString("\n# --- plugin auto append end ---")
	// 2.4 打印
	doExec("cat", fp)
	return err
}

func abs(s string) string {
	a, _ := filepath.Abs(s)
	return a
}

func doExec(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	return cmd.Run()
}

func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
