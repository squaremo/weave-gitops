package git_test

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/weave-gitops/pkg/git"
)

var (
	gitClient git.Git
	dir       string
	err       error
)

var _ = BeforeEach(func() {
	gitClient = git.New(nil)

	dir, err = ioutil.TempDir("", "wego-git-test-")
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterEach(func() {
	Expect(os.RemoveAll(dir)).To(Succeed())
})

var _ = Describe("Init", func() {
	It("creates an empty repository", func() {
		_, err := gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())

		_, err = gitClient.Open(dir)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("does nothing when the repository has already been initialized", func() {
		init, err := gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(init).To(BeTrue())

		init, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(init).To(BeFalse())
	})

	It("returns an error when the directory already contains a repository", func() {
		_, err := gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())

		init, err := git.New(nil).Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).Should(MatchError("repository already exists"))
		Expect(init).To(BeFalse())
	})
})

var _ = Describe("Open", func() {
	It("fails when the directory is an empty directory", func() {
		_, err := gitClient.Open(dir)
		Expect(err).To(MatchError("repository does not exist"))
	})
})

var _ = Describe("Clone", func() {
	It("clones a given repository", func() {
		_, err := gitClient.Clone(context.Background(), dir, "https://github.com/githubtraining/hellogitworld", "master")
		Expect(err).ShouldNot(HaveOccurred())

		_, err = os.Stat(dir + "/README.txt")
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("initialize a given repository if remote branch not found", func() {
		_, err := gitClient.Clone(context.Background(), dir, "https://github.com/githubtraining/hellogitworld", "new-branch")
		Expect(err).ShouldNot(HaveOccurred())

		_, err = gitClient.Open(dir)
		Expect(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("Write", func() {
	It("writes a file into a given repository", func() {
		_, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())
		filePath := "/test.txt"
		content := []byte("testing")
		Expect(gitClient.Write(filePath, content)).To(Succeed())

		fileContent, err := ioutil.ReadFile(dir + filePath)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(content).To(Equal(fileContent))
	})

	It("fails when the repository has not been initialized", func() {
		_, err := gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())

		gc := git.New(nil)
		err = gc.Write("testing.txt", []byte("testing"))
		Expect(err).To(MatchError("no git repository"))
	})

	It("returns an error if the repository is bare", func() {
		_, err = gogit.PlainInit(dir, true)
		Expect(err).ShouldNot(HaveOccurred())
		_, err = gitClient.Open(dir)
		Expect(err).ShouldNot(HaveOccurred())

		err := gitClient.Write("testing.txt", []byte("testing"))
		Expect(err).Should(MatchError("failed to open the worktree: worktree not available in a bare repository"))
	})
})

var _ = Describe("Commit", func() {
	It("commits into a given repository", func() {
		_, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())

		filePath := "/test.txt"
		content := []byte("testing")
		Expect(gitClient.Write(filePath, content)).To(Succeed())

		_, err = gitClient.Commit(git.Commit{
			Author:  git.Author{Name: "test", Email: "test@example.com"},
			Message: "test commit",
		})
		Expect(err).ShouldNot(HaveOccurred())

		out := executeCommand(dir, "sh", "-c", `git log -1 --pretty=%B`)
		Expect(out).To(ContainSubstring("test commit"))
	})

	It("commits into a given repository skipping filtered files", func() {
		_, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())

		filePath := "/test.txt"
		content := []byte("testing")
		Expect(gitClient.Write(filePath, content)).To(Succeed())

		_, err = gitClient.Commit(git.Commit{
			Author:  git.Author{Name: "test", Email: "test@example.com"},
			Message: "test commit",
		},
			func(fname string) bool {
				return false
			})
		Expect(err).Should(HaveOccurred())
		isClean, err := gitClient.Status()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(isClean).To(BeFalse())
	})

	It("commits into a given repository skipping filtered files on .wego folder", func() {
		_, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())

		filePath := ".wego/test.txt"
		content := []byte("testing")
		Expect(gitClient.Write(filePath, content)).To(Succeed())

		_, err = gitClient.Commit(git.Commit{
			Author:  git.Author{Name: "test", Email: "test@example.com"},
			Message: "test commit",
		},
			func(fname string) bool {
				return strings.HasPrefix(fname, ".wego")
			})
		Expect(err).ShouldNot(HaveOccurred())
		isClean, err := gitClient.Status()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(isClean).To(BeTrue())
	})

	It("fails if there are no changes to commit", func() {
		_, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())
		filePath := "/test.txt"
		content := []byte("testing")
		Expect(gitClient.Write(filePath, content)).To(Succeed())

		_, err = gitClient.Commit(git.Commit{
			Author:  git.Author{Name: "test", Email: "test@example.com"},
			Message: "test commit",
		})
		Expect(err).ShouldNot(HaveOccurred())

		out := executeCommand(dir, "sh", "-c", `git log -1 --pretty=%B`)
		Expect(out).To(ContainSubstring("test commit"))

		_, err = gitClient.Commit(git.Commit{
			Author:  git.Author{Name: "test", Email: "test@example.com"},
			Message: "test commit",
		})
		Expect(err).Should(MatchError("no staged files"))
	})
})

var _ = Describe("Status", func() {
	It("returns true if no files have been changed in the repository", func() {
		_, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())

		isClean, err := gitClient.Status()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(isClean).To(BeTrue())
	})

	It("returns false if a file has been changed in the repository", func() {
		_, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())

		filePath := "/test.txt"
		content := []byte("testing")
		Expect(gitClient.Write(filePath, content)).To(Succeed())

		isClean, err := gitClient.Status()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(isClean).To(BeFalse())
	})

	It("fails if not initialized", func() {
		_, err := gitClient.Status()
		Expect(err).Should(MatchError("no git repository"))
	})

	It("returns an error if the repository is bare", func() {
		_, err = gogit.PlainInit(dir, true)
		Expect(err).ShouldNot(HaveOccurred())
		_, err = gitClient.Open(dir)
		Expect(err).ShouldNot(HaveOccurred())
		_, err := gitClient.Status()
		Expect(err).Should(MatchError("failed to open the worktree: worktree not available in a bare repository"))
	})
})

