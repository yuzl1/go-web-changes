package engine

import (
	"context"
	"errors"
	"fmt"
	"monitor-server/internal/model"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type StepResult struct {
	StepOrder int    `json:"step_order"`
	Status    string `json:"status"`
	ElapsedMs int64  `json:"elapsed_ms"`
	Output    string `json:"output"`
	Count     int    `json:"count"`    // 匹配元素数（自动检测）
	Error     string `json:"error,omitempty"`
}

func browserOpts() []chromedp.ExecAllocatorOption {
	return append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
}

func newBrowserCtx(timeout time.Duration) (context.Context, context.CancelFunc) {
	allocCtx, ca := chromedp.NewExecAllocator(context.Background(), browserOpts()...)
	ctx, cb := chromedp.NewContext(allocCtx)
	ctx, cc := context.WithTimeout(ctx, timeout)
	return ctx, func() { cc(); cb(); ca() }
}

// ExecuteScan 扫描：加载页面→注入jQuery→执行脚本→返回结果
func ExecuteScan(url string, rules []model.ScanRule) (string, error) {
	ctx, cancel := newBrowserCtx(30 * time.Second)
	defer cancel()
	if err := navigateAndInject(ctx, url); err != nil {
		return "", err
	}
	return executeChain(ctx, rules)
}

// ExecuteTest 测试规则
func ExecuteTest(url string, rules []model.ScanRule) ([]StepResult, string, error) {
	ctx, cancel := newBrowserCtx(20 * time.Second)
	defer cancel()
	if err := navigateAndInject(ctx, url); err != nil {
		return nil, "", err
	}
	return executeChainWithSteps(ctx, rules)
}

func navigateAndInject(ctx context.Context, url string) error {
	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		return errors.New("页面加载失败: " + err.Error())
	}
	chromedp.Run(ctx, chromedp.Sleep(500*time.Millisecond))
	// 注入 jQuery
	var hasJQ bool
	chromedp.Run(ctx, chromedp.Evaluate(`typeof jQuery!=='undefined'`, &hasJQ))
	if !hasJQ {
		chromedp.Run(ctx, chromedp.Evaluate(
			`(function(){var s=document.createElement('script');s.src='https://cdn.jsdelivr.net/npm/jquery@3.7.1/dist/jquery.min.js';document.head.appendChild(s);})()`, nil))
		for i := 0; i < 20; i++ {
			chromedp.Run(ctx, chromedp.Sleep(500*time.Millisecond))
			var loaded bool
			chromedp.Run(ctx, chromedp.Evaluate(`typeof jQuery!=='undefined'`, &loaded))
			if loaded {
				break
			}
		}
	}
	return nil
}

// escapeJS 转义字符串用于嵌入 JS 单引号字符串
func escapeJS(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}

func executeChain(ctx context.Context, rules []model.ScanRule) (string, error) {
	var final string
	for _, r := range rules {
		var result string
		if err := chromedp.Run(ctx, chromedp.Evaluate(r.RuleContent, &result)); err != nil {
			if r.RuleMode == 2 { continue }
			return "", fmt.Errorf("步骤%d: %v", r.StepOrder, err)
		}
		final = strings.TrimSpace(result)
		escaped := escapeJS(final)
		chromedp.Run(ctx, chromedp.Evaluate(
			fmt.Sprintf("window.__prev_result=%q;window.__prev_html=%q", escaped, escaped), nil))
	}
	return final, nil
}

func executeChainWithSteps(ctx context.Context, rules []model.ScanRule) ([]StepResult, string, error) {
	// 初始化 __prev
	chromedp.Run(ctx, chromedp.Evaluate("window.__prev_result='';window.__prev_html=''", nil))
	var results []StepResult
	var prev string
	hasError := false
	for _, r := range rules {
		start := time.Now()
		script := r.RuleContent
		var result string
		err := chromedp.Run(ctx, chromedp.Evaluate(script, &result))
		elapsed := time.Since(start).Milliseconds()
		sr := StepResult{StepOrder: r.StepOrder, ElapsedMs: elapsed}
		if err != nil {
			sr.Status = "error"
			sr.Error = err.Error()
			hasError = true
			if r.RuleMode == 2 {
				sr.Status = "skipped"
				continue
			}
		} else {
			sr.Status = "success"
			sr.Output = strings.TrimSpace(result)
			sr.Count = detectElementCount(ctx, r.RuleContent)
			// 存入 window 供下一步引用
			escaped := escapeJS(sr.Output)
			chromedp.Run(ctx, chromedp.Evaluate(
				fmt.Sprintf("window.__prev_result=%q;window.__prev_html=%q", escaped, escaped), nil))
			prev = sr.Output
		}
		results = append(results, sr)
	}
	if hasError {
		return results, prev, errors.New("存在失败步骤")
	}
	return results, prev, nil
}

// jqSelectorRE 匹配 jQuery 选择器: $('...') 或 $("...")
var jqSelectorRE1 = regexp.MustCompile(`\$\s*\(\s*'([^']+)'\s*\)`)
var jqSelectorRE2 = regexp.MustCompile(`\$\s*\(\s*"([^"]+)"\s*\)`)

// detectElementCount 自动检测脚本中jQuery选择器匹配的元素数
func detectElementCount(ctx context.Context, script string) int {
	// 尝试匹配 $('...') 或 $("...")
	m := jqSelectorRE1.FindStringSubmatch(script)
	if len(m) < 2 {
		m = jqSelectorRE2.FindStringSubmatch(script)
	}
	if len(m) < 2 {
		return -1
	}
	selector := m[1]
	checkJS := fmt.Sprintf(`(function(){try{return $('%s').length}catch(e){return -1}})()`, strings.ReplaceAll(selector, "'", "\\'"))
	var count int
	if err := chromedp.Run(ctx, chromedp.Evaluate(checkJS, &count)); err != nil {
		return -1
	}
	return count
}

func IsChanged(a, b string) bool { return strings.TrimSpace(a) != strings.TrimSpace(b) }

func FetchPageHTML(url string) (string, error) {
	ctx, cancel := newBrowserCtx(25 * time.Second)
	defer cancel()
	if err := navigateAndInject(ctx, url); err != nil {
		return "", err
	}
	var html string
	chromedp.Run(ctx, chromedp.Evaluate("document.documentElement.outerHTML", &html))
	return html, nil
}
