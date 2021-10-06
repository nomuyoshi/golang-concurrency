package main

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
)

type result struct {
	path string
	sum  [md5.Size]byte
	err  error // ファイル読み取り時のエラー
}

func main() {
	ctx := context.Background()
	defer func() {
		fmt.Println("goroutine数 = ", runtime.NumGoroutine())
	}()
	m, err := MD5All(ctx, os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x  %s\n", m[path], path)
	}
}

func walkFiles(ctx context.Context, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(paths)
		// errc はバッファありなので送信操作はブロックされない
		errc <- filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-ctx.Done():
				fmt.Println("walk canceled")
				return errors.New("walk canceled")
			}
			return nil
		})
	}()

	return paths, errc
}

func digester(ctx context.Context, paths <-chan string, c chan<- result) {
	for path := range paths {
		data, err := ioutil.ReadFile(path)
		select {
		case c <- result{path: path, sum: md5.Sum(data), err: err}:
		case <-ctx.Done():
			fmt.Println("digester canceled")
			return
		}
	}
}

func MD5All(pctx context.Context, root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	ctx, cancel := context.WithCancel(pctx)
	c := make(chan result)
	defer cancel()

	paths, errc := walkFiles(ctx, root)

	const concurrency = 10
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			digester(ctx, paths, c)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for r := range c {
		if r.err != nil {
			return nil, r.err
		}

		m[r.path] = r.sum
	}

	if err := <-errc; err != nil {
		return nil, err
	}

	return m, nil
}
