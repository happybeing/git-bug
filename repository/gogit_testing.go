package repository

import (
	"io/ioutil"
	"log"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
)

// This is intended for testing only

func CreateGoGitTestRepo(bare bool) TestedRepo {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	var creator func(string, billy.Filesystem) (*GoGitRepo, error)

	if bare {
		creator = InitBareGoGitRepo
	} else {
		creator = InitGoGitRepo
	}

	fs := memfs.New()
	repo, err := creator(dir, fs)
	if err != nil {
		log.Fatal(err)
	}

	config := repo.LocalConfig()
	if err := config.StoreString("user.name", "testuser"); err != nil {
		log.Fatal("failed to set user.name for test repository: ", err)
	}
	if err := config.StoreString("user.email", "testuser@example.com"); err != nil {
		log.Fatal("failed to set user.email for test repository: ", err)
	}

	return repo
}

func SetupGoGitReposAndRemote() (repoA, repoB, remote TestedRepo) {
	repoA = CreateGoGitTestRepo(false)
	repoB = CreateGoGitTestRepo(false)
	remote = CreateGoGitTestRepo(true)

	remoteAddr := "file://" + remote.GetPath()

	err := repoA.AddRemote("origin", remoteAddr)
	if err != nil {
		log.Fatal(err)
	}

	err = repoB.AddRemote("origin", remoteAddr)
	if err != nil {
		log.Fatal(err)
	}

	return repoA, repoB, remote
}
