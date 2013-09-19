package trie

import (
	"github.com/fvbock/uds-go/set"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(1)
}

func TestAddSingle(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	if !tr.Root.End {
		t.Error("Expected Root End to be true")
	}
}

func TestAddBigSmall(t *testing.T) {
	tr := NewTrie()
	tr.Add("testing")
	tr.Add("tests")
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
}

func TestAddSmallBig(t *testing.T) {
	tr := NewTrie()
	tr.Add("tests")
	tr.Add("testing")
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
}

func TestAddTestFirst(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	tr.Add("testing")
	tr.Add("tests")
	if !tr.Root.End {
		t.Error("Expected Root End to be true")
	}
	if !tr.Root.End {
		t.Error("Expected trunk End to be true")
	}
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
}

func TestAddTestLast(t *testing.T) {
	tr := NewTrie()
	tr.Add("testing")
	tr.Add("tests")
	tr.Add("test")
	if !tr.Root.End {
		t.Error("Expected Root End to be true")
	}
	if !tr.Root.Branches['i'].End {
		t.Error("Expected 'i' End to be true")
	}
	if !tr.Root.Branches['s'].End {
		t.Error("Expected 's' End to be true")
	}
}

func TestDump(t *testing.T) {
	tr := NewTrie()
	tr.Add("teased")
	tr.Add("test")
	tr.Add("testing")
	t.Logf("\n%s", tr.Dump())
}

func TestHasPrefixEmpty(t *testing.T) {
	tr := NewTrie()
	if tr.HasPrefix("test") {
		t.Error("Expected no prefix test")
	}
}

func TestHasPrefixOne(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	if !tr.HasPrefix("test") {
		t.Error("Expected prefix test")
	}
}

func TestHasPrefixMany(t *testing.T) {
	tr := NewTrie()
	tr.Add("tease")
	tr.Add("teases")
	tr.Add("teased")
	tr.Add("teaser")
	tr.Add("tests")
	tr.Add("test")
	tr.Add("tested")
	tr.Add("testing")
	if tr.HasPrefix("ted") {
		t.Error("Expected no prefix ted")
	}
	if !tr.HasPrefix("tease") {
		t.Error("Expected prefix tease")
	}
	if !tr.HasPrefix("testing") {
		t.Error("Expected prefix testing")
	}
}

func TestHasEmpty(t *testing.T) {
	tr := NewTrie()
	if tr.Has("test") {
		t.Error("Expected no test")
	}
}

func TestHasOne(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	if !tr.Has("test") {
		t.Error("Expected test")
	}
}

func TestHasMany(t *testing.T) {
	tr := NewTrie()
	tr.Add("tease")
	tr.Add("teases")
	tr.Add("teased")
	tr.Add("teaser")
	tr.Add("tests")
	tr.Add("test")
	tr.Add("tested")
	tr.Add("testing")
	if tr.Has("testi") {
		t.Error("Expected no testi")
	}
	if tr.Has("te") {
		t.Error("Expected no te")
	}
	if !tr.Has("tease") {
		t.Error("Expected tease")
	}
	if !tr.Has("testing") {
		t.Error("Expected testing")
	}
}

func TestHasPrefixManyMultibyte(t *testing.T) {
	tr := NewTrie()
	tr.Add("日本人")
	tr.Add("人")
	tr.Add("日本")
	tr.Add("日本語学校")
	tr.Add("学校")
	tr.Add("日本語")
	if tr.HasPrefix("ä") {
		t.Error("Expected no prefix ä")
	}
	if tr.HasPrefix("無い") {
		t.Error("Expected no prefix 無い")
	}
	if !tr.HasPrefix("日本語") {
		t.Error("Expected prefix 日本語")
	}
	if !tr.HasPrefix("日") {
		t.Error("Expected prefix 日")
	}
}

