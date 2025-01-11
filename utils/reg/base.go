package reg

import (
	"github.com/dlclark/regexp2"
)

type Reg struct {
	*regexp2.Regexp
}

func (r *Reg) MatchString(s string) (isMatch bool) {
	isMatch, _ = r.Regexp.MatchString(s)
	return
}

func (r *Reg) FindString(s string) (dest string) {
	match, _ := r.Regexp.FindStringMatch(s)
	if match == nil {
		return
	}
	return match.String()
}

func (r *Reg) FindAllString(s string) (dest []string) {
	match, _ := r.Regexp.FindStringMatch(s)
	if match == nil {
		return
	}
	for match != nil {
		dest = append(dest, match.String())
		match, _ = r.Regexp.FindNextMatch(match)
	}
	return
}

func (r *Reg) ReplaceAllString(src string, repl string) (dest string) {
	dest, _ = r.Regexp.Replace(src, repl, -1, -1)
	return
}

func (r *Reg) FindStringIndex(s string) (loc, length int) {
	match, _ := r.Regexp.FindStringMatch(s)
	if match == nil {
		return
	}
	return match.Index, match.Length

}

func (r *Reg) MustCompile(s string) {
	r.Regexp = regexp2.MustCompile(s, regexp2.None)
}
