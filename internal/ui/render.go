package ui

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/daihao4371/hostmanager/internal/theme"
)

// 高级绘制引擎
type RenderEngine struct {
	theme       *theme.Theme
	frameBuffer [][]termbox.Cell // 帧缓冲区
	width       int
	height      int
	animations  []*Animation // 活跃动画列表
}

// 创建渲染引擎
func NewRenderEngine(t *theme.Theme) *RenderEngine {
	width, height := termbox.Size()

	// 初始化帧缓冲区
	buffer := make([][]termbox.Cell, height)
	for i := range buffer {
		buffer[i] = make([]termbox.Cell, width)
	}

	return &RenderEngine{
		theme:       t,
		frameBuffer: buffer,
		width:       width,
		height:      height,
		animations:  []*Animation{},
	}
}

// 渲染卡片组件
func (r *RenderEngine) RenderCard(x, y, width, height int, style ComponentStyle, content string, state ComponentState) {
	// 绘制阴影效果
	if style.Shadow.Enabled {
		r.drawShadow(x+style.Shadow.OffsetX, y+style.Shadow.OffsetY, width, height)
	}

	// 绘制渐变背景
	r.drawGradientBackground(x, y, width, height, style.Background)

	// 绘制圆角边框
	r.drawRoundedBorder(x, y, width, height, style.Border)

	// 绘制内容区域
	contentX := x + style.Padding.Left
	contentY := y + style.Padding.Top
	contentWidth := width - style.Padding.Left - style.Padding.Right

	// 根据状态调整颜色
	textColor := r.getStateColor(r.theme.Foreground, state)
	r.drawTextWithStyle(contentX, contentY, content, textColor, contentWidth)
}

// 智能状态指示器增强版
func (r *RenderEngine) RenderEnhancedStatusBadge(x, y int, statusType StatusType, text string, iconSet *IconSet, animated bool, frame int) {
	var bgColor, fgColor termbox.Attribute
	var icon string

	// 获取对应图标
	if ico, exists := iconSet.Status[statusType]; exists {
		icon = ico
	} else {
		icon = "○"
	}

	// 根据状态类型设置颜色
	switch statusType {
	case StatusOnline:
		bgColor = r.theme.Success
		fgColor = termbox.ColorWhite
	case StatusOffline:
		bgColor = r.theme.Error
		fgColor = termbox.ColorWhite
	case StatusLoading, StatusConnecting:
		bgColor = r.theme.Warning
		fgColor = termbox.ColorBlack
		// 动画效果
		if animated {
			loadingIcons := []string{"◐", "◓", "◑", "◒"}
			icon = loadingIcons[frame%len(loadingIcons)]
		}
	case StatusError:
		bgColor = r.theme.Error
		fgColor = termbox.ColorWhite
	case StatusWarning:
		bgColor = r.theme.Warning
		fgColor = termbox.ColorBlack
	case StatusMaintenance:
		bgColor = r.theme.Info
		fgColor = termbox.ColorWhite
	default:
		bgColor = r.theme.Border
		fgColor = r.theme.Foreground
	}

	// 绘制状态指示器背景
	r.setCell(x, y, ' ', fgColor, bgColor)
	r.setCell(x+1, y, ' ', fgColor, bgColor)

	// 绘制图标
	r.setCell(x, y, rune(icon[0]), fgColor, bgColor)

	// 绘制状态文本
	if text != "" {
		r.drawText(x+3, y, text, r.theme.Foreground)
	}
}

// 绘制状态指示器（兼容性保持）
func (r *RenderEngine) RenderStatusBadge(x, y int, status, text string) {
	var statusType StatusType
	switch status {
	case "online":
		statusType = StatusOnline
	case "offline":
		statusType = StatusOffline
	case "loading":
		statusType = StatusLoading
	default:
		statusType = StatusIdle
	}

	iconSet := CreateIconSet()
	r.RenderEnhancedStatusBadge(x, y, statusType, text, iconSet, true, 0)
}