func TestHasManyMultibyte(t *testing.T) {
	tr := NewTrie()
	tr.Add("日本人")
	tr.Add("人")
	tr.Add("日本")
	tr.Add("日本語学校")
	tr.Add("学校")
	tr.Add("日本語")
	if tr.Has("ä") {
		t.Error("Expected no ä")
	}
	if tr.Has("無い") {
		t.Error("Expected no 無い")
	}
	if tr.Has("日") {
		t.Error("Expected no 日")
	}
	if !tr.Has("日本語") {
		t.Error("Expected 日本語")
	}
	if !tr.Has("学校") {
		t.Error("Expected 学校")
	}
}

func TestDeleteEmpty(t *testing.T) {
	tr := NewTrie()
	if tr.Delete("test") {
		t.Error("Expected false for tr.Delete('test')")
	}
}

func TestDeleteOne(t *testing.T) {
	tr := NewTrie()
	tr.Add("test")
	if !tr.Delete("test") {
		t.Error("Expected true for tr.Delete('test')")
	}
}

func TestDeleteMany(t *testing.T) {
	tr := NewTrie()
	tr.Add("tease")
	tr.Add("teases")
	tr.Add("teased")
	tr.Add("test")

	if tr.Delete("te") {
		t.Error("Expected false for tr.Delete('te')")
	}
	if !tr.Delete("test") {
		t.Error("Expected true for tr.Delete('test')")
	}

	expectedMembers := make(map[string]bool)
	expectedMembers["tease"] = true
	expectedMembers["teases"] = true
	expectedMembers["teased"] = true
	// expectedMembers["test"] = true
	for _, m := range tr.Members() {
		if m == "test" {
			t.Error("Expected 'test' to be deleted")
		} else {
			delete(expectedMembers, m)
		}
	}

	if len(expectedMembers) != 0 {
		t.Error("Deletion seems to have deleted more than just 'test'.", expectedMembers)
	}
	if !tr.Delete("tease") {
		t.Error("Expected true for tr.Delete('tease')")
	}
	if !tr.Delete("teases") {
		t.Error("Expected true for tr.Delete('tease')")
	}
	if !tr.Delete("teased") {
		t.Error("Expected true for tr.Delete('tease')")
	}

	// tr.PrintDump()

	if len(tr.Root.Branches) != 0 {
		t.Error("Expected 0 Branches on Root")
	}
	if len(tr.Root.LeafValue) != 0 {
		t.Error("Expected no LeafValue on Root")
	}
	if tr.Root.End {
		t.Error("Expected End to be false on Root")
	}
}

func TestDeleteManyRandom_az(t *testing.T) {
	tr := NewTrie()
	var prefix = "prefix"
	var words []string
	var str []byte
	var n = 0
	for n < 100 {
		i := 0
		str = []byte{}
		for i < 10 {
			rn := 0
			for rn < 97 {
				rn = rand.Intn(123)
			}
			str = append(str, byte(rn))
			i++
		}
		if rand.Intn(2) == 1 {
			words = append(words, prefix+string(str))
			tr.Add(prefix + string(str))
		} else {
			words = append(words, string(str))
			tr.Add(string(str))
		}
		n++
	}
	// t.Log(words)
	// tr.PrintDump()
	for wi, w := range words {
		if !tr.Delete(w) {
			t.Errorf("Expected true for tr.Delete('%s')", w)
		}
		// expect to still find the rest
		if wi+1 < len(words) {
			for _, ow := range words[wi+1:] {
				// t.Logf("Checking for %s", ow)
				if !tr.Has(ow) {
					t.Errorf("Expected to still find %s", ow)
				}
			}
		}
	}
	tr.PrintDump()
	if len(tr.Root.Branches) != 0 {
		t.Error("Expected 0 Branches on Root")
	}
	if len(tr.Root.LeafValue) != 0 {
		t.Error("Expected no LeafValue on Root")
	}
	if tr.Root.End {
		t.Error("Expected End to be false on Root")
	}
}

func _TestMultiAdd(t *testing.T) {
	tr := NewTrie()
	words := []string{"foodie", "foods", "foodchain", "foodcrave", "food", "人", "日本", "日本語学校", "学校", "日本語"}
	wg := sync.WaitGroup{}
	for _, w := range words {
		wg.Add(1)
		go func(word string) {
			tr.Add(word)
			wg.Done()
		}(w)
	}
	wg.Wait()
	tr.PrintDump()
}