var _ = Describe("Head", func() {
	It("returns if the working tree is clean", func() {
		_, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())

		filePath := "/test.txt"
		content := []byte("testing")
		Expect(gitClient.Write(filePath, content)).To(Succeed())

		_, err = gitClient.Commit(git.Commit{
			Author:  git.Author{Name: "test", Email: "test@example.com"},
			Message: "test commit",
		})
		Expect(err).ShouldNot(HaveOccurred())

		hash, err := gitClient.Head()
		Expect(err).ShouldNot(HaveOccurred())

		out := executeCommand(dir, "sh", "-c", `git log -1`)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).To(ContainSubstring(hash))
	})

	It("fails if not initialized", func() {
		_, err := gitClient.Head()
		Expect(err).Should(MatchError("no git repository"))
	})

	It("fails if no commits in the git repository", func() {
		_, err := gogit.PlainInit(dir, true)
		Expect(err).ShouldNot(HaveOccurred())
		_, err = gitClient.Open(dir)
		Expect(err).ShouldNot(HaveOccurred())

		_, err = gitClient.Head()
		Expect(err).Should(MatchError("reference not found"))
	})
})

var _ = Describe("Push", func() {
	It("fails if not access", func() {
		_, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())

		filePath := "/test.txt"
		content := []byte("testing")
		Expect(gitClient.Write(filePath, content)).To(Succeed())

		_, err = gitClient.Commit(git.Commit{
			Author:  git.Author{Name: "test", Email: "test@example.com"},
			Message: "test commit",
		})
		Expect(err).ShouldNot(HaveOccurred())

		err := gitClient.Push(context.Background())
		Expect(err).Should(HaveOccurred())
	})

	It("fails if not initialized", func() {
		err := gitClient.Push(context.Background())
		Expect(err).Should(MatchError("no git repository"))
	})

	It("fails if no commits in the git repository", func() {
		_, err := gogit.PlainInit(dir, true)
		Expect(err).ShouldNot(HaveOccurred())
		_, err = gitClient.Open(dir)
		Expect(err).ShouldNot(HaveOccurred())

		err = gitClient.Push(context.Background())
		Expect(err).Should(MatchError("remote not found"))
	})
})

var _ = Describe("Remove", func() {
	It("fails if no file present at path in the git repository", func() {
		_, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(gitClient.Remove("foo")).ShouldNot(Succeed())
	})

	It("succeeds if file present at path in the git repository", func() {
		_, err = gitClient.Init(dir, "https://github.com/github/gitignore", "master")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(gitClient.Write("foo", []byte("bar"))).To(Succeed())
		Expect(gitClient.Remove("foo")).To(Succeed())
	})
})

func executeCommand(workingDir, cmd string, args ...string) []byte {
	c := exec.Command(cmd, args...)
	c.Dir = workingDir
	out, err := c.Output()
	Expect(err).ShouldNot(HaveOccurred())
	return out
}