// 绘制高级进度条
func (r *RenderEngine) RenderAdvancedProgressBar(x, y, width int, progress float32, segments int, animated bool) {
	// 背景条
	bgChar := '░'
	for i := 0; i < width; i++ {
		r.setCell(x+i, y, bgChar, r.theme.Border, r.theme.Background)
	}

	if segments > 1 {
		// 分段进度条
		segmentWidth := width / segments
		for seg := 0; seg < segments; seg++ {
			segStart := x + seg*segmentWidth
			segEnd := segStart + segmentWidth - 1

			// 计算当前段的进度
			segmentProgress := progress*float32(segments) - float32(seg)
			if segmentProgress > 1.0 {
				segmentProgress = 1.0
			} else if segmentProgress < 0.0 {
				segmentProgress = 0.0
			}

			// 绘制段进度
			segFillWidth := int(float32(segmentWidth-1) * segmentProgress)
			fillChar := '█'

			for i := 0; i < segFillWidth; i++ {
				if segStart+i < segEnd {
					color := r.getProgressColor(float32(seg) / float32(segments-1))
					r.setCell(segStart+i, y, fillChar, color, r.theme.Background)
				}
			}

			// 段分隔符
			if seg < segments-1 {
				r.setCell(segEnd, y, '│', r.theme.Border, r.theme.Background)
			}
		}
	} else {
		// 标准进度条
		fillWidth := int(float32(width) * progress)
		fillChar := '█'

		for i := 0; i < fillWidth; i++ {
			color := r.getProgressColor(float32(i) / float32(width))
			r.setCell(x+i, y, fillChar, color, r.theme.Background)
		}
	}

	// 动画光泽效果
	if animated && progress > 0 && progress < 1.0 {
		glowPos := int(float32(width) * progress)
		if glowPos > 0 && glowPos < width {
			r.setCell(x+glowPos, y, '◆', termbox.ColorWhite, r.theme.Background)
		}
	}

	// 进度文本
	progressText := fmt.Sprintf("%.1f%%", progress*100)
	textX := x + (width-len(progressText))/2
	r.drawText(textX, y+1, progressText, r.theme.Foreground)
}

// 获取渐变进度颜色
func (r *RenderEngine) getProgressColor(ratio float32) termbox.Attribute {
	if ratio < 0.3 {
		return r.theme.Error
	} else if ratio < 0.7 {
		return r.theme.Warning
	} else {
		return r.theme.Success
	}
}

// 绘制进度条
func (r *RenderEngine) RenderProgressBar(x, y, width int, progress float32) {
	r.RenderAdvancedProgressBar(x, y, width, progress, 1, true)
}

// 私有辅助方法

// 设置单元格
func (r *RenderEngine) setCell(x, y int, ch rune, fg, bg termbox.Attribute) {
	if x >= 0 && x < r.width && y >= 0 && y < r.height {
		termbox.SetCell(x, y, ch, fg, bg)
	}
}

// 绘制文本
func (r *RenderEngine) drawText(x, y int, text string, color termbox.Attribute) {
	for i, ch := range text {
		r.setCell(x+i, y, ch, color, r.theme.Background)
	}
}

// 绘制居中文本
func (r *RenderEngine) drawCenteredText(x, y, width int, text string, color termbox.Attribute) {
	textX := x + (width-len(text))/2
	r.drawText(textX, y, text, color)
}

// 绘制带样式的文本
func (r *RenderEngine) drawTextWithStyle(x, y int, text string, color termbox.Attribute, maxWidth int) {
	if len(text) > maxWidth {
		text = text[:maxWidth-3] + "..."
	}
	r.drawText(x, y, text, color)
}

