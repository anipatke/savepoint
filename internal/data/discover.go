package data

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ReleaseInfo struct {
	ID    string
	Path  string
	Epics []EpicInfo
}

type EpicInfo struct {
	ID    string
	Path  string
	Tasks []TaskInfo
}

type TaskInfo struct {
	ID   string
	Path string
}

// Discover is the legacy V1 directory walker. It is retained for explicit
// migration and frozen historical fixtures only; live commands use the V2
// identity-keyed discovery functions below.
type Discover struct{}

func NewDiscover() *Discover {
	return &Discover{}
}

func (d *Discover) FindSavepointRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		savepointPath := filepath.Join(dir, ".savepoint")
		info, err := os.Stat(savepointPath)
		if err == nil && info.IsDir() {
			return savepointPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrSavepointDirectoryMissing
		}
		dir = parent
	}
}

func (d *Discover) ListReleases(root string) ([]ReleaseInfo, error) {
	releasesPath := filepath.Join(root, "releases")
	info, err := os.Stat(releasesPath)
	if err != nil {
		return nil, fmt.Errorf("releases directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("releases is not a directory")
	}

	entries, err := os.ReadDir(releasesPath)
	if err != nil {
		return nil, err
	}

	var releases []ReleaseInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id := entry.Name()
		releases = append(releases, ReleaseInfo{
			ID:   id,
			Path: filepath.Join(releasesPath, id),
		})
	}

	sort.Slice(releases, func(i, j int) bool {
		return releases[i].ID < releases[j].ID
	})
	return releases, nil
}

// ListRootDirs returns sorted child directory names directly under root.
func (d *Discover) ListRootDirs(root string) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", root)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	sort.Strings(dirs)
	return dirs, nil
}

func (d *Discover) ListEpics(root, release string) ([]EpicInfo, error) {
	epicsPath := filepath.Join(root, "releases", release, "epics")
	info, err := os.Stat(epicsPath)
	if err != nil {
		return nil, fmt.Errorf("epics directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("epics is not a directory")
	}

	entries, err := os.ReadDir(epicsPath)
	if err != nil {
		return nil, err
	}

	var epics []EpicInfo
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		id := entry.Name()
		epics = append(epics, EpicInfo{
			ID:   id,
			Path: filepath.Join(epicsPath, id),
		})
	}

	sort.Slice(epics, func(i, j int) bool {
		return epics[i].ID < epics[j].ID
	})
	return epics, nil
}

func (d *Discover) ListDefects(root, release string) ([]DefectInfo, error) {
	defectsPath := filepath.Join(root, "releases", release, "defects")
	info, err := os.Stat(defectsPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("stat defects dir: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("defects is not a directory")
	}

	entries, err := os.ReadDir(defectsPath)
	if err != nil {
		return nil, fmt.Errorf("read defects dir: %w", err)
	}

	var defects []DefectInfo
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		id := entry.Name()[:len(entry.Name())-3]
		defects = append(defects, DefectInfo{
			ID:   id,
			Path: filepath.Join(defectsPath, entry.Name()),
		})
	}

	sort.Slice(defects, func(i, j int) bool {
		return defects[i].ID < defects[j].ID
	})
	return defects, nil
}

// V2 record discovery layout, fixed by the V2 file model:
//
//	objectives/O###-slug/Objective.md
//	objectives/O###-slug/tasks/T###-slug.md
const (
	v2ObjectivesDirName = "objectives"
	v2ObjectiveFileName = "Objective.md"
	v2TasksDirName      = "tasks"
	v2ReleasesDirName   = "releases"
	v2ReleaseFileName   = "Release.md"
	v2ChecksDirName     = "checks"
	v2IssuesDirName     = "issues"
)