func TestDumpToFileLoadFromFile(t *testing.T) {
	tr := NewTrie()
	var prefix = "prefix"
	var words []string
	var str []byte
	var n = 0
	for n < 100 {
		i := 0
		str = []byte{}
		for i < 10 {
			rn := 0
			for rn < 97 {
				rn = rand.Intn(123)
			}
			str = append(str, byte(rn))
			i++
		}
		if rand.Intn(2) == 1 {
			words = append(words, prefix+string(str))
			tr.Add(prefix + string(str))
		} else {
			words = append(words, string(str))
			tr.Add(string(str))
		}
		n++
	}
	tr.DumpToFile("testfiles/TestDumpToFileLoadFromFile")

	loadedTrie, err := LoadFromFile("testfiles/TestDumpToFileLoadFromFile")
	if err != nil {
		t.Errorf("Failed to load Trie from file: %v", err)
	}
	for _, w := range words {
		// t.Logf("Checking for %s", w)
		if !loadedTrie.Has(w) {
			t.Errorf("Expected to find %s", w)
		}
	}

	trMembers := set.NewStringSet(tr.Members()...)
	loadedTrieMembers := set.NewStringSet(loadedTrie.Members()...)

	t.Log("trMembers.IsEqual(loadedTrieMembers):", trMembers.IsEqual(loadedTrieMembers))

	diff := trMembers.Difference(loadedTrieMembers)
	if diff.Len() > 0 {
		t.Error("Dump() of the original and the LoadFromFile() version of the Trie are different.")
	}
}

func TestLoadFromFileEmpty(t *testing.T) {
	loadedTrie, err := LoadFromFile("testfiles/empty")
	if err != nil {
		t.Errorf("Failed to load Trie from file: %v", err)
	}

	loadedTrieMembers := set.NewStringSet(loadedTrie.Members()...)
	t.Log(loadedTrieMembers)
	t.Log(loadedTrieMembers.Len())
	if loadedTrieMembers.Len() > 0 {
		t.Error("Expected 0 Members from LoadFromFile() with an empty file.")
	}
}

// some simple benchmarks

func BenchmarkBenchAdd(b *testing.B) {
	for x := 0; x < b.N; x++ {
		tr := NewTrie()
		tr.Add("foodie")
		tr.Add("foods")
		tr.Add("foodchain")
		tr.Add("foodcrave")
		tr.Add("food")
		tr.Add("人")
		tr.Add("日本")
		tr.Add("日本語学校")
		tr.Add("学校")
		tr.Add("日本語")
	}
}

func BenchmarkBenchHasPrefix(b *testing.B) {
	tr := NewTrie()
	tr.Add("foodie")
	tr.Add("foods")
	tr.Add("foodchain")
	tr.Add("foodcrave")
	tr.Add("food")
	tr.Add("人")
	tr.Add("日本")
	tr.Add("日本語学校")
	tr.Add("学校")
	tr.Add("日本語")

	for x := 0; x < b.N; x++ {
		tr.HasPrefix("foodie")
		tr.HasPrefix("foodcrave")
		tr.HasPrefix("日本")
		tr.HasPrefix("学校")
		tr.HasPrefix("thisisnotinhere")
		tr.HasPrefix("学日本校")
	}
}

func BenchmarkBenchHas(b *testing.B) {
	tr := NewTrie()
	tr.Add("foodie")
	tr.Add("foods")
	tr.Add("foodchain")
	tr.Add("foodcrave")
	tr.Add("food")
	tr.Add("人")
	tr.Add("日本")
	tr.Add("日本語学校")
	tr.Add("学校")
	tr.Add("日本語")

	for x := 0; x < b.N; x++ {
		tr.Has("foodie")
		tr.Has("foodcrave")
		tr.Has("日本")
		tr.Has("学校")
		tr.Has("thisisnotinhere")
		tr.Has("学日本校")
	}
}