// 绘制阴影
func (r *RenderEngine) drawShadow(x, y, width, height int) {
	shadowChar := '░'
	shadowColor := termbox.Attribute(8) // 深灰色

	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			if x+dx < r.width && y+dy < r.height {
				r.setCell(x+dx, y+dy, shadowChar, shadowColor, r.theme.Background)
			}
		}
	}
}

// 绘制渐变背景
func (r *RenderEngine) drawGradientBackground(x, y, width, height int, bg BackgroundStyle) {
	if bg.Type == BackgroundGradient {
		for dy := 0; dy < height; dy++ {
			// 计算渐变比例
			ratio := float32(dy) / float32(height-1)
			// 简化的颜色插值（使用交替字符模拟渐变）
			char := ' '
			if int(ratio*4)%2 == 0 {
				char = '░'
			}

			for dx := 0; dx < width; dx++ {
				if x+dx < r.width && y+dy < r.height {
					r.setCell(x+dx, y+dy, char, bg.Color, r.theme.Background)
				}
			}
		}
	} else {
		// 实心背景
		for dy := 0; dy < height; dy++ {
			for dx := 0; dx < width; dx++ {
				if x+dx < r.width && y+dy < r.height {
					r.setCell(x+dx, y+dy, ' ', bg.Color, bg.Color)
				}
			}
		}
	}
}

// 绘制圆角边框
func (r *RenderEngine) drawRoundedBorder(x, y, width, height int, border BorderStyle) {
	if border.Type == BorderNone {
		return
	}

	var borderChars struct {
		horizontal, vertical, topLeft, topRight, bottomLeft, bottomRight rune
	}

	switch border.Type {
	case BorderSolid:
		borderChars = struct {
			horizontal, vertical, topLeft, topRight, bottomLeft, bottomRight rune
		}{'─', '│', '┌', '┐', '└', '┘'}
	case BorderDouble:
		borderChars = struct {
			horizontal, vertical, topLeft, topRight, bottomLeft, bottomRight rune
		}{'═', '║', '╔', '╗', '╚', '╝'}
	case BorderRounded:
		borderChars = struct {
			horizontal, vertical, topLeft, topRight, bottomLeft, bottomRight rune
		}{'─', '│', '╭', '╮', '╰', '╯'}
	default:
		borderChars = struct {
			horizontal, vertical, topLeft, topRight, bottomLeft, bottomRight rune
		}{'─', '│', '┌', '┐', '└', '┘'}
	}

	// 绘制水平线
	for dx := 1; dx < width-1; dx++ {
		r.setCell(x+dx, y, borderChars.horizontal, border.Color, r.theme.Background)
		r.setCell(x+dx, y+height-1, borderChars.horizontal, border.Color, r.theme.Background)
	}

	// 绘制垂直线
	for dy := 1; dy < height-1; dy++ {
		r.setCell(x, y+dy, borderChars.vertical, border.Color, r.theme.Background)
		r.setCell(x+width-1, y+dy, borderChars.vertical, border.Color, r.theme.Background)
	}

	// 绘制角点
	r.setCell(x, y, borderChars.topLeft, border.Color, r.theme.Background)
	r.setCell(x+width-1, y, borderChars.topRight, border.Color, r.theme.Background)
	r.setCell(x, y+height-1, borderChars.bottomLeft, border.Color, r.theme.Background)
	r.setCell(x+width-1, y+height-1, borderChars.bottomRight, border.Color, r.theme.Background)
}

// 根据状态获取颜色
func (r *RenderEngine) getStateColor(baseColor termbox.Attribute, state ComponentState) termbox.Attribute {
	switch state {
	case StateHover:
		return r.theme.Highlight
	case StatePressed:
		return r.theme.Warning
	case StateFocused:
		return r.theme.Info
	case StateDisabled:
		return termbox.Attribute(8) // 灰色
	case StateLoading:
		return r.theme.Warning
	default:
		return baseColor
	}
}