// DiscoverV2Releases confines discovery to root/releases and decodes the
// optional Release record family. A release directory is identified by its
// own Release.md record; the directory slug never supplies or changes the
// declared R### identity. Missing releases/ and directories without a
// Release.md are valid, which keeps release-free V2 projects unchanged.
func DiscoverV2Releases(root string) (map[string]*ReleaseV2, error) {
	releases := map[string]*ReleaseV2{}
	releasesPath := filepath.Join(root, v2ReleasesDirName)
	if _, statErr := os.Lstat(releasesPath); os.IsNotExist(statErr) {
		return releases, nil
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	confined := &v2PathConfiner{rootAbs: rootAbs, seen: map[string]string{}}

	releasesInfo, err := confined.stat(releasesPath, v2ReleasesDirName)
	if err != nil {
		return nil, err
	}
	if !releasesInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", releasesPath)
	}

	entries, err := os.ReadDir(releasesPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		relDir := filepath.Join(v2ReleasesDirName, entry.Name())
		dirPath := filepath.Join(releasesPath, entry.Name())
		dirInfo, err := confined.stat(dirPath, relDir)
		if err != nil {
			return nil, err
		}
		if !dirInfo.IsDir() {
			continue
		}

		relFile := filepath.Join(relDir, v2ReleaseFileName)
		filePath := filepath.Join(dirPath, v2ReleaseFileName)
		fileInfo, statErr := confined.stat(filePath, relFile)
		if statErr != nil {
			if os.IsNotExist(statErr) {
				continue
			}
			return nil, statErr
		}
		if fileInfo.IsDir() {
			return nil, fmt.Errorf("%s is a directory, want a file", filePath)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", filePath, err)
		}
		release, err := DecodeReleaseV2(relFile, string(content))
		if err != nil {
			return nil, err
		}
		release.Source.ProjectRoot = rootAbs

		if !releasePathMatchesID(entry.Name(), release.ID) {
			return nil, fmt.Errorf("%w: %s: release %s directory name %q does not match its id", ErrV2PathMismatch, relDir, release.ID, entry.Name())
		}

		id := string(release.ID)
		if existing, ok := releases[id]; ok {
			return nil, fmt.Errorf("%w: release %s declared at both %s and %s", ErrV2DuplicateID, id, existing.Source.Path, relFile)
		}
		releases[id] = release
	}

	return releases, nil
}

func releasePathMatchesID(name string, id ReleaseID) bool {
	return name == string(id) || strings.HasPrefix(name, string(id)+"-")
}

// DiscoverV2Records confines discovery to root/objectives and decodes every
// Objective and Task record found there into identity-keyed maps. It never
// consults record content to build filesystem paths: every path walked comes
// from a real directory entry, and every path resolved through symlinks or
// aliasing case collisions is rejected rather than followed. It also rejects
// duplicate global IDs and a directory/file name that does not start with
// the ID the record declares. root is the .savepoint directory.
//
// DiscoverV2Records validates structure only: it does not infer Check
// clearance, owner acceptance, or dependency satisfaction, and it does not
// validate ownership or reference graphs — that is dependency.go's job once
// the caller has both maps.
func DiscoverV2Records(root string) (objectives map[string]*ObjectiveV2, tasks map[string]*TaskV2, err error) {
	objectives = map[string]*ObjectiveV2{}
	tasks = map[string]*TaskV2{}

	objectivesPath := filepath.Join(root, v2ObjectivesDirName)
	if _, statErr := os.Lstat(objectivesPath); os.IsNotExist(statErr) {
		return objectives, tasks, nil
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, nil, err
	}

	confined := &v2PathConfiner{rootAbs: rootAbs, seen: map[string]string{}}

	objectivesInfo, err := confined.stat(objectivesPath, v2ObjectivesDirName)
	if err != nil {
		return nil, nil, err
	}
	if !objectivesInfo.IsDir() {
		return nil, nil, fmt.Errorf("%s is not a directory", objectivesPath)
	}

	objEntries, err := os.ReadDir(objectivesPath)
	if err != nil {
		return nil, nil, err
	}

	for _, objEntry := range objEntries {
		objDirPath := filepath.Join(objectivesPath, objEntry.Name())
		relObjDir := filepath.Join(v2ObjectivesDirName, objEntry.Name())

		objDirInfo, err := confined.stat(objDirPath, relObjDir)
		if err != nil {
			return nil, nil, err
		}
		if !objDirInfo.IsDir() {
			// A stray file directly under objectives/; not an Objective.
			continue
		}

		objFilePath := filepath.Join(objDirPath, v2ObjectiveFileName)
		relObjFile := filepath.Join(relObjDir, v2ObjectiveFileName)
		objFileInfo, statErr := confined.stat(objFilePath, relObjFile)
		if statErr != nil {
			if os.IsNotExist(statErr) {
				// Tasks are independently owned by their explicit objective
				// field. Discover them even when the containing Objective record
				// is absent so LoadV2Index can report the missing owner instead
				// of silently dropping the work.
				if err := discoverV2Tasks(confined, objDirPath, relObjDir, tasks); err != nil {
					return nil, nil, err
				}
				continue
			}
			return nil, nil, statErr
		}
		if objFileInfo.IsDir() {
			return nil, nil, fmt.Errorf("%s is a directory, want a file", objFilePath)
		}

		content, err := os.ReadFile(objFilePath)
		if err != nil {
			return nil, nil, fmt.Errorf("read %s: %w", objFilePath, err)
		}

		objective, err := DecodeObjectiveV2(relObjFile, string(content))
		if err != nil {
			return nil, nil, err
		}
		objective.Source.ProjectRoot = rootAbs

		if !strings.HasPrefix(objEntry.Name(), objective.ID) {
			return nil, nil, fmt.Errorf("%w: %s: objective %s directory name %q does not start with its id", ErrV2PathMismatch, relObjDir, objective.ID, objEntry.Name())
		}

		if existing, ok := objectives[objective.ID]; ok {
			return nil, nil, fmt.Errorf("%w: objective %s declared at both %s and %s", ErrV2DuplicateID, objective.ID, existing.Source.Path, relObjFile)
		}
		objectives[objective.ID] = objective

		if err := discoverV2Tasks(confined, objDirPath, relObjDir, tasks); err != nil {
			return nil, nil, err
		}
	}

	return objectives, tasks, nil
}

