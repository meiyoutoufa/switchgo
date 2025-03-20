package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var switchGoCmd = &cobra.Command{
	Use:   "to",
	Short: "切换到新的版本",
	Run: func(cmd *cobra.Command, args []string) {
		// 从命令行参数读取环境变量名
		version, _ := cmd.Flags().GetString("version")
		goPath, _ := cmd.Flags().GetString("gopath")
		if goPath == "" {
			goPath = os.Getenv("GOPATH")
		}

		// 检查是否存在 PATH 环境变量
		if path, exists := os.LookupEnv("PATH"); !exists {
			fmt.Println("PATH 未设置")
			return
		} else {
			paths := strings.Split(path, ";")

			for i, path := range paths {
				if strings.Contains(path, goPath) {
					path := buildPath(goPath, version)
					if !checkExistDir(path) {
						fmt.Println("没有这个go版本的包, 请确认")
						return
					}
					paths[i] = path

				}
			}
			newPath := strings.Join(paths, ";")
			setPath(newPath)
		}
	},
}

func setPath(newPath string) {
	// 根据操作系统调用不同方法
	switch runtime.GOOS {
	case "windows":
		// 用户级生效
		if err := setPathWindows(newPath, true); err != nil {
			fmt.Println("修改失败:", err)
		} else {
			fmt.Println("PATH 已更新，请重启终端或重新登录生效")
		}

	case "linux", "darwin":
		// 用户级生效（系统级需设置 systemLevel=true 并用 sudo 运行）
		// if err := setPathUnix(newPath, false); err != nil {
		// 	fmt.Println("修改失败:", err)
		// } else {
		// 	fmt.Println("PATH 已更新，请执行 source ~/.bashrc 或重启终端生效")
		// }

	default:
		fmt.Println("不支持的操作系统")
	}
}

func setPathUnix(newPath string, systemLevel bool) error {
	configFile := ""
	if systemLevel {
		configFile = "/etc/environment" // 系统级（需root权限）
	} else {
		// 获取当前用户的默认Shell配置文件
		homeDir, _ := os.UserHomeDir()
		configFile = fmt.Sprintf("%s/.bashrc", homeDir) // 可扩展为 .zshrc 等
	}

	// 追加PATH到配置文件
	line := fmt.Sprintf("\nexport PATH=\"%s:$PATH\"\n", newPath)
	return os.WriteFile(configFile, []byte(line), 0644)
}

func setPathWindows(newPath string, systemLevel bool) error {
	// 构建 setx 命令
	args := []string{"/c", "setx", "PATH", newPath}
	if systemLevel {
		args = append(args, "/M") // 修改系统级变量（需管理员权限）
	}

	cmd := exec.Command("cmd", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	return cmd.Run()
}

func buildPath(goPath string, vsersion string) string {
	// 获取操作系统
	os := runtime.GOOS
	// 根据不同系统生成路径
	var configPath string
	switch os {
	case "windows":
		configPath = filepath.Join(goPath, "go"+vsersion, "bin")
	case "darwin":
		configPath = filepath.Join(goPath, "go"+vsersion, "bin")
	case "linux":
		configPath = filepath.Join(goPath, "go"+vsersion, "bin")
	default:
		configPath = filepath.Join(goPath, "go"+vsersion, "bin")
	}
	return configPath
}

func checkExistDir(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("%s 不存在\n", path)
		return false
	}
	if err != nil {
		fmt.Printf("检查 %s 出错: %v\n", path, err)
		return false
	}

	// if info.IsDir() {
	//     fmt.Printf("%s 是目录\n", path)
	// } else {
	//     fmt.Printf("%s 是文件，大小: %d 字节\n", path, info.Size())
	// }
	return true
}

func init() {
	switchGoCmd.Flags().StringP("version", "v", "1.18.10", "切换go的版本")
	switchGoCmd.Flags().StringP("gopath", "p", "", "GOPATH的地址")
}
