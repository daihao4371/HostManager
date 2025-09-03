package ui

import (
	"github.com/nsf/termbox-go"
	"time"
)

// 交互管理器
type InteractionManager struct {
	feedbackConfig FeedbackConfig
	hoverTimer     *time.Timer
	hoverTarget    *UIComponent
	clickFeedbacks []ClickFeedback
	tooltips       []Tooltip
}

// 点击反馈效果
type ClickFeedback struct {
	X, Y      int
	StartTime time.Time
	Frame     int
	Active    bool
}

// 工具提示信息
type Tooltip struct {
	X, Y      int
	Text      string
	Direction string
	StartTime time.Time
	Visible   bool
	Target    *UIComponent
}

// 创建交互管理器
func NewInteractionManager(config FeedbackConfig) *InteractionManager {
	return &InteractionManager{
		feedbackConfig: config,
		clickFeedbacks: make([]ClickFeedback, 0),
		tooltips:       make([]Tooltip, 0),
	}
}

// 处理鼠标悬停事件
func (im *InteractionManager) HandleMouseHover(x, y int, component *UIComponent) {
	// 检查是否需要显示工具提示
	if component != nil && component.ID != "" {
		// 设置悬停延迟
		if im.hoverTimer != nil {
			im.hoverTimer.Stop()
		}

		im.hoverTarget = component
		im.hoverTimer = time.AfterFunc(im.feedbackConfig.HoverDelay, func() {
			im.showTooltip(x, y, component)
		})
	} else {
		// 清除悬停状态
		im.clearHover()
	}
}

// 处理点击事件
func (im *InteractionManager) HandleMouseClick(x, y int, component *UIComponent) {
	// 添加点击反馈效果
	if im.feedbackConfig.ClickFeedback {
		feedback := ClickFeedback{
			X:         x,
			Y:         y,
			StartTime: time.Now(),
			Frame:     0,
			Active:    true,
		}
		im.clickFeedbacks = append(im.clickFeedbacks, feedback)
	}

	// 清除工具提示
	im.hideAllTooltips()
}

// 处理键盘事件
func (im *InteractionManager) HandleKeyPress(key termbox.Key, ch rune) bool {
	// 处理ESC键关闭工具提示
	if key == termbox.KeyEsc {
		im.hideAllTooltips()
		return true
	}

	// 处理F1显示帮助
	if key == termbox.KeyF1 {
		im.showKeyboardHints()
		return true
	}

	return false
}

// 更新交互效果
func (im *InteractionManager) Update() {
	// 更新点击反馈动画
	activeClickFeedbacks := make([]ClickFeedback, 0)
	for _, feedback := range im.clickFeedbacks {
		if feedback.Active {
			elapsed := time.Since(feedback.StartTime)
			if elapsed < 500*time.Millisecond {
				feedback.Frame = int(elapsed / (50 * time.Millisecond))
				activeClickFeedbacks = append(activeClickFeedbacks, feedback)
			}
		}
	}
	im.clickFeedbacks = activeClickFeedbacks

	// 更新工具提示
	activeTooltips := make([]Tooltip, 0)
	for _, tooltip := range im.tooltips {
		if tooltip.Visible {
			// 工具提示在5秒后自动消失
			if time.Since(tooltip.StartTime) < 5*time.Second {
				activeTooltips = append(activeTooltips, tooltip)
			}
		}
	}
	im.tooltips = activeTooltips
}

// 渲染交互效果
func (im *InteractionManager) Render(renderer *RenderEngine) {
	// 渲染点击反馈效果
	for _, feedback := range im.clickFeedbacks {
		if feedback.Active {
			renderer.RenderClickFeedback(feedback.X, feedback.Y, feedback.Frame)
		}
	}

	// 渲染工具提示
	for _, tooltip := range im.tooltips {
		if tooltip.Visible {
			renderer.RenderTooltip(tooltip.X, tooltip.Y, tooltip.Text)
		}
	}
}

// 显示工具提示
func (im *InteractionManager) showTooltip(x, y int, component *UIComponent) {
	tooltipText := im.getTooltipText(component)
	if tooltipText == "" {
		return
	}

	tooltip := Tooltip{
		X:         x,
		Y:         y,
		Text:      tooltipText,
		Direction: "bottom",
		StartTime: time.Now(),
		Visible:   true,
		Target:    component,
	}

	im.tooltips = append(im.tooltips, tooltip)
}

// 获取工具提示文本
func (im *InteractionManager) getTooltipText(component *UIComponent) string {
	switch component.Type {
	case ComponentButton:
		return "点击执行操作"
	case ComponentList:
		return "使用上下键选择，回车确认"
	case ComponentCard:
		return "双击查看详情"
	case ComponentModal:
		return "按ESC键关闭"
	default:
		return component.ID
	}
}

// 显示键盘快捷键提示
func (im *InteractionManager) showKeyboardHints() {
	hints := []struct {
		key    string
		action string
	}{
		{"↑/↓", "选择项目"},
		{"Enter", "确认选择"},
		{"ESC", "取消/返回"},
		{"Tab", "切换焦点"},
		{"F1", "显示帮助"},
		{"Ctrl+C", "退出程序"},
		{"Space", "切换选择"},
		{"Del", "删除项目"},
	}

	// 创建提示面板
	tooltip := Tooltip{
		X:         5,
		Y:         5,
		Text:      im.formatKeyHints(hints),
		Direction: "bottom",
		StartTime: time.Now(),
		Visible:   true,
		Target:    nil,
	}

	im.tooltips = append(im.tooltips, tooltip)
}

// 格式化快捷键提示
func (im *InteractionManager) formatKeyHints(hints []struct {
	key    string
	action string
}) string {
	result := "快捷键帮助:\n"
	for _, hint := range hints {
		result += hint.key + ": " + hint.action + "\n"
	}
	return result
}

// 清除悬停状态
func (im *InteractionManager) clearHover() {
	if im.hoverTimer != nil {
		im.hoverTimer.Stop()
		im.hoverTimer = nil
	}
	im.hoverTarget = nil
}

// 隐藏所有工具提示
func (im *InteractionManager) hideAllTooltips() {
	for i := range im.tooltips {
		im.tooltips[i].Visible = false
	}
	im.tooltips = make([]Tooltip, 0)
}

// 检查是否有活跃的交互效果
func (im *InteractionManager) HasActiveEffects() bool {
	return len(im.clickFeedbacks) > 0 || len(im.tooltips) > 0
}

// 获取当前悬停的组件
func (im *InteractionManager) GetHoverTarget() *UIComponent {
	return im.hoverTarget
}

// 设置反馈配置
func (im *InteractionManager) SetFeedbackConfig(config FeedbackConfig) {
	im.feedbackConfig = config
}