// discoverV2Tasks walks objDirPath/tasks, decoding every *.md file into
// tasks. A missing tasks directory is not an error: an Objective may not yet
// own any Task.
func discoverV2Tasks(confined *v2PathConfiner, objDirPath, relObjDir string, tasks map[string]*TaskV2) error {
	tasksPath := filepath.Join(objDirPath, v2TasksDirName)
	relTasksDir := filepath.Join(relObjDir, v2TasksDirName)

	tasksInfo, err := confined.stat(tasksPath, relTasksDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !tasksInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", tasksPath)
	}

	entries, err := os.ReadDir(tasksPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		taskFilePath := filepath.Join(tasksPath, entry.Name())
		relTaskFile := filepath.Join(relTasksDir, entry.Name())
		taskFileInfo, err := confined.stat(taskFilePath, relTaskFile)
		if err != nil {
			return err
		}
		if taskFileInfo.IsDir() {
			continue
		}

		content, err := os.ReadFile(taskFilePath)
		if err != nil {
			return fmt.Errorf("read %s: %w", taskFilePath, err)
		}

		task, err := DecodeTaskV2(relTaskFile, string(content))
		if err != nil {
			return err
		}
		task.Source.ProjectRoot = confined.rootAbs

		baseName := strings.TrimSuffix(entry.Name(), ".md")
		if !strings.HasPrefix(baseName, task.ID) {
			return fmt.Errorf("%w: %s: task %s file name %q does not start with its id", ErrV2PathMismatch, relTaskFile, task.ID, entry.Name())
		}

		if existing, ok := tasks[task.ID]; ok {
			return fmt.Errorf("%w: task %s declared at both %s and %s", ErrV2DuplicateID, task.ID, existing.Source.Path, relTaskFile)
		}
		tasks[task.ID] = task
	}

	return nil
}

// DiscoverV2Checks confines discovery to root/checks and decodes every Check
// record found there into an identity-keyed map. root is the .savepoint
// directory. A missing checks/ directory is not an error: a project may not
// yet have any recorded evaluations.
//
// DiscoverV2Checks validates structure only: it does not resolve a Check's
// scope target against the rest of the index, or validate supersedes chains
// — that is project.go's job once the caller has the full V2Index.
func DiscoverV2Checks(root string) (map[string]*CheckV2, error) {
	return discoverV2FlatDir(root, v2ChecksDirName, "check", DecodeCheckV2)
}

// DiscoverV2Issues confines discovery to root/issues and decodes every Issue
// record found there into an identity-keyed map, under the same confinement,
// filename-prefix, and duplicate-identity rules Checks use. root is the
// .savepoint directory. A missing issues/ directory is not an error: a
// project may have no durable follow-up yet.
//
// DiscoverV2Issues validates structure only: it does not resolve an Issue's
// task, check, or duplicate_of references against the rest of the index —
// that is project.go's job once the caller has the full V2Index.
func DiscoverV2Issues(root string) (map[string]*IssueV2, error) {
	return discoverV2FlatDir(root, v2IssuesDirName, "issue", DecodeIssueV2)
}

// v2FlatRecord is what a record family filed directly under one flat
// directory must expose for discoverV2FlatDir to index it: the global ID the
// record declares, and a handle to the source document whose resolution
// context discovery fills in.
type v2FlatRecord interface {
	recordID() string
	sourceDocument() *V2SourceDocument
}

func (c *CheckV2) recordID() string                  { return c.ID }
func (c *CheckV2) sourceDocument() *V2SourceDocument { return &c.Source }

