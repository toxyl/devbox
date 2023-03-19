package tar

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func FromDir(srcDir, dstFile string) error {
	src, err := filepath.Abs(srcDir)
	if err != nil {
		return err
	}
	dst, err := filepath.Abs(dstFile)
	if err != nil {
		return err
	}
	ext := filepath.Ext(dst)
	flags := ""
	switch ext {
	case ".xz":
		flags = "-cJf"
	case ".gz":
		flags = "-czf"
	default:
		return fmt.Errorf("unsupported tar type: " + ext)
	}
	err = os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		return err
	}
	err = os.Chdir(src)
	if err != nil {
		return err
	}
	_ = os.Remove(dst) // in case the file already exists
	cmd := exec.Command("tar", flags, dst, "-C", src, "--checkpoint=.100", ".")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Release()
		}
		fmt.Println()
	}()
	return cmd.Run()
}

func ToDir(srcFile, dstDir string) error {
	src, err := filepath.Abs(srcFile)
	if err != nil {
		return err
	}
	dst, err := filepath.Abs(dstDir)
	if err != nil {
		return err
	}
	ext := filepath.Ext(src)
	flags := ""
	switch ext {
	case ".xz":
		flags = "-xSf"
	case ".gz":
		flags = "-xzSf"
	default:
		return fmt.Errorf("unsupported tar type: " + ext)
	}
	_ = os.RemoveAll(dst) // make sure the target directory is gone so we get a clean copy
	err = os.MkdirAll(dst, 0777)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		return err
	}
	os.Chdir(dst)
	if err != nil {
		return err
	}
	cmd := exec.Command("tar", flags, src, "--checkpoint=.100")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Release()
		}
		fmt.Println()
	}()
	return cmd.Run()
}
