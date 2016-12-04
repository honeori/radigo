package radigo

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/yyoshiki41/radigo/internal"
)

var aacResultFile string

func initTempAACDir() (string, error) {
	aacDir, err := ioutil.TempDir(radigoPath, "aac")
	if err != nil {
		return "", err
	}

	aacResultFile = path.Join(aacDir, "result.aac")
	return aacDir, nil
}

func outputMP3(aacDir, outputFile string) error {
	if err := createConcatedAACFile(aacDir); err != nil {
		return err
	}

	if err := convertAACToMP3(outputFile); err != nil {
		return err
	}

	return nil
}

func createConcatedAACFile(aacDir string) error {
	name, err := internal.ConcatFileNames(aacDir)
	if err != nil {
		return err
	}

	f, err := newFfmpeg(fmt.Sprintf("concat:%s", name))
	if err != nil {
		return err
	}

	f.setDir(aacDir)
	f.setArgs("-c", "copy")
	// TODO: Run 結果の標準出力を拾う
	return f.run(aacResultFile)
}

func convertAACToMP3(outputFile string) error {
	f, err := newFfmpeg(aacResultFile)
	if err != nil {
		return err
	}

	f.setDir(radigoPath)
	f.setArgs(
		"-c:a", "libmp3lame",
		"-ac", "2",
		"-q:a", "2",
	)
	// TODO: Run 結果の標準出力を拾う
	return f.run(outputFile)
}