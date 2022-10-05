// Copyright 2018 Shannon Pekary. All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

package sys

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeTempDir() (dir string, err error) {
	dir, err = os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	dir = filepath.Join(dir, "gofileTest")
	err = os.Mkdir(dir, 0777)
	if err != nil {
		_ = os.RemoveAll(dir)
		return
	}
	return
}

func testDataDir1() string {
	return filepath.Join("testdata", "dir1")
}
func testDataDir2() string {
	return filepath.Join("testdata", "dir2")
}

func TestCopyFile(t *testing.T) {
	tempDir, err := makeTempDir()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.RemoveAll(tempDir)

	src := filepath.Join(tempDir, "test.txt")

	testContent := "I am a test \n and the 2nd line"

	err = os.WriteFile(src, []byte(testContent), 0755)
	if err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(tempDir, "test2.txt")

	time.Sleep(2 * time.Second) // delay to make sure file is created later

	err = CopyFiles(dst, CopyDoNotOverwrite, src)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := os.ReadFile(dst)
	if string(bytes) != testContent {
		t.Error("File content does not match")
	}

	time.Sleep(2 * time.Second) // delay to make sure file is created later

	// test grow
	testContent2 := "I am a test \n and the 2nd line\nand a 3rd"
	src2 := filepath.Join(tempDir, "test2-2.txt")
	err = os.WriteFile(src2, []byte(testContent2), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// We should not overwrite a file that already exists
	err = CopyFiles(dst, CopyDoNotOverwrite, src2)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err = os.ReadFile(dst)
	if string(bytes) != testContent {
		t.Error("File content should not have changed")
	}

	// We should overwrite a file that already exists here
	err = CopyFiles(dst, CopyOverwriteOnlyIfNewer, src2)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err = os.ReadFile(dst)
	if string(bytes) != testContent2 {
		t.Error("File should overwrite")
	}

	// We should not overwrite a newer file with an older file here
	err = CopyFiles(dst, CopyOverwriteOnlyIfNewer, src)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err = os.ReadFile(dst)
	if string(bytes) != testContent2 {
		t.Error("File should not overwrite")
	}

	// We should overwrite the file here
	err = CopyFiles(dst, CopyOverwrite, src)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err = os.ReadFile(dst)
	if string(bytes) != testContent {
		t.Error("File should overwrite")
	}

	// test shrink of file after copying over
	testContent3 := "I am a smaller test"
	src3 := filepath.Join(tempDir, "test3.txt")
	err = os.WriteFile(src3, []byte(testContent3), 0755)

	err = CopyFiles(dst, CopyOverwrite, src3)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err = os.ReadFile(dst)
	if string(bytes) != testContent3 {
		t.Error("File content does not match")
	}

	// Test the copy directory capability of CopyFiles
	dst2 := filepath.Join(os.TempDir(), "gofileTest2")
	defer os.RemoveAll(dst2)

	err = CopyFiles(dst2, CopyDoNotOverwrite, tempDir)
	if err == nil {
		t.Error("Copying a directory to a non-existant location should fail")
	}

	err = CopyFiles(dst2+string(filepath.Separator), CopyDoNotOverwrite, src)
	if err == nil {
		t.Error("Copying a directory to a non-existant location should fail")
	}

	err = os.Mkdir(dst2, 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = CopyFiles(dst2, CopyDoNotOverwrite, tempDir)
	if err != nil {
		t.Fatal(err)
	}

	items, _ := os.ReadDir(dst2)
	if items[0].Name() != "gofileTest" {
		t.Fatal("First item in directory is not gofileTest")
	}

	// Now copy individual items
	var items2 []string
	items, _ = os.ReadDir(tempDir)
	for _, item := range items {
		items2 = append(items2, filepath.Join(tempDir, item.Name()))
	}
	err = CopyFiles(dst2, CopyDoNotOverwrite, items2...)
	if err != nil {
		t.Fatal(err)
	}
	items, _ = os.ReadDir(dst2)
	if len(items) != 5 {
		t.Error("Items were not copied")
	}

	_ = os.RemoveAll(dst2)
	_ = os.Mkdir(dst2, 0777)

	err = CopyFiles(dst2+string(filepath.Separator), CopyOverwrite, src)
	if err != nil {
		t.Fatal(err)
	}
	items, _ = os.ReadDir(dst2)
	if items[0].Name() != "test.txt" {
		t.Fatal("First item in directory is not gofileTest")
	}

	// Check some error states
	err = CopyFiles("", CopyDoNotOverwrite, items2...)
	if err == nil {
		t.Error("Error expected")
	}
	err = CopyFiles(dst2, CopyDoNotOverwrite, "")
	if err == nil {
		t.Error("Error expected")
	}
	err = CopyFiles(dst2, CopyDoNotOverwrite, "random")
	if err == nil {
		t.Error("Error expected")
	}
	err = CopyFiles(src3, CopyDoNotOverwrite, items2...)
	if err == nil {
		t.Error("Error expected")
	}

}

func TestDirectoryCopy(t *testing.T) {
	dir1 := filepath.Join(os.TempDir(), "dir1")
	dir2 := filepath.Join(os.TempDir(), "dir2")

	err := os.Mkdir(dir1, 0777)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.RemoveAll(dir1)

	if !IsDir(dir1) {
		t.Fatal("directory not detected")
	}

	err = os.Mkdir(dir2, 0777)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.RemoveAll(dir2)

	subdir := filepath.Join(dir1, "subdir1")
	if err := os.Mkdir(subdir, 0777); err != nil {
		t.Fatal(err)
	}

	// set up the test directory
	testContent := "I am a test"
	if err := os.WriteFile(filepath.Join(dir1, "test1"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir1, "test2"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "test3"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}

	if err := CopyDirectory(dir1, dir2, CopyOverwrite); err != nil {
		t.Fatal(err)
	}

	items, err := os.ReadDir(dir2)
	if err != nil {
		t.Fatal(err)
	}
	if items[0].Name() != "dir1" {
		t.Fatal("First item in directory is not dir1")
	}
	items, err = os.ReadDir(filepath.Join(dir2, "dir1"))
	if err != nil {
		t.Fatal(err)
	}

	if items[0].Name() != "subdir1" {
		t.Fatal("First item in directory is not subdir1, but rather: " + items[0].Name())
	}
	if items[1].Name() != "test1" {
		t.Fatal("Second item in directory is not test1, but rather: " + items[1].Name())
	}
}

func TestSplitList(t *testing.T) {
	s := "a;b"
	l := SplitList(s)
	if len(l) != 2 {
		t.Error("List is not split correctly")
	}

	s = "a:b:c"
	l = SplitList(s)
	if len(l) != 3 {
		t.Error("List is not split correctly")
	}

}

func TestDirectoryCopyEx(t *testing.T) {
	dir1 := filepath.Join(os.TempDir(), "dir1")
	dir2 := filepath.Join(os.TempDir(), "dir2")
	os.RemoveAll(dir1)
	err := os.Mkdir(dir1, 0777)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.RemoveAll(dir1)

	if !IsDir(dir1) {
		t.Fatal("directory not detected")
	}
	os.RemoveAll(dir2)
	err = os.Mkdir(dir2, 0777)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.RemoveAll(dir2)

	subdir := filepath.Join(dir1, "subdir1")
	if err := os.Mkdir(subdir, 0777); err != nil {
		t.Fatal(err)
	}

	// set up the test directory
	testContent := "I am a test"
	if err := os.WriteFile(filepath.Join(dir1, "test1"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir1, "test2"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir1, "test3.abc"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "test4"), []byte(testContent), 0755); err != nil {
		t.Fatal(err)
	}

	if err := CopyDirectoryEx(dir1, dir2, CopyOverwrite, []string{"*.abc"}); err != nil {
		t.Fatal(err)
	}

	items, err := os.ReadDir(dir2)
	if err != nil {
		t.Fatal(err)
	}
	if items[0].Name() != "dir1" {
		t.Fatal("First item in directory is not dir1")
	}
	items, err = os.ReadDir(filepath.Join(dir2, "dir1"))
	if err != nil {
		t.Fatal(err)
	}

	if items[0].Name() != "subdir1" {
		t.Fatal("First item in directory is not subdir1, but rather: " + items[0].Name())
	}
	if items[1].Name() != "test1" {
		t.Fatal("Second item in directory is not test1, but rather: " + items[1].Name())
	}
	if len(items) != 3 {
		t.Fatalf("Did not copy the correct number of items: %d", len(items))
	}
}

func TestModuleExpandFileList(t *testing.T) {
	modules, err := ModulePaths()
	if err != nil {
		t.Fatal(err)
	}

	fileList := []string{
		"github.com/goradd/gofile/*.md", // readme file
		"github.com/goradd/gofile/LICENSE",
		"/a/b/c", // non-existent items should be left alone
	}
	l := ModuleExpandFileList(fileList, modules)

	if len(l) != 3 {
		t.Error("Not the correct list size.")
	}

	for _, i := range l {
		if i[:6] == "github" {
			t.Error("Item 1 or 2 was not changed.")
		}
	}

	if !listContains(l, filepath.FromSlash("/a/b/c")) {
		t.Error("Item 3 was changed.")
	}
}

func listContains(list []string, val string) bool {
	for _, l := range list {
		if l == val {
			return true
		}
	}
	return false
}

func Test_copyFileTo(t *testing.T) {
	dir, err := makeTempDir()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer os.RemoveAll(dir)

	src := filepath.Join(dir, "test.txt")

	testContent := "I am a test \n and the 2nd line"

	err = os.WriteFile(src, []byte(testContent), 0755)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		src       string
		destDir   string
		name      string
		overwrite CopyOverwriteType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"bad src", args{src: "random", destDir: "random", name: "random", overwrite: CopyOverwrite}, true},
		{"bad dest", args{src: dir, destDir: "random", name: "random", overwrite: CopyOverwrite}, true},
		{"copy once", args{src: src, destDir: dir, name: "test2", overwrite: CopyOverwrite}, false},
		{"copy twice", args{src: src, destDir: dir, name: "test2", overwrite: CopyOverwrite}, false},
		{"copy no overwrite", args{src: src, destDir: dir, name: "test2", overwrite: CopyDoNotOverwrite}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := copyFileTo(tt.args.src, tt.args.destDir, tt.args.name, tt.args.overwrite); (err != nil) != tt.wantErr {
				t.Errorf("copyFileTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCopyFilesExErrors(t *testing.T) {
	type args struct {
		dst        string
		overwrite  CopyOverwriteType
		exclusions []string
		src        []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no source", args{testDataDir1(), CopyDoNotOverwrite, []string{}, []string{}}, true},
		{"no dest", args{"", CopyDoNotOverwrite, []string{}, []string{testDataDir1()}}, true},
		{"bad src", args{testDataDir1(), CopyDoNotOverwrite, []string{}, []string{"random"}}, true},
		{"bad dest", args{"random", CopyDoNotOverwrite, []string{}, []string{filepath.Join("testdata", "t1.txt"), filepath.Join("testdata", "t2.txt")}}, true},
		{"file dest", args{filepath.Join("testdata", "t1.txt"), CopyDoNotOverwrite, []string{}, []string{testDataDir1()}}, true},
		{"bad dir dest", args{"random1/", CopyDoNotOverwrite, []string{}, []string{filepath.Join("testdata", "t1.txt")}}, true},
		{"bad parent dir", args{"random1/bad2", CopyDoNotOverwrite, []string{}, []string{filepath.Join("testdata", "t1.txt")}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyFilesEx(tt.args.dst, tt.args.overwrite, tt.args.exclusions, tt.args.src...); (err != nil) != tt.wantErr {
				t.Errorf("CopyFilesEx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCopyFilesEx(t *testing.T) {

	t.Run("copy dir", func(t *testing.T) {
		tempDir, err := makeTempDir()
		if err != nil {
			t.Fatal(err)
			return
		}
		defer os.RemoveAll(tempDir)

		if err = CopyFilesEx(tempDir, CopyOverwrite, []string{"*.abc"}, filepath.Join(testDataDir1(), "a")); err != nil {
			t.Errorf("CopyFilesEx() error = %v", err)
		}
		items, _ := os.ReadDir(tempDir)
		if len(items) != 1 {
			t.Errorf("CopyFilesEx() did not copy correct number of directories. Wanted 1, found %d", len(items))
		}

		items, _ = os.ReadDir(filepath.Join(tempDir, "a"))
		if len(items) != 2 {
			t.Errorf("CopyFilesEx() did not copy correct number of files. Wanted 2, found %d", len(items))
		}

	})

	t.Run("copy one", func(t *testing.T) {
		tempDir, err := makeTempDir()
		if err != nil {
			t.Fatal(err)
			return
		}
		defer os.RemoveAll(tempDir)

		if err = CopyFilesEx(tempDir+string(filepath.Separator), CopyOverwrite, []string{"*.abc"}, filepath.Join("testdata", "t1.txt")); err != nil {
			t.Errorf("CopyFilesEx() error = %v", err)
		}
		items, _ := os.ReadDir(tempDir)
		if len(items) != 1 {
			t.Errorf("CopyFilesEx() did not copy correct number of items. Wanted 1, found %d", len(items))
		}
	})

	t.Run("copy sub", func(t *testing.T) {
		tempDir, err := makeTempDir()
		if err != nil {
			t.Fatal(err)
			return
		}
		defer os.RemoveAll(tempDir)

		if err = CopyFilesEx(filepath.Join(tempDir, "test1"), CopyOverwrite, []string{"*.abc"}, filepath.Join("testdata", "t1.txt")); err != nil {
			t.Errorf("CopyFilesEx() error = %v", err)
		}
		items, _ := os.ReadDir(tempDir)
		if len(items) != 1 {
			t.Errorf("CopyFilesEx() did not copy correct number of items. Wanted 1, found %d", len(items))
		}
	})

	t.Run("copy sub excluded", func(t *testing.T) {
		tempDir, err := makeTempDir()
		if err != nil {
			t.Fatal(err)
			return
		}
		defer os.RemoveAll(tempDir)

		if err = CopyFilesEx(filepath.Join(tempDir, "test1"), CopyOverwrite, []string{"*.abc"}, filepath.Join("testdata", "t3.abc")); err != nil {
			t.Errorf("CopyFilesEx() error = %v", err)
		}
		items, _ := os.ReadDir(tempDir)
		if len(items) != 0 {
			t.Errorf("CopyFilesEx() did not copy correct number of items. Wanted 0, found %d", len(items))
		}
	})
	t.Run("copy sub excluded", func(t *testing.T) {
		tempDir, err := makeTempDir()
		if err != nil {
			t.Fatal(err)
			return
		}
		defer os.RemoveAll(tempDir)

		if err = CopyFilesEx(filepath.Join(tempDir, "test1"), CopyOverwrite, []string{"*.abc"}, filepath.Join("testdata", "t3.abc")); err != nil {
			t.Errorf("CopyFilesEx() error = %v", err)
		}
		items, _ := os.ReadDir(tempDir)
		if len(items) != 0 {
			t.Errorf("CopyFilesEx() did not copy correct number of items. Wanted 0, found %d", len(items))
		}
	})

	t.Run("copy file exists", func(t *testing.T) {
		tempDir, err := makeTempDir()
		if err != nil {
			t.Fatal(err)
			return
		}
		defer os.RemoveAll(tempDir)

		if err = CopyFilesEx(filepath.Join(tempDir, "test1"), CopyOverwrite, []string{"*.abc"}, filepath.Join("testdata", "t1.txt")); err != nil {
			t.Errorf("CopyFilesEx() error = %v", err)
		}
		items, _ := os.ReadDir(tempDir)
		if len(items) != 1 {
			t.Errorf("CopyFilesEx() did not copy correct number of items. Wanted 1, found %d", len(items))
		}

		if err = CopyFilesEx(filepath.Join(tempDir, "test1"), CopyOverwrite, []string{"*.abc"}, filepath.Join("testdata", "t1.txt")); err != nil {
			t.Errorf("CopyFilesEx() error = %v", err)
		}
		items, _ = os.ReadDir(tempDir)
		if len(items) != 1 {
			t.Errorf("CopyFilesEx() did not copy correct number of items. Wanted 1, found %d", len(items))
		}

	})

	t.Run("copy one to dir existing", func(t *testing.T) {
		tempDir, err := makeTempDir()
		if err != nil {
			t.Fatal(err)
			return
		}
		defer os.RemoveAll(tempDir)

		if err = CopyFilesEx(tempDir, CopyOverwrite, []string{"*.abc"}, filepath.Join("testdata", "t1.txt")); err != nil {
			t.Errorf("CopyFilesEx() error = %v", err)
		}
		items, _ := os.ReadDir(tempDir)
		if len(items) != 1 {
			t.Errorf("CopyFilesEx() did not copy correct number of directories. Wanted 1, found %d", len(items))
		}

		if err = CopyFilesEx(tempDir, CopyOverwrite, []string{"*.abc"}, filepath.Join("testdata", "t1.txt")); err != nil {
			t.Errorf("CopyFilesEx() error = %v", err)
		}
		items, _ = os.ReadDir(tempDir)
		if len(items) != 1 {
			t.Errorf("CopyFilesEx() did not copy correct number of directories. Wanted 1, found %d", len(items))
		}

	})

	t.Run("copy one to dir existing excluded", func(t *testing.T) {
		tempDir, err := makeTempDir()
		if err != nil {
			t.Fatal(err)
			return
		}
		defer os.RemoveAll(tempDir)

		if err = CopyFilesEx(tempDir, CopyOverwrite, []string{}, filepath.Join("testdata", "t1.txt")); err != nil {
			t.Errorf("CopyFilesEx() error = %v", err)
		}
		items, _ := os.ReadDir(tempDir)
		if len(items) != 1 {
			t.Errorf("CopyFilesEx() did not copy correct number of directories. Wanted 1, found %d", len(items))
		}

		if err = CopyFilesEx(tempDir, CopyOverwrite, []string{"*.txt"}, filepath.Join("testdata", "t2.txt")); err != nil {
			t.Errorf("CopyFilesEx() error = %v", err)
		}
		items, _ = os.ReadDir(tempDir)
		if len(items) != 1 {
			t.Errorf("CopyFilesEx() did not copy correct number of directories. Wanted 1, found %d", len(items))
		}

	})

}
