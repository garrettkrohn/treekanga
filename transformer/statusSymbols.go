package transformer

import (
	"fmt"
	"strings"

	"github.com/garrettkrohn/treekanga/models"
)

// StatusLegend documents the compact symbols rendered by WorktreeStatusSymbols.
const StatusLegend = "status legend: + staged, * modified, ? untracked, ↑/↓ ahead/behind default branch, ⇡/⇣ ahead/behind remote, ✓ merged"

// WorktreeStatusSymbols renders a worktree's R1-R4 status fields as a
// compact, worktrunk-style symbol string. Indicators that carry no signal
// (e.g. no upstream, not merged) are omitted rather than shown as "empty".
func WorktreeStatusSymbols(worktree models.Worktree) string {
	var parts []string

	if dirty := DirtySymbols(worktree); dirty != "" {
		parts = append(parts, dirty)
	}
	if ahead := DefaultAheadBehindSymbols(worktree); ahead != "" {
		parts = append(parts, ahead)
	}
	if remote := RemoteAheadBehindSymbols(worktree); remote != "" {
		parts = append(parts, remote)
	}
	if merged := MergedSymbol(worktree); merged != "" {
		parts = append(parts, merged)
	}

	return strings.Join(parts, " ")
}

// DirtySymbols renders the R1 working-tree indicator: staged/modified/untracked.
func DirtySymbols(worktree models.Worktree) string {
	var b strings.Builder
	if worktree.HasStaged {
		b.WriteString("+")
	}
	if worktree.HasModified {
		b.WriteString("*")
	}
	if worktree.HasUntracked {
		b.WriteString("?")
	}
	return b.String()
}

// DefaultAheadBehindSymbols renders the R2 indicator: commits ahead/behind
// the default branch.
func DefaultAheadBehindSymbols(worktree models.Worktree) string {
	return aheadBehindSymbols('↑', '↓', worktree.AheadDefault, worktree.BehindDefault)
}

// RemoteAheadBehindSymbols renders the R3 indicator: commits ahead/behind
// the remote tracking branch. Returns "" when no upstream is configured.
func RemoteAheadBehindSymbols(worktree models.Worktree) string {
	if !worktree.HasUpstream {
		return ""
	}
	return aheadBehindSymbols('⇡', '⇣', worktree.AheadRemote, worktree.BehindRemote)
}

// MergedSymbol renders the R4 indicator: a check mark when the branch's
// content is already present in the default branch.
func MergedSymbol(worktree models.Worktree) string {
	if worktree.Merged == models.MergeStatusMerged {
		return "✓"
	}
	return ""
}

func aheadBehindSymbols(aheadGlyph, behindGlyph rune, ahead, behind int) string {
	var b strings.Builder
	if ahead > 0 {
		fmt.Fprintf(&b, "%c%d", aheadGlyph, ahead)
	}
	if behind > 0 {
		fmt.Fprintf(&b, "%c%d", behindGlyph, behind)
	}
	return b.String()
}