// 绘制高级加载动画
func (r *RenderEngine) RenderAdvancedSpinner(x, y int, spinnerType string, frame int, iconSet *IconSet) {
	var icon string

	switch spinnerType {
	case "dots":
		spinnerFrames := []string{"◐", "◓", "◑", "◒"}
		icon = spinnerFrames[frame%len(spinnerFrames)]
	case "bars":
		spinnerFrames := []string{"▁", "▃", "▅", "▇", "▅", "▃"}
		icon = spinnerFrames[frame%len(spinnerFrames)]
	case "arrows":
		spinnerFrames := []string{"↑", "↗", "→", "↘", "↓", "↙", "←", "↖"}
		icon = spinnerFrames[frame%len(spinnerFrames)]
	case "pulse":
		spinnerFrames := []string{"●", "◉", "○", "◉"}
		icon = spinnerFrames[frame%len(spinnerFrames)]
	default:
		spinnerFrames := []string{"◐", "◓", "◑", "◒"}
		icon = spinnerFrames[frame%len(spinnerFrames)]
	}

	// 动态颜色循环
	colors := []termbox.Attribute{r.theme.Info, r.theme.Success, r.theme.Warning}
	color := colors[frame%len(colors)]

	r.setCell(x, y, rune(icon[0]), color, r.theme.Background)
}

// 绘制加载动画（兼容性保持）
func (r *RenderEngine) RenderSpinner(x, y int, frame int) {
	iconSet := CreateIconSet()
	r.RenderAdvancedSpinner(x, y, "dots", frame, iconSet)
}

// 绘制工具提示
func (r *RenderEngine) RenderTooltip(x, y int, text string) {
	tooltipWidth := len(text) + 4
	tooltipHeight := 3

	// 确保在屏幕边界内
	if x < 0 {
		x = 0
	}
	if x+tooltipWidth >= r.width {
		x = r.width - tooltipWidth - 1
	}
	if y < 0 {
		y = 0
	}
	if y+tooltipHeight >= r.height {
		y = r.height - tooltipHeight - 1
	}

	// 绘制工具提示背景
	for dy := 0; dy < tooltipHeight; dy++ {
		for dx := 0; dx < tooltipWidth; dx++ {
			r.setCell(x+dx, y+dy, ' ', r.theme.Foreground, r.theme.Background)
		}
	}

	// 绘制提示文本
	textX := x + 2
	textY := y + 1
	r.drawText(textX, textY, text, r.theme.Foreground)
}

// 绘制交互反馈效果
func (r *RenderEngine) RenderClickFeedback(x, y int, frame int) {
	if frame > 10 {
		return
	}

	// 简单的扩散效果
	char := '●'
	color := r.theme.Info
	r.setCell(x, y, char, color, r.theme.Background)
}

// 启动动画
func (r *RenderEngine) StartAnimation(component *UIComponent) {
	if component.Animation != nil && component.Animation.Active {
		return
	}

	animation := &Animation{
		Config:    component.Style.Animation,
		StartTime: time.Now(),
		Progress:  0.0,
		Active:    true,
	}

	component.Animation = animation
	r.animations = append(r.animations, animation)
}

// 更新所有动画
func (r *RenderEngine) UpdateAnimations() {
	currentTime := time.Now()
	var activeAnimations []*Animation

	for _, anim := range r.animations {
		if !anim.Active {
			continue
		}

		elapsed := currentTime.Sub(anim.StartTime)
		if elapsed >= anim.Config.Duration {
			anim.Progress = 1.0
			anim.Active = false
		} else {
			rawProgress := float32(elapsed) / float32(anim.Config.Duration)
			anim.Progress = applyEasing(rawProgress, anim.Config.Easing)
		}

		if anim.Active {
			activeAnimations = append(activeAnimations, anim)
		}
	}

	r.animations = activeAnimations
}

// 检查是否有活跃动画
func (r *RenderEngine) HasActiveAnimations() bool {
	return len(r.animations) > 0
}