func (i *IssueV2) recordID() string                  { return i.ID }
func (i *IssueV2) sourceDocument() *V2SourceDocument { return &i.Source }

// discoverV2FlatDir walks root/dirName and decodes every *.md file there
// through decode, keying the result by declared ID. Like DiscoverV2Records it
// never consults record content to build filesystem paths, rejects a path
// resolving outside the project root or aliasing an already-visited path on a
// case-insensitive filesystem, and rejects a duplicate global ID and a
// filename that does not start with the ID the record declares. recordKind
// names the family in diagnostics. A missing directory is not an error.
func discoverV2FlatDir[T v2FlatRecord](root, dirName, recordKind string, decode func(path, content string) (T, error)) (map[string]T, error) {
	records := map[string]T{}

	dirPath := filepath.Join(root, dirName)
	if _, statErr := os.Lstat(dirPath); os.IsNotExist(statErr) {
		return records, nil
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	confined := &v2PathConfiner{rootAbs: rootAbs, seen: map[string]string{}}

	dirInfo, err := confined.stat(dirPath, dirName)
	if err != nil {
		return nil, err
	}
	if !dirInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dirPath)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		relFile := filepath.Join(dirName, entry.Name())
		fileInfo, err := confined.stat(filePath, relFile)
		if err != nil {
			return nil, err
		}
		if fileInfo.IsDir() {
			continue
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", filePath, err)
		}

		record, err := decode(relFile, string(content))
		if err != nil {
			return nil, err
		}
		record.sourceDocument().ProjectRoot = rootAbs

		id := record.recordID()
		baseName := strings.TrimSuffix(entry.Name(), ".md")
		if !strings.HasPrefix(baseName, id) {
			return nil, fmt.Errorf("%w: %s: %s %s file name %q does not start with its id", ErrV2PathMismatch, relFile, recordKind, id, entry.Name())
		}

		if existing, ok := records[id]; ok {
			return nil, fmt.Errorf("%w: %s %s declared at both %s and %s", ErrV2DuplicateID, recordKind, id, existing.sourceDocument().Path, relFile)
		}
		records[id] = record
	}

	return records, nil
}

// v2PathConfiner rejects two ways a discovered path could escape the safe,
// project-relative interpretation V2 discovery requires: resolving (through
// one or more symlinks, at any path component) outside the project root, and
// aliasing an already-visited project-relative path on a case-insensitive
// filesystem. stat is the only way discovery inspects a path, so neither
// check can be bypassed by an entry that is itself a symlink.
type v2PathConfiner struct {
	rootAbs string
	seen    map[string]string // lowercased relative path -> original relative path
}

// stat confines absPath (identified for diagnostics by project-relative
// relPath) and, once confined, returns the FileInfo of what it resolves to.
// A not-exist error from the initial Lstat is returned unwrapped so callers
// can keep using os.IsNotExist.
func (c *v2PathConfiner) stat(absPath, relPath string) (os.FileInfo, error) {
	if _, err := os.Lstat(absPath); err != nil {
		return nil, err
	}

	real, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return nil, fmt.Errorf("resolve %s: %w", absPath, err)
	}
	realAbs, err := filepath.Abs(real)
	if err != nil {
		return nil, err
	}
	if realAbs != c.rootAbs && !strings.HasPrefix(realAbs, c.rootAbs+string(filepath.Separator)) {
		return nil, fmt.Errorf("%w: %s: resolves outside the project root", ErrV2UnsafePath, relPath)
	}

	key := strings.ToLower(relPath)
	if existing, ok := c.seen[key]; ok && existing != relPath {
		return nil, fmt.Errorf("%w: %s and %s alias on a case-insensitive filesystem", ErrV2UnsafePath, existing, relPath)
	}
	c.seen[key] = relPath

	return os.Stat(absPath)
}

func (d *Discover) ListTasks(root, release, epic string) ([]TaskInfo, error) {
	tasksPath := filepath.Join(root, "releases", release, "epics", epic, "tasks")
	info, err := os.Stat(tasksPath)
	if err != nil {
		return nil, fmt.Errorf("tasks directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("tasks is not a directory")
	}

	entries, err := os.ReadDir(tasksPath)
	if err != nil {
		return nil, err
	}

	var tasks []TaskInfo
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		id := entry.Name()[:len(entry.Name())-3]
		tasks = append(tasks, TaskInfo{
			ID:   id,
			Path: filepath.Join(tasksPath, entry.Name()),
		})
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})
	return tasks, nil
}
